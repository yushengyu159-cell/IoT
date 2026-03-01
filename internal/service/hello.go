package service

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

type HelloService struct{}

var Hello = new(HelloService)

// GetInfo 获取系统信息
func (s *HelloService) GetInfo(ctx context.Context) (map[string]interface{}, error) {
	info := map[string]interface{}{
		"framework": "GoFrame v2",
		"project":   "fabric-sdk",
		"version":   "1.0.0",
		"goVersion": "1.22.2",
		"platform":  "linux/amd64",
	}

	g.Log().Info(ctx, "获取系统信息成功")
	return info, nil
}
