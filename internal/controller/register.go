package controller

import (
	"fabric-sdk/internal/model"
	"fabric-sdk/internal/service"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

type RegisterController struct{}

var Register = new(RegisterController)

// POST /api/register/step1
// 第一步：基础信息注册（邮箱、密码、姓名）
func (c *RegisterController) Step1(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	// 验证必填字段
	if req.Email == "" || req.Password == "" || req.FullName == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱、密码、姓名不能为空",
		})
		return
	}

	// 检查邮箱是否已注册
	if service.IsEmailRegistered(req.Email) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "该邮箱已注册",
		})
		return
	}

	// 验证密码强度
	if !service.ValidatePassword(req.Password) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "密码长度必须在8-16位之间",
		})
		return
	}

	// 存储第一步信息到数据库
	user := &service.UserRegistration{
		Email:     req.Email,
		Password:  service.HashPassword(req.Password),
		FullName:  req.FullName,
		Role:      "owner",                        // 默认角色
		Status:    "pending",                      // 待完成状态
		DID:       service.GenerateDID(req.Email), // 生成DID
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 存储到数据库
	if err := service.SaveUserToDatabase(user); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "存储用户信息失败: " + err.Error(),
		})
		return
	}

	// 跳转到角色选择页面
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "基础信息验证成功，请选择角色",
		Data: map[string]interface{}{
			"redirect_url": "/static/register-3.html",
			"email":        req.Email,
		},
	})
}

// POST /api/register/step2
// 第二步：角色选择和额外信息
func (c *RegisterController) Step2(r *ghttp.Request) {
	var req struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		FullName     string `json:"fullName"`
		Role         string `json:"role"`
		Phone        string `json:"phone"`
		PropertyName string `json:"propertyName,omitempty"`
		Occupation   string `json:"occupation,omitempty"`
		Institution  string `json:"institution,omitempty"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	// 验证必填字段
	if req.Email == "" || req.Password == "" || req.FullName == "" || req.Role == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "必填字段不能为空",
		})
		return
	}

	// 根据角色处理
	switch req.Role {
	case "owner":
		// Owner需要进入建筑信息收集流程
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    200,
			Message: "Owner角色，请填写建筑信息",
			Data: map[string]interface{}{
				"redirect_url": "/static/register-3-1.html",
				"role":         req.Role,
			},
		})

	case "property_manager", "institution":
		// 直接完成注册
		registerReq := &model.RegisterRequest{
			Email:        req.Email,
			Password:     req.Password,
			FullName:     req.FullName,
			Role:         req.Role,
			Phone:        req.Phone, // 确保手机号字段正确传递
			PropertyName: req.PropertyName,
			Occupation:   req.Occupation,
			Institution:  req.Institution,
		}

		response, err := service.RegisterUser(registerReq)
		if err != nil {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    500,
				Message: "注册失败",
			})
			return
		}

		if response.Success {
			// 注册成功后直接完成，无需审核，跳转到仪表板界面
			// 获取链码返回的数据
			chaincodeData := map[string]interface{}{
				"assetID":   req.Email,                                // 资产ID（邮箱）
				"status":    "success",                                // 状态
				"message":   "Asset created successfully",             // 消息
				"timestamp": time.Now().Format("2006-01-02 15:04:05"), // 时间戳
				"txID":      fmt.Sprintf("tx_%d", time.Now().Unix()),  // 模拟交易ID
			}

			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    200,
				Message: "注册成功",
				Data: map[string]interface{}{
					"success":       true,
					"redirect_url":  "/static/dashboard.html",
					"message":       "注册成功",
					"RegisterTime":  response.RegisterTime,
					"chaincodeData": chaincodeData,
				},
			})
		} else {
			// 将邮箱已存在的情况视为幂等成功
			if response.ErrorDetails == "EMAIL_ALREADY_EXISTS" {
				r.Response.WriteJson(ghttp.DefaultHandlerResponse{
					Code:    200,
					Message: "账号已存在，视为已注册",
					Data: map[string]interface{}{
						"success":      true,
						"redirect_url": "/static/dashboard.html",
						"message":      "账号已存在，已为您跳转",
					},
				})
				return
			}
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    400,
				Message: response.Message,
			})
		}

	default:
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "无效的用户角色",
		})
	}
}

// POST /api/register/step3
// 第三步：建筑信息收集（Owner角色）
func (c *RegisterController) Step3(r *ghttp.Request) {
	var req struct {
		Email        string `json:"email"`
		Password     string `json:"password"`
		FullName     string `json:"fullName"`
		Role         string `json:"role"`
		Phone        string `json:"phone"`
		BuildingName string `json:"buildingName"`
		BuildingAddr string `json:"buildingAddr"`
		BuildingType string `json:"buildingType"`
	}

	// 记录接收到的原始请求数据
	glog.Info(nil, "🔍 Step3 API 接收到的请求数据:")
	glog.Info(nil, "   - 请求头:", r.Header)
	glog.Info(nil, "   - 请求体:", r.GetBody())
	glog.Info(nil, "   - Content-Type:", r.Header.Get("Content-Type"))

	if err := r.Parse(&req); err != nil {
		glog.Error(nil, "❌ Step3 参数解析失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	// 记录解析后的请求数据
	glog.Info(nil, "✅ Step3 参数解析成功:")
	glog.Info(nil, "   - Email:", req.Email)
	glog.Info(nil, "   - Password:", req.Password)
	glog.Info(nil, "   - FullName:", req.FullName)
	glog.Info(nil, "   - Role:", req.Role)
	glog.Info(nil, "   - Phone:", req.Phone)
	glog.Info(nil, "   - BuildingName:", req.BuildingName)
	glog.Info(nil, "   - BuildingAddr:", req.BuildingAddr)
	glog.Info(nil, "   - BuildingType:", req.BuildingType)

	// 验证必填字段
	glog.Info(nil, "🔍 Step3 开始验证必填字段...")

	if req.Email == "" {
		glog.Error(nil, "❌ Email 字段为空")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Email 字段不能为空",
		})
		return
	}

	if req.Password == "" {
		glog.Error(nil, "❌ Password 字段为空")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Password 字段不能为空",
		})
		return
	}

	if req.FullName == "" {
		glog.Error(nil, "❌ FullName 字段为空")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "FullName 字段不能为空",
		})
		return
	}

	if req.Role == "" {
		glog.Error(nil, "❌ Role 字段为空")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "Role 字段不能为空",
		})
		return
	}

	glog.Info(nil, "✅ 基础字段验证通过")

	// 验证建筑信息
	glog.Info(nil, "🔍 Step3 开始验证建筑信息字段...")

	if req.BuildingName == "" {
		glog.Error(nil, "❌ BuildingName 字段为空")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "建筑名称不能为空",
		})
		return
	}

	if req.BuildingAddr == "" {
		glog.Error(nil, "❌ BuildingAddr 字段为空")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "建筑地址不能为空",
		})
		return
	}

	if req.BuildingType == "" {
		glog.Error(nil, "❌ BuildingType 字段为空")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "建筑类型不能为空",
		})
		return
	}

	glog.Info(nil, "✅ 建筑信息字段验证通过")

	// 完成Owner注册
	registerReq := &model.RegisterRequest{
		Email:        req.Email,
		Password:     req.Password,
		FullName:     req.FullName,
		Role:         req.Role,
		Phone:        req.Phone,
		BuildingName: req.BuildingName,
		BuildingAddr: req.BuildingAddr,
		BuildingType: req.BuildingType,
	}

	response, err := service.RegisterUser(registerReq)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "注册失败",
		})
		return
	}

	if response.Success {
		// 使用真实的链码返回数据
		chaincodeData := response.ChaincodeData
		if chaincodeData == nil {
			// 如果没有链码数据，使用默认值
			chaincodeData = map[string]interface{}{
				"assetID":   req.Email,                                // 资产ID（邮箱）
				"status":    "success",                                // 状态
				"message":   "Asset created successfully",             // 消息
				"timestamp": time.Now().Format("2006-01-02 15:04:05"), // 时间戳
				"txID":      fmt.Sprintf("tx_%d", time.Now().Unix()),  // 模拟交易ID
			}
		}

		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    200,
			Message: "注册成功",
			Data: map[string]interface{}{
				"success":       true,
				"redirect_url":  "/static/dashboard.html",
				"message":       "注册成功",
				"RegisterTime":  response.RegisterTime,
				"chaincodeData": chaincodeData,
			},
		})
	} else {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: response.Message,
		})
	}
}

// POST /api/register/complete
// 完成注册
func (c *RegisterController) Complete(r *ghttp.Request) {
	var req struct {
		Email string `json:"email"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	// 完成用户注册
	err := service.CompleteUserRegistration(req.Email)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "完成注册失败",
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "注册完成",
		Data: map[string]interface{}{
			"redirect_url": "/static/dashboard.html",
		},
	})
}

// GET /api/register/status
// 获取注册状态
func (c *RegisterController) GetStatus(r *ghttp.Request) {
	email := r.Get("email").String()
	did := r.Get("did").String()
	if email == "" && did == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "邮箱或DID至少提供一个"})
		return
	}

	var user *service.UserRegistration
	var exists bool
	if email != "" {
		user, exists = service.GetUserRegistration(email)
	} else {
		user, exists = service.GetUserRegistrationByDID(did)
		if exists {
			email = user.Email
		}
	}
	if !exists {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 404, Message: "用户不存在"})
		return
	}

	// 仅基于MySQL构造链码相关展示数据，避免读取链码失败
	registerTimeStr := user.CreatedAt.Format("2006-01-02 15:04:05")
	chaincodeData := map[string]interface{}{
		"status":    user.Status,
		"timestamp": registerTimeStr,
	}

	resp := map[string]interface{}{
		"Email":         user.Email,
		"FullName":      user.FullName,
		"Role":          user.Role,
		"Phone":         user.Phone,
		"Age":           user.Age,
		"did":           user.DID,
		"DID":           user.DID, // 兼容前端大小写读取
		"building_name": user.BuildingName,
		"building_addr": user.BuildingAddr,
		"building_type": user.BuildingType,
		"property_name": user.PropertyName,
		"occupation":    user.Occupation,
		"institution":   user.Institution,
		"status":        user.Status,
		"assetID":       email,
		"RegisterTime":  registerTimeStr,
		"timestamp":     registerTimeStr,
		"chaincodeData": chaincodeData,
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "获取注册状态成功", Data: resp})
}

// POST /api/register/update
// 更新注册信息
func (c *RegisterController) Update(r *ghttp.Request) {
	var req struct {
		Email        string `json:"email"`
		BuildingName string `json:"buildingName,omitempty"`
		BuildingAddr string `json:"buildingAddr,omitempty"`
		BuildingType string `json:"buildingType,omitempty"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	// 更新用户信息
	updates := make(map[string]interface{})
	if req.BuildingName != "" {
		updates["buildingName"] = req.BuildingName
	}
	if req.BuildingAddr != "" {
		updates["buildingAddr"] = req.BuildingAddr
	}
	if req.BuildingType != "" {
		updates["buildingType"] = req.BuildingType
	}

	err := service.UpdateUserRegistration(req.Email, updates)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "更新失败",
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "更新成功",
	})
}

// POST /api/register/verify
// 用户登录验证
func (c *RegisterController) Verify(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email" v:"required|email"`
		Password string `json:"password" v:"required"`
	}

	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "参数解析失败",
		})
		return
	}

	// 验证用户登录
	user, exists := service.GetUserRegistration(req.Email)
	if !exists {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	// 验证密码
	if !service.ValidatePasswordHash(req.Password, user.Password) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "密码错误",
		})
		return
	}

	// 登录成功，返回用户信息
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "登录成功",
		Data: map[string]interface{}{
			"valid":        true,
			"email":        user.Email,
			"fullName":     user.FullName,
			"role":         user.Role,
			"phone":        user.Phone,
			"age":          user.Age,
			"did":          user.DID,
			"registerTime": user.CreatedAt,
		},
	})
}

// POST /api/register/check-email
// 检查邮箱是否已注册
func (c *RegisterController) CheckEmail(r *ghttp.Request) {
	var req struct {
		Email        string `json:"email"`
		ForceRefresh bool   `json:"force_refresh,omitempty"` // 添加强制刷新参数
	}

	if err := r.Parse(&req); err != nil || req.Email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱地址不能为空",
		})
		return
	}

	// 如果请求强制刷新，先清理缓存和数据库连接
	if req.ForceRefresh {
		// 强制刷新数据库连接
		service.ForceRefreshDatabase()
		// 刷新用户缓存
		service.RefreshUserCache(req.Email)
		glog.New().Infof(r.Context(), "🔄 强制刷新数据库连接和用户缓存: %s", req.Email)
	}

	// 检查邮箱是否已注册
	exists := service.IsEmailRegistered(req.Email)

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "邮箱检查完成",
		Data: map[string]interface{}{
			"exists":    exists,
			"email":     req.Email,
			"refreshed": req.ForceRefresh,
		},
	})
}

// POST /api/register/reset-password
// 重置密码（需客户端已完成验证码校验）
func (c *RegisterController) ResetPassword(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := r.Parse(&req); err != nil || req.Email == "" || req.Password == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "邮箱与新密码不能为空"})
		return
	}
	if !service.ValidatePassword(req.Password) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "密码长度必须在8-16位之间"})
		return
	}
	if err := service.ResetPassword(req.Email, req.Password); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 404, Message: err.Error()})
		return
	}
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "密码重置成功"})
}

// POST /api/register/clear-cache
// 清理指定用户的缓存
func (c *RegisterController) ClearCache(r *ghttp.Request) {
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

	// 清理用户缓存
	service.ClearUserFromCache(req.Email)

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "缓存清理成功",
		Data: map[string]interface{}{
			"email":   req.Email,
			"cleared": true,
		},
	})
}

// GET /api/register/emails
// 获取所有注册成功的邮箱账号
func (c *RegisterController) GetAllEmails(r *ghttp.Request) {
	// 获取所有注册成功的邮箱列表
	response, err := service.GetAllRegisteredEmails()
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取邮箱列表失败",
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取邮箱列表成功",
		Data:    response,
	})
}

// GET /api/register/admin/users
// 获取所有用户（管理员用，只显示已审核通过的用户）
func (c *RegisterController) GetAllUsersForAdmin(r *ghttp.Request) {
	// 获取已审核通过的用户列表
	response, err := service.GetAllUsersForAdmin()
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取所有用户列表失败",
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取所有用户列表成功",
		Data:    response,
	})
}

// GET /api/register/admin/pending-users
// 获取待审核用户（管理员审核用）
func (c *RegisterController) GetPendingUsersForAdmin(r *ghttp.Request) {
	// 获取待审核用户列表
	response, err := service.GetPendingUsersForAdmin()
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取待审核用户列表失败",
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取待审核用户列表成功",
		Data:    response,
	})
}

// POST /api/register/admin/cleanup
// 清理未完成注册的账号
func (c *RegisterController) CleanupIncompleteRegistrations(r *ghttp.Request) {
	// 执行清理操作
	err := service.CleanupIncompleteRegistrations()
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "清理未完成注册账号失败",
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "清理未完成注册账号成功",
	})
}

// GET /api/register/user-detail
// 根据邮箱获取用户详细信息
func (c *RegisterController) GetUserDetail(r *ghttp.Request) {
	// 获取邮箱参数
	email := r.Get("email").String()

	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱地址不能为空",
		})
		return
	}

	// 获取用户详细信息
	response, err := service.GetUserDetailByEmail(email)
	if err != nil {
		if err.Error() == "用户不存在" {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    404,
				Message: "用户不存在",
			})
		} else {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    500,
				Message: "获取用户详细信息失败",
			})
		}
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取用户详细信息成功",
		Data:    response,
	})
}
