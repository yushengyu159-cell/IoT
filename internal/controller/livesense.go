package controller

import (
	"fabric-sdk/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type LiveSenseController struct{}

var LiveSense = new(LiveSenseController)

// CheckBelucksAPIUser 检查用户是否有权限使用LiveSense API（中间件）
// 支持多个用户，只要用户在数据库中的name字段为"belucksapi"
func (c *LiveSenseController) CheckBelucksAPIUser(r *ghttp.Request) {
	// 从请求头、查询参数或Cookie获取用户邮箱
	email := r.Get("email").String()
	if email == "" {
		email = r.GetHeader("X-User-Email")
	}
	if email == "" {
		// 尝试从Cookie获取
		cookieEmail := r.Cookie.Get("user_email")
		if cookieEmail != nil {
			email = cookieEmail.String()
		}
	}
	
	// 如果还是没有，尝试从localStorage（前端传递的请求头）
	if email == "" {
		email = r.GetHeader("X-User-Email-From-LocalStorage")
	}

	if email == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    401,
			Message: "未提供用户邮箱，请先登录",
		})
		r.Exit()
		return
	}

	// 验证用户是否有权限使用LiveSense API
	// 检查数据库中是否存在该用户且名字为belucksapi
	db := service.GetDB()
	if db == nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "数据库连接失败",
		})
		r.Exit()
		return
	}

	var count int64
	result := db.Table("dids").Where("email = ? AND name = ?", email, "belucksapi").Count(&count)
	if result.Error != nil || count == 0 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    403,
			Message: "无权限访问此接口，仅belucksapi用户可使用。请联系管理员开通权限。",
		})
		r.Exit()
		return
	}

	// 将用户邮箱存储到请求上下文中，供后续使用
	r.SetCtxVar("user_email", email)
	r.Middleware.Next()
}

// GetSensorValues 获取传感器数值
// GET /api/livesense/values
func (c *LiveSenseController) GetSensorValues(r *ghttp.Request) {
	ctx := r.Context()

	sensorID := r.Get("id").String()
	from := r.Get("from").String()
	to := r.Get("to").String()
	aggregation := r.Get("aggregation").String()
	group := r.Get("group").String()
	limit := r.Get("limit").String()

	if sensorID == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "传感器ID不能为空",
		})
		return
	}

	values, err := service.LiveSenseAPI.GetSensorValues(ctx, sensorID, from, to, aggregation, group, limit)
	if err != nil {
		g.Log().Error(ctx, "获取传感器数值失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取传感器数值失败: " + err.Error(),
		})
		return
	}

	// 使用GetValuesList方法处理values为对象、数组或null的情况
	valuesList := values.GetValuesList()
	if len(valuesList) == 0 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    200,
			Message: "获取成功，但无数据",
			Data: g.Map{
				"values": nil,
				"note":   "指定时间范围内无传感器数据",
			},
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取成功",
		Data: g.Map{
			"values": valuesList,
		},
	})
}

// GetSensorsUnderContext 获取上下文下的传感器列表
// GET /api/livesense/sensors
func (c *LiveSenseController) GetSensorsUnderContext(r *ghttp.Request) {
	ctx := r.Context()

	contextID := r.Get("contextId").Int()
	contextTypeID := r.Get("contextTypeId").String()
	groupByTypeID := r.Get("groupByTypeId").String()

	if contextID == 0 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "上下文ID不能为空",
		})
		return
	}

	sensors, err := service.LiveSenseAPI.GetSensorsUnderContext(ctx, contextID, contextTypeID, groupByTypeID)
	if err != nil {
		g.Log().Error(ctx, "获取传感器列表失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "获取传感器列表失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "获取成功",
		Data:    sensors,
	})
}

// Authenticate 认证接口（用于测试）
// POST /api/livesense/authenticate
func (c *LiveSenseController) Authenticate(r *ghttp.Request) {
	ctx := r.Context()

	// 手动触发认证
	err := service.LiveSenseAPI.Authenticate(ctx)
	if err != nil {
		g.Log().Error(ctx, "认证失败:", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    500,
			Message: "认证失败: " + err.Error(),
		})
		return
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "认证成功",
		Data: g.Map{
			"authenticated": true,
		},
	})
}

