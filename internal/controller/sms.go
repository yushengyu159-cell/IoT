package controller

import (
	"fabric-sdk/internal/service"
	"fmt"
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type SMSController struct{}

var SMS = new(SMSController)

// 存储验证码（生产环境应使用Redis）
var verificationCodes = make(map[string]*VerificationCode)

type VerificationCode struct {
	Code      string
	Email     string
	Phone     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// SendCode 发送验证码
// POST /api/sms/send-code
func (c *SMSController) SendCode(r *ghttp.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := r.Parse(&req); err != nil || req.Email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Email is required",
		})
		return
	}

	// 检查邮箱是否已注册
	if !service.IsEmailRegistered(req.Email) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    404,
			Message: "Email not registered",
			Data: g.Map{
				"exists": false,
			},
		})
		return
	}

	// 获取用户信息（包括手机号）
	user, exists := service.GetUserRegistration(req.Email)
	if !exists || user.Phone == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "No phone number associated with this account",
		})
		return
	}

	// 生成6位数字验证码
	code := generateVerificationCode()

	// 存储验证码（5分钟有效期）
	verificationCodes[req.Email] = &VerificationCode{
		Code:      code,
		Email:     req.Email,
		Phone:     user.Phone,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		CreatedAt: time.Now(),
	}

	g.Log().Infof(r.Context(), "SMS verification code generated for %s: %s", req.Email, code)

	// 发送短信（这里使用模拟，实际需要接入短信服务商）
	err := sendSMS(user.Phone, code)
	if err != nil {
		g.Log().Errorf(r.Context(), "Failed to send SMS: %v", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "Failed to send verification code",
		})
		return
	}

	// 返回成功，包含脱敏的手机号
	maskedPhone := maskPhoneNumber(user.Phone)

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Verification code sent successfully",
		Data: g.Map{
			"debugCode": code,
			"phone":    maskedPhone,
			"expiresIn": 300,
		},
	})
}

// VerifyCode 验证验证码并登录
// POST /api/sms/verify-code
func (c *SMSController) VerifyCode(r *ghttp.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := r.Parse(&req); err != nil || req.Email == "" || req.Code == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Email and verification code are required",
		})
		return
	}

	// 验证码长度检查
	if len(req.Code) != 6 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Invalid verification code format",
		})
		return
	}

	// 检查验证码
	verification, exists := verificationCodes[req.Email]
	if !exists {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Verification code not found or expired",
		})
		return
	}

	// 检查是否过期
	if time.Now().After(verification.ExpiresAt) {
		delete(verificationCodes, req.Email)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Verification code has expired",
		})
		return
	}

	// 验证码匹配
	if verification.Code != req.Code {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Invalid verification code",
		})
		return
	}

	// 验证成功，清除验证码
	delete(verificationCodes, req.Email)

	// 获取用户完整信息
	user, exists := service.GetUserRegistration(req.Email)
	if !exists {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    404,
			Message: "User not found",
		})
		return
	}

	// 构建用户信息（与邮箱登录格式一致）
	userInfo := map[string]interface{}{
		"Email":        user.Email,
		"Name":         user.FullName,
		"Role":         user.Role,
		"Phone":        user.Phone,
		"Age":          user.Age,
		"DID":          user.DID,
		"RegisterTime": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"isLoggedIn":   true,
	}

	// 判断跳转页面
	var redirectUrl string
	if req.Email == "esgvisa@gmail.com" {
		redirectUrl = "/static/control1.html"
	} else {
		redirectUrl = "/static/dashboard.html"
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Login successful",
		Data: g.Map{
			"userInfo":    userInfo,
			"redirectUrl": redirectUrl,
		},
	})
}

// generateVerificationCode 生成6位数字验证码
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// maskPhoneNumber 脱敏手机号
func maskPhoneNumber(phone string) string {
	if len(phone) <= 7 {
		return phone
	}
	// 显示前3位和后4位，中间用*代替
	return phone[:3] + "****" + phone[len(phone)-4:]
}

// sendSMS 发送短信（模拟）
// TODO: 接入真实的短信服务商（阿里云短信、腾讯云短信等）
func sendSMS(phone, code string) error {
	// 模拟发送短信
	// 实际项目中应该调用短信服务商的API
	g.Log().Infof(nil, "[SMS] Sending verification code %s to phone %s", code, phone)

	// 示例：阿里云短信API调用
	// import "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	// client, _ := dysmsapi.NewClientWithAccessKey(...)
	// request := dysmsapi.CreateSendSmsRequest()
	// request.Scheme = "https://"
	// request.PhoneNumbers = phone
	// request.SignName = "ESG VISA"
	// request.TemplateCode = "SMS_TEMPLATE_ID"
	// request.TemplateParam = `{"code":"` + code + `"}`
	// response, err := client.SendSms(request)

	return nil
}
