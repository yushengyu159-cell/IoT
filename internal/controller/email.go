package controller

import (
	"fabric-sdk/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

type EmailController struct{}

var Email = new(EmailController)

// POST /api/email/send-code
// body: { email, language? }
func (c *EmailController) SendVerificationCode(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email"`
		Language string `json:"language,omitempty"`
	}
	
	if err := r.Parse(&req); err != nil || req.Email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱地址不能为空",
		})
		return
	}
	
	// 发送验证码（支持多语言）
	code, err := service.SendVerificationEmailWithLanguage(req.Email, req.Language)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	
	// 开发环境下返回验证码（生产环境应移除）
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "验证码发送成功",
		Data: map[string]interface{}{
			"email":    req.Email,
			"language": req.Language,
			"code":     code, // 开发环境显示，生产环境应隐藏
			"note":     "开发环境显示验证码，生产环境应隐藏",
		},
	})
}

// POST /api/email/verify-code
// body: { email, code }
func (c *EmailController) VerifyCode(r *ghttp.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	
	if err := r.Parse(&req); err != nil || req.Email == "" || req.Code == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱和验证码不能为空",
		})
		return
	}
	
	// 验证验证码
	valid, err := service.VerifyEmailCode(req.Email, req.Code)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: err.Error(),
		})
		return
	}
	
	if !valid {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "验证码验证失败",
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "验证码验证成功",
		Data: map[string]interface{}{
			"email": req.Email,
			"valid": true,
		},
	})
}

// GET /api/email/status
// query: email
func (c *EmailController) GetVerificationStatus(r *ghttp.Request) {
	email := r.Get("email").String()
	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱地址不能为空",
		})
		return
	}
	
	verification, exists := service.GetVerificationStatus(email)
	if !exists {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    404,
			Message: "验证码不存在",
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取验证状态成功",
		Data: map[string]interface{}{
			"email":      verification.Email,
			"expires_at": verification.ExpiresAt.Format("2006-01-02 15:04:05"),
			"used":       verification.Used,
			"remaining":  verification.ExpiresAt.Sub(verification.ExpiresAt).String(),
		},
	})
}

// POST /api/email/resend-code
// body: { email }
func (c *EmailController) ResendCode(r *ghttp.Request) {
	var req struct {
		Email string `json:"email"`
	}
	
	if err := r.Parse(&req); err != nil || req.Email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱地址不能为空",
		})
		return
	}
	
	// 检查是否在冷却期内（防止频繁发送）
	verification, exists := service.GetVerificationStatus(req.Email)
	if exists && !verification.Used {
		// 如果验证码未过期且未使用，不允许重新发送
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    429,
			Message: "验证码仍在有效期内，请稍后再试",
		})
		return
	}
	
	// 重新发送验证码
	code, err := service.SendVerificationEmail(req.Email)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "验证码重新发送成功",
		Data: map[string]interface{}{
			"email": req.Email,
			"code":  code, // 开发环境显示
			"note":  "开发环境显示验证码，生产环境应隐藏",
		},
	})
}

// POST /api/email/send-password-reset
// body: { email, language? }
func (c *EmailController) SendPasswordResetCode(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email"`
		Language string `json:"language,omitempty"`
	}
	
	if err := r.Parse(&req); err != nil || req.Email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱地址不能为空",
		})
		return
	}
	
	// 发送密码重置验证码（支持多语言）
	code, err := service.SendPasswordResetEmailWithLanguage(req.Email, req.Language)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: err.Error(),
		})
		return
	}
	
	// 开发环境下返回验证码（生产环境应移除）
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "密码重置验证码发送成功",
		Data: map[string]interface{}{
			"email":    req.Email,
			"language": req.Language,
			"code":     code, // 开发环境显示，生产环境应隐藏
			"note":     "开发环境显示验证码，生产环境应隐藏",
		},
	})
}

// CheckEmailExists 检查邮箱是否已注册
func (c *EmailController) CheckEmailExists(r *ghttp.Request) {
	var req struct {
		Email string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确"`
	}
	
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数错误",
			Data:    err.Error(),
		})
		return
	}
	
	exists := service.CheckEmailExists(req.Email)
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "检查完成",
		Data: map[string]interface{}{
			"email":  req.Email,
			"exists": exists,
		},
	})
}

// TestConnection 测试SMTP连接
func (c *EmailController) TestConnection(r *ghttp.Request) {
	var req struct {
		Email string `json:"email" v:"required|email#邮箱不能为空|邮箱格式不正确"`
	}
	
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数错误",
			Data:    err.Error(),
		})
		return
	}
	
	err := service.TestSMTPConnection(req.Email)
	
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "连接测试失败",
			Data: map[string]interface{}{
				"email": req.Email,
				"error": err.Error(),
			},
		})
		return
	}
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "连接测试成功",
		Data: map[string]interface{}{
			"email": req.Email,
			"status": "connected",
		},
	})
}
