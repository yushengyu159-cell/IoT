package controller

import (
	"fabric-sdk/internal/service"
	"fabric-sdk/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type FailedRegistrationController struct{}

// AddFailedRegistration 添加失败注册记录
func (c *FailedRegistrationController) AddFailedRegistration(r *ghttp.Request) {
	var req model.FailedRegistrationRequest
	if err := r.Parse(&req); err != nil {
		g.Log().Error(nil, "❌ 解析请求参数失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "请求参数错误",
		})
		return
	}

	// 添加失败注册记录
	err := service.AddFailedRegistration(req.Email, "", req.Reason)
	if err != nil {
		g.Log().Error(nil, "❌ 添加失败注册记录失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "添加失败注册记录失败",
		})
		return
	}

	g.Log().Info(nil, "✅ 添加失败注册记录成功:", req.Email)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "添加失败注册记录成功",
	})
}

// GetFailedRegistration 获取失败注册记录
func (c *FailedRegistrationController) GetFailedRegistration(r *ghttp.Request) {
	email := r.Get("email").String()
	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱参数不能为空",
		})
		return
	}

	failedReg, err := service.GetFailedRegistration(email)
	if err != nil {
		if err.Error() == "record not found" {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    404,
				Message: "未找到失败注册记录",
			})
			return
		}
		g.Log().Error(nil, "❌ 获取失败注册记录失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取失败注册记录失败",
		})
		return
	}

	response := model.FailedRegistrationResponse{
		ID:          failedReg.ID,
		Email:       failedReg.Email,
		DID:         failedReg.DID,
		Reason:      failedReg.Reason,
		Status:      failedReg.Status,
		FailedAt:    failedReg.FailedAt,
		RetryCount:  failedReg.RetryCount,
		LastRetryAt: failedReg.LastRetryAt,
		CreatedAt:   failedReg.CreatedAt,
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取失败注册记录成功",
		Data:    response,
	})
}

// GetAllFailedRegistrations 获取所有失败注册记录
func (c *FailedRegistrationController) GetAllFailedRegistrations(r *ghttp.Request) {
	failedRegs, err := service.GetAllFailedRegistrations()
	if err != nil {
		g.Log().Error(nil, "❌ 获取失败注册记录列表失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取失败注册记录列表失败",
		})
		return
	}

	response := model.FailedRegistrationListResponse{
		Total: len(failedRegs),
		List:  failedRegs,
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取失败注册记录列表成功",
		Data:    response,
	})
}

// UpdateFailedRegistrationStatus 更新失败注册状态
func (c *FailedRegistrationController) UpdateFailedRegistrationStatus(r *ghttp.Request) {
	email := r.Get("email").String()
	status := r.Get("status").String()
	
	if email == "" || status == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱和状态参数不能为空",
		})
		return
	}

	err := service.UpdateFailedRegistrationStatus(email, status)
	if err != nil {
		g.Log().Error(nil, "❌ 更新失败注册状态失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "更新失败注册状态失败",
		})
		return
	}

	g.Log().Info(nil, "✅ 更新失败注册状态成功:", email, "状态:", status)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "更新失败注册状态成功",
	})
}

// DeleteFailedRegistration 删除失败注册记录
func (c *FailedRegistrationController) DeleteFailedRegistration(r *ghttp.Request) {
	email := r.Get("email").String()
	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "邮箱参数不能为空",
		})
		return
	}

	err := service.DeleteFailedRegistration(email)
	if err != nil {
		g.Log().Error(nil, "❌ 删除失败注册记录失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "删除失败注册记录失败",
		})
		return
	}

	g.Log().Info(nil, "✅ 删除失败注册记录成功:", email)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "删除失败注册记录成功",
	})
}

// GetFailedRegistrationStats 获取失败注册统计
func (c *FailedRegistrationController) GetFailedRegistrationStats(r *ghttp.Request) {
	stats, err := service.GetFailedRegistrationStats()
	if err != nil {
		g.Log().Error(nil, "❌ 获取失败注册统计失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取失败注册统计失败",
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取失败注册统计成功",
		Data:    stats,
	})
}