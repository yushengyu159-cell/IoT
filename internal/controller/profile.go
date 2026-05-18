package controller

import (
	"fabric-sdk/internal/model"
	"fabric-sdk/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ProfileController struct{}

var Profile = new(ProfileController)

// extractEmail gets user email from query params, headers, or cookies
func extractEmail(r *ghttp.Request) string {
	email := r.Get("email").String()
	if email == "" {
		email = r.GetHeader("X-User-Email")
	}
	if email == "" {
		if cookieEmail := r.Cookie.Get("user_email"); cookieEmail != nil {
			email = cookieEmail.String()
		}
	}
	return email
}

// CreateBuildingProfile POST /api/profile/building
func (c *ProfileController) CreateBuildingProfile(r *ghttp.Request) {
	var req model.BuildingProfileRequest
	if err := r.Parse(&req); err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "Invalid parameters"})
		return
	}

	email := req.Email
	if email == "" {
		email = extractEmail(r)
	}
	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "email is required"})
		return
	}
	req.Email = email

	result, err := service.UpsertBuildingProfile(r.Context(), &req)
	if err != nil {
		g.Log().Error(r.Context(), "Save building profile failed:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "Failed to save: " + err.Error()})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code: 200, Message: "Success", Data: result,
	})
}

// GetBuildingProfile GET /api/profile/building?email=X
func (c *ProfileController) GetBuildingProfile(r *ghttp.Request) {
	email := extractEmail(r)
	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "email is required"})
		return
	}

	result, err := service.AutoInitFromRegistration(r.Context(), email)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 404, Message: "Building profile not found: " + err.Error()})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code: 200, Message: "Success", Data: result,
	})
}

// GetBlockchainRecord GET /api/profile/blockchain-record?email=X
func (c *ProfileController) GetBlockchainRecord(r *ghttp.Request) {
	email := extractEmail(r)
	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "email is required"})
		return
	}

	chainData, err := service.GetBlockchainRecord(r.Context(), email)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "Blockchain read failed: " + err.Error()})
		return
	}

	verified, _ := service.VerifyDataIntegrity(r.Context(), email)
	chainData["verified"] = verified

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code: 200, Message: "Success", Data: chainData,
	})
}

// GetCertificate GET /api/profile/certificate?email=X
func (c *ProfileController) GetCertificate(r *ghttp.Request) {
	email := extractEmail(r)
	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "email is required"})
		return
	}

	asset, err := service.GetBuildingProfile(email)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 404, Message: "No building profile found"})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code: 200, Message: "Success",
		Data: map[string]interface{}{
			"buildingName":   asset.BuildingName,
			"address":        asset.BuildingAddr,
			"certifications": asset.Certifications,
			"grade":          asset.Grade,
			"accreditedBy":   asset.AccreditedBy,
			"chainTxId":      asset.ChainTxID,
			"verified":       asset.Verified,
			"createdAt":      asset.CreatedAt,
			"updatedAt":      asset.UpdatedAt,
		},
	})
}
