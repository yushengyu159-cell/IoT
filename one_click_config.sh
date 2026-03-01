#!/bin/bash

echo "🚀 ESG VISA 邮箱一键配置脚本"
echo "============================="
echo ""

# 检查配置文件是否存在
if [ ! -f "configs/email.yaml" ]; then
    echo "❌ 配置文件 configs/email.yaml 不存在"
    exit 1
fi

echo "📧 开始配置邮箱SMTP授权码..."
echo ""

# 备份原配置文件
cp configs/email.yaml configs/email.yaml.backup.$(date +%Y%m%d_%H%M%S)
echo "✅ 已备份原配置文件"

# 获取163邮箱SMTP授权码
echo "🔐 163邮箱SMTP授权码配置"
echo "邮箱地址: tj18832045990@163.com"
echo "获取步骤：登录163邮箱 → 设置 → POP3/SMTP/IMAP → 开启SMTP服务 → 设置授权码"
echo ""
read -p "请输入163邮箱SMTP授权码: " smtp_163

if [ -n "$smtp_163" ]; then
    echo "✅ 163邮箱SMTP授权码已输入"
    
    # 更新配置文件
    sed -i "s/REPLACE_WITH_163_SMTP_PASSWORD/$smtp_163/g" configs/email.yaml
    echo "✅ 163邮箱配置已更新"
else
    echo "⚠️  跳过163邮箱配置"
fi

echo ""

# 获取QQ邮箱SMTP授权码
echo "🔐 QQ邮箱SMTP授权码配置"
echo "邮箱地址: tj18832045990@qq.com"
echo "获取步骤：登录QQ邮箱 → 设置 → 账户 → POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务 → 开启SMTP服务 → 获取授权码"
echo ""
read -p "请输入QQ邮箱SMTP授权码: " smtp_qq

if [ -n "$smtp_qq" ]; then
    echo "✅ QQ邮箱SMTP授权码已输入"
    
    # 更新配置文件
    sed -i "s/REPLACE_WITH_QQ_SMTP_PASSWORD/$smtp_qq/g" configs/email.yaml
    echo "✅ QQ邮箱配置已更新"
else
    echo "⚠️  跳过QQ邮箱配置"
fi

echo ""

echo "🎯 配置完成！"
echo "============="

# 检查配置状态
echo "📋 当前配置状态："
if grep -q "REPLACE_WITH_163_SMTP_PASSWORD" configs/email.yaml; then
    echo "❌ 163邮箱密码未配置"
else
    echo "✅ 163邮箱密码已配置"
fi

if grep -q "REPLACE_WITH_QQ_SMTP_PASSWORD" configs/email.yaml; then
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

echo "💡 测试命令："
echo "============="
echo "1. 启动服务：go run main.go"
echo "2. 测试连接：curl -X POST http://localhost:8199/api/email/test-connection -H 'Content-Type: application/json' -d '{\"email\":\"tj18832045990@163.com\"}'"
echo "3. 发送验证码：curl -X POST http://localhost:8199/api/email/send-code -H 'Content-Type: application/json' -d '{\"email\":\"tj18832045990@163.com\"}'"
echo ""

echo "✨ 配置完成！现在可以测试真实邮件发送功能了！"
echo ""
echo "📝 注意：配置完成后，系统将自动发送真实邮件到邮箱，验证码不再显示在日志中"
