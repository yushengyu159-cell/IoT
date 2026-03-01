package controller

import (
	"context"
	"fabric-sdk/internal/logic"

	"github.com/gogf/gf/v2/net/ghttp"
)

type HelloController struct{}

var Hello = new(HelloController)

// Index 首页
func (c *HelloController) Index(r *ghttp.Request) {
	ctx := context.Background()

	info, err := logic.Hello.GetSystemInfo(ctx)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    1,
			Message: "获取系统信息失败",
			Data:    nil,
		})
		return
	}

	// 添加系统状态信息
	info["system_status"] = map[string]interface{}{
		"status":  "success",
		"message": "系统运行正常",
		"time":    "2025-07-20 15:24:22",
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    0,
		Message: "Hello Fabric SDK with GoFrame!",
		Data:    info,
	})
}
