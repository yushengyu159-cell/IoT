#!/bin/bash

echo "🚀 ESG VISA 邮箱快速配置脚本"
echo "============================="
echo ""

# 检查配置文件是否存在
if [ ! -f "configs/email.yaml" ]; then
    echo "❌ 配置文件 configs/email.yaml 不存在"
    exit 1
fi

echo "📧 请按照提示输入邮箱SMTP授权码："
echo ""

# 获取163邮箱SMTP授权码
echo "🔐 163邮箱SMTP授权码配置"
echo "邮箱地址: tj18832045990@163.com"
read -p "请输入163邮箱SMTP授权码: " smtp_163

if [ -n "$smtp_163" ]; then
    echo "✅ 163邮箱SMTP授权码已输入"
    
    # 更新配置文件
    sed -i "s/your-163-smtp-password/$smtp_163/g" configs/email.yaml
    echo "✅ 163邮箱配置已更新"
else
    echo "⚠️  跳过163邮箱配置"
fi

echo ""

# 获取QQ邮箱SMTP授权码
echo "🔐 QQ邮箱SMTP授权码配置"
echo "邮箱地址: tj18832045990@qq.com"
read -p "请输入QQ邮箱SMTP授权码: " smtp_qq

if [ -n "$smtp_qq" ]; then
    echo "✅ QQ邮箱SMTP授权码已输入"
    
    # 更新配置文件
    sed -i "s/your-qq-smtp-password/$smtp_qq/g" configs/email.yaml
    echo "✅ QQ邮箱配置已更新"
else
    echo "⚠️  跳过QQ邮箱配置"
fi

echo ""

echo "🎯 配置完成！"
echo "============="

# 检查配置状态
echo "📋 当前配置状态："
if grep -q "your-163-smtp-password" configs/email.yaml; then
    echo "❌ 163邮箱密码未配置"
else
    echo "✅ 163邮箱密码已配置"
fi

if grep -q "your-qq-smtp-password" configs/email.yaml; then
    echo "❌ QQ邮箱密码未配置"
else
    echo "✅ QQ邮箱密码已配置"
fi

echo ""

echo "🔧 下一步操作："
echo "==============="
echo "1. 重启后端服务：go run main.go"
echo "2. 测试邮件发送功能"
echo "3. 检查邮箱是否收到验证码"
echo ""

echo "💡 提示："
echo "========="
echo "- 如果不知道如何获取SMTP授权码，请查看 configs/email_real_template.yaml 文件"
echo "- 配置完成后，系统将自动发送真实邮件到您的邮箱"
echo "- 验证码将不再显示在日志中，而是直接发送到邮箱"
echo ""

echo "✨ 配置完成！现在可以测试真实邮件发送功能了！"
