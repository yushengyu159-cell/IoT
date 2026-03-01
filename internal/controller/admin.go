package controller

import (
	"fabric-sdk/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

type AdminController struct{}

var Admin = new(AdminController)

// POST /api/admin/review-approve
func (c *AdminController) ReviewApprove(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email"`
		Reason   string `json:"reason"`
		Language string `json:"language,omitempty"`
	}
	if err := r.Parse(&req); err != nil || req.Email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "邮箱不能为空"})
		return
	}
	if err := service.AdminApproveUserWithLanguage(req.Email, req.Reason, req.Language); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: err.Error()})
		return
	}
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "审核通过"})
}

// POST /api/admin/review-reject
func (c *AdminController) ReviewReject(r *ghttp.Request) {
	var req struct {
		Email    string `json:"email"`
		Reason   string `json:"reason"`
		Language string `json:"language,omitempty"`
	}
	if err := r.Parse(&req); err != nil || req.Email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "邮箱不能为空"})
		return
	}
	if err := service.AdminRejectUserWithLanguage(req.Email, req.Reason, req.Language); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: err.Error()})
		return
	}
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "已拒绝"})
}
