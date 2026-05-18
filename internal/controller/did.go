package controller

import (
	"fabric-sdk/internal/middleware"
	"fabric-sdk/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

type DIDController struct{}

var DID = new(DIDController)

// POST /api/did/register
// body: { email, addresses, phone, password, info }
func (c *DIDController) RegisterDID(r *ghttp.Request) {
	var req struct {
		Email     string `json:"email"`
		Addresses string `json:"addresses"` // 逗号分隔
		Phone     string `json:"phone"`
		Password  string `json:"password"`
		Info      string `json:"info"`
	}
	if err := r.Parse(&req); err != nil || req.Email == "" || req.Password == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "参数不完整"})
		return
	}
	// 调用链码（使用与ESG文件上传相同的方式）
	result, err := service.Chaincode.CreateAssetWithMetadata(r.Context(),
		req.Email,  // 资产ID（邮箱）
		"USER_DID", // 资产类型
		0,          // 大小
		req.Email,  // 所有者（邮箱）
		0)          // 评估值
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: err.Error()})
		return
	}
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "DID登记成功(上链)", Data: result})
}

// POST /api/did/verify
// body: { email, password } 或 { did, password }
func (c *DIDController) VerifyDID(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		DID      string `json:"did"`
	}
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "参数解析失败"})
		return
	}

	// 支持 email 或 did 参数
	var email string
	var user *service.UserRegistration
	var exists bool
	if req.Email != "" {
		email = req.Email
		user, exists = service.GetUserRegistration(email)
	} else if req.DID != "" {
		// 如果提供的是 DID，则通过 DID 反查用户
		user, exists = service.GetUserRegistrationByDID(req.DID)
		if exists {
			email = user.Email
		}
	} else {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "请提供 email 或 did 参数"})
		return
	}

	if req.Password == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "密码不能为空"})
		return
	}

	// 首先从数据库获取用户信息进行密码验证（若上一步已拿到user则跳过再查）
	if user == nil {
		user, exists = service.GetUserRegistration(email)
	}
	if !exists {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "用户不存在"})
		return
	}

	// 验证密码是否正确
	if !service.ValidatePasswordHash(req.Password, user.Password) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "密码错误"})
		return
	}

	// 审核拦截：仅 completed 可登录
	if user.Status != "completed" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 403, Message: "账号未审核通过或已被拒绝", Data: map[string]interface{}{"status": user.Status}})
		return
	}
	// 密码验证成功，返回用户信息和token
	token := middleware.GenerateToken(user.Email)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "DID验证成功",
		Data: map[string]interface{}{
			"email":    user.Email,
			"fullName": user.FullName,
			"role":     user.Role,
			"phone":    user.Phone,
			"age":      user.Age,
			"did":      user.DID,
			"status":   user.Status,
			"valid":    true,
			"token":    token,
		},
	})
}
