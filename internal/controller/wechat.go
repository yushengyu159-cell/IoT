package controller

import (
	"context"
	"encoding/base64"
	"fabric-sdk/internal/service"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"
)

type WechatController struct{}

var Wechat = &WechatController{}

// CheckAndGenerateQRCodeReq 检查邮箱并生成二维码请求
type CheckAndGenerateQRCodeReq struct {
	Email string `json:"email" v:"required|email"`
}

// CheckAndGenerateQRCodeResp 检查邮箱并生成二维码响应
type CheckAndGenerateQRCodeResp struct {
	Exists     bool   `json:"exists"`
	Phone      string `json:"phone"`
	QRCodeURL  string `json:"qrcode_url"`
	LoginToken string `json:"login_token"`
}

// CheckAndGenerateQRCode 检查邮箱并生成微信二维码
func (c *WechatController) CheckAndGenerateQRCode(r *ghttp.Request) {
	ctx := context.Background()
	
	var req CheckAndGenerateQRCodeReq
	if err := r.Parse(&req); err != nil {
		g.Log().Error(ctx, "解析请求参数失败:", err)
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 检查邮箱是否存在
	user, exists := service.GetUserRegistration(req.Email)
	if !exists || user == nil {
		// 邮箱不存在
		r.Response.WriteJson(g.Map{
			"code":    200,
			"message": "邮箱不存在",
			"data": CheckAndGenerateQRCodeResp{
				Exists: false,
			},
		})
		return
	}

	// 生成登录token
	loginToken := uuid.New().String()
	
	// 生成二维码
	qrCodeURL, err := generateWeChatQRCode(loginToken)
	if err != nil {
		g.Log().Error(ctx, "生成二维码失败:", err)
		r.Response.WriteJson(g.Map{
			"code":    500,
			"message": "生成二维码失败",
			"data":    nil,
		})
		return
	}

	// 保存登录token到内存缓存（30分钟过期）
	loginData := g.Map{
		"status":     "waiting",
		"email":      req.Email,
		"phone":      user.Phone,
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	
	cacheKey := fmt.Sprintf("wechat_login:%s", loginToken)
	gcache.Set(ctx, cacheKey, loginData, time.Minute*30)

	g.Log().Info(ctx, "生成微信登录二维码成功:", req.Email, "token:", loginToken)

	// 返回结果
	r.Response.WriteJson(g.Map{
		"code":    200,
		"message": "成功",
		"data": CheckAndGenerateQRCodeResp{
			Exists:     true,
			Phone:      user.Phone,
			QRCodeURL:  qrCodeURL,
			LoginToken: loginToken,
		},
	})
}

// CheckLoginStatus 检查微信登录状态
func (c *WechatController) CheckLoginStatus(r *ghttp.Request) {
	ctx := context.Background()
	
	token := r.Get("token").String()
	email := r.Get("email").String()

	if token == "" {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "token不能为空",
			"data":    nil,
		})
		return
	}

	// 从内存缓存获取登录状态
	cacheKey := fmt.Sprintf("wechat_login:%s", token)
	value, err := gcache.Get(ctx, cacheKey)
	if err != nil || value.IsNil() || value.IsEmpty() {
		// token不存在或已过期
		r.Response.WriteJson(g.Map{
			"code":    200,
			"message": "二维码已过期",
			"data": g.Map{
				"status": "expired",
			},
		})
		return
	}

	var loginData map[string]interface{}
	if err := value.Struct(&loginData); err != nil {
		// 如果解析失败，尝试使用interface{}断言
		if interfaceValue, ok := value.Val().(map[string]interface{}); ok {
			loginData = interfaceValue
		} else {
			// 如果都失败，才使用默认值
			g.Log().Warning(ctx, "无法解析登录数据，使用默认值，token:", token, "error:", err)
			loginData = make(map[string]interface{})
			loginData["status"] = "waiting"
			loginData["email"] = email
		}
	}

	// 安全获取status字段
	status, ok := loginData["status"].(string)
	if !ok {
		status = "waiting"
		loginData["status"] = status
	}
	
	// 如果没有email字段，使用请求中的email
	if _, ok := loginData["email"]; !ok {
		loginData["email"] = email
	}

	g.Log().Info(ctx, fmt.Sprintf("查询登录状态，token: %s, status: %s, email: %s", token, status, email))

	r.Response.WriteJson(g.Map{
		"code":    200,
		"message": "成功",
		"data": g.Map{
			"status": status,
			"email":  loginData["email"],
		},
	})
}

// ScanQRCode 扫码后更新状态为已扫描
func (c *WechatController) ScanQRCode(r *ghttp.Request) {
	ctx := context.Background()
	
	var req struct {
		Token string `json:"token" v:"required"`
	}
	
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 从内存缓存获取登录token
	cacheKey := fmt.Sprintf("wechat_login:%s", req.Token)
	value, err := gcache.Get(ctx, cacheKey)
	if err != nil || value.IsNil() {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "token不存在或已过期",
			"data":    nil,
		})
		return
	}

	var existingData map[string]interface{}
	if err := value.Struct(&existingData); err != nil {
		r.Response.WriteJson(g.Map{
			"code":    500,
			"message": "解析数据失败",
			"data":    nil,
		})
		return
	}

	// 扫描后直接确认为已确认状态，跳过确认界面
	existingData["status"] = "scanned"
	existingData["scanned_at"] = time.Now().Format("2006-01-02 15:04:05")
	existingData["confirmed_at"] = time.Now().Format("2006-01-02 15:04:05")

	// 更新内存缓存
	gcache.Set(ctx, cacheKey, existingData, time.Minute*30)

	g.Log().Info(ctx, "微信扫码成功，已自动确认:", req.Token)

	r.Response.WriteJson(g.Map{
		"code":    200,
		"message": "扫码成功，已自动确认",
		"data":    nil,
	})
}

// generateWeChatQRCode 生成微信登录二维码
func generateWeChatQRCode(token string) (string, error) {
	// 生成二维码内容 - 使用HTTP链接，方便微信扫码后打开
	// 格式：http://域名/wechat/confirm?token=xxx
	qrContent := fmt.Sprintf("http://47.238.159.234:8199/static/wechat-confirm.html?token=%s", token)
	
	// 生成PNG图片
	png, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	
	// 转换为base64
	encoded := base64.StdEncoding.EncodeToString(png)
	qrCodeURL := fmt.Sprintf("data:image/png;base64,%s", encoded)
	
	return qrCodeURL, nil
}

// ConfirmWeChatLogin 确认微信登录（模拟接口，实际应该由微信回调）
func (c *WechatController) ConfirmWeChatLogin(r *ghttp.Request) {
	ctx := context.Background()
	
	var req struct {
		Token string `json:"token" v:"required"`
	}
	
	if err := r.Parse(&req); err != nil {
		g.Log().Error(ctx, "解析请求参数失败:", err)
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "请求参数错误",
			"data":    nil,
		})
		return
	}

	// 从内存缓存获取登录token
	cacheKey := fmt.Sprintf("wechat_login:%s", req.Token)
	value, err := gcache.Get(ctx, cacheKey)
	if err != nil || value.IsNil() {
		r.Response.WriteJson(g.Map{
			"code":    400,
			"message": "token不存在或已过期",
			"data":    nil,
		})
		return
	}

	// 更新状态为已确认
	loginData := g.Map{
		"status":       "confirmed",
		"confirmed_at": time.Now().Format("2006-01-02 15:04:05"),
	}
	
	var existingData map[string]interface{}
	if err := value.Struct(&existingData); err == nil {
		if email, ok := existingData["email"].(string); ok {
			loginData["email"] = email
		}
		if phone, ok := existingData["phone"].(string); ok {
			loginData["phone"] = phone
		}
	}

	// 更新内存缓存中的状态
	gcache.Set(ctx, cacheKey, loginData, time.Minute*30)

	g.Log().Info(ctx, "微信登录确认成功:", req.Token)

	r.Response.WriteJson(g.Map{
		"code":    200,
		"message": "登录确认成功",
		"data":    nil,
	})
}

