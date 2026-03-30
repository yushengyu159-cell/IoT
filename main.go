package main

import (
	"fmt"
	"fabric-sdk/internal/cmd"
	"fabric-sdk/internal/service"

	"github.com/gogf/gf/v2/os/gctx"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
)

func main() {
	fmt.Println("🚀 启动 fabric-sdk 服务...")
	fmt.Println("🧹 清理未完成注册的账号...")
	if err := service.CleanupIncompleteRegistrations(); err != nil {
		fmt.Printf("⚠️ 清理未完成注册账号失败: %v\n", err)
	} else {
		fmt.Println("✅ 清理未完成注册账号完成")
	}
	cmd.Main.Run(gctx.New())
}
