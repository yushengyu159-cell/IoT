package logic

import (
	"context"
	"fabric-sdk/internal/service"
)

type HelloLogic struct{}

var Hello = new(HelloLogic)

// GetSystemInfo 获取系统信息
func (l *HelloLogic) GetSystemInfo(ctx context.Context) (map[string]interface{}, error) {
	return service.Hello.GetInfo(ctx)
}
