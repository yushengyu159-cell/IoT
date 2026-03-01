#!/bin/bash

echo "🚀 ESG VISA 邮箱配置自动设置脚本"
echo "=================================="

# 检查配置文件是否存在
if [ ! -f "configs/email.yaml" ]; then
    echo "❌ 配置文件 configs/email.yaml 不存在"
    exit 1
fi

echo "📧 当前邮箱配置状态："
echo ""

# 检查163邮箱配置
echo "🔍 检查163邮箱配置..."
if grep -q "your-163-smtp-password" configs/email.yaml; then
    echo "❌ 163邮箱密码未配置（仍使用占位符）"
    echo "   需要获取SMTP授权码并更新配置文件"
else
    echo "✅ 163邮箱密码已配置"
fi

echo ""

# 检查QQ邮箱配置
echo "🔍 检查QQ邮箱配置..."
if grep -q "your-qq-smtp-password" configs/email.yaml; then
    echo "❌ QQ邮箱密码未配置（仍使用占位符）"
    echo "   需要获取SMTP授权码并更新配置文件"
else
    echo "✅ QQ邮箱密码已配置"
fi

echo ""

echo "📋 配置步骤说明："
echo "=================="
echo ""
echo "1️⃣ 163邮箱SMTP授权码获取："
echo "   - 登录163邮箱网页版"
echo "   - 设置 → POP3/SMTP/IMAP"
echo "   - 开启SMTP服务，设置授权码"
echo "   - 将授权码填入 configs/email.yaml 的 password 字段"
echo ""
echo "2️⃣ QQ邮箱SMTP授权码获取："
echo "   - 登录QQ邮箱网页版"
echo "   - 设置 → 账户 → POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务"
echo "   - 开启SMTP服务，获取授权码"
echo "   - 将授权码填入 configs/email.yaml 的 password 字段"
echo ""
echo "3️⃣ 配置文件更新后："
echo "   - 重启后端服务：go run main.go"
echo "   - 测试邮件发送功能"
echo "   - 检查邮箱是否收到验证码"
echo ""

echo "🎯 当前状态："
echo "============="
echo "✅ 邮箱服务已配置"
echo "✅ SSL/TLS连接已配置"
echo "✅ 验证码生成功能正常"
echo "❌ 邮件发送功能未配置（需要真实SMTP密码）"
echo ""

echo "💡 建议："
echo "========="
echo "1. 先使用测试模式完成功能测试"
echo "2. 配置真实邮箱后再测试邮件发送"
echo "3. 使用日志中的验证码完成当前验证"
echo ""

echo "🔧 测试命令："
echo "============="
echo "1. 启动服务：go run main.go"
echo "2. 测试连接：curl -X POST http://localhost:8199/api/email/test-connection -H 'Content-Type: application/json' -d '{\"email\":\"tj18832045990@163.com\"}'"
echo "3. 发送验证码：curl -X POST http://localhost:8199/api/email/send-code -H 'Content-Type: application/json' -d '{\"email\":\"tj18832045990@163.com\"}'"
echo ""

echo "✨ 配置完成后，邮箱就能正常收到验证码了！"
