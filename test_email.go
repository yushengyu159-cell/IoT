//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"fabric-sdk/internal/service"
)

func main() {
	fmt.Println("🧪 测试邮箱配置和发送功能...")
	
	// 测试邮箱配置加载
	fmt.Println("\n1. 测试邮箱配置加载...")
	
	// 测试QQ邮箱配置
	qqConfig, err := service.GetEmailConfig("test@qq.com")
	if err != nil {
		fmt.Printf("❌ QQ邮箱配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ QQ邮箱配置成功: %s:%s\n", qqConfig.Host, qqConfig.Port)
	}
	
	// 测试163邮箱配置
	config163, err := service.GetEmailConfig("test@163.com")
	if err != nil {
		fmt.Printf("❌ 163邮箱配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ 163邮箱配置成功: %s:%s\n", config163.Host, config163.Port)
	}
	
	// 测试Gmail配置
	gmailConfig, err := service.GetEmailConfig("test@gmail.com")
	if err != nil {
		fmt.Printf("❌ Gmail配置失败: %v\n", err)
	} else {
		fmt.Printf("✅ Gmail配置成功: %s:%s\n", gmailConfig.Host, gmailConfig.Port)
	}
	
	fmt.Println("\n2. 测试验证码生成...")
	code := service.GenerateVerificationCode()
	fmt.Printf("✅ 生成的验证码: %s\n", code)
	
	fmt.Println("\n3. 测试邮件发送（需要真实邮箱配置）...")
	fmt.Println("⚠️  注意：邮件发送测试需要配置真实的邮箱和SMTP密码")
	fmt.Println("   请先配置 configs/email.yaml 中的真实邮箱信息")
	
	// 这里可以添加实际的邮件发送测试
	// 但需要真实的邮箱配置
}
