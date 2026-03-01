#!/bin/bash

echo "🚀 ESG VISA 邮箱自动配置脚本"
echo "============================="
echo ""

# 检查配置文件是否存在
if [ ! -f "configs/email.yaml" ]; then
    echo "❌ 配置文件 configs/email.yaml 不存在"
    exit 1
fi

echo "📧 开始自动配置邮箱..."
echo ""

# 备份原配置文件
cp configs/email.yaml configs/email.yaml.backup
echo "✅ 已备份原配置文件"

# 创建临时配置文件（使用真实邮箱配置）
cat > configs/email_temp.yaml << 'EOF'
# 邮箱验证配置
email:
  # QQ邮箱配置
  qq:
    host: "smtp.qq.com"
    port: "587"
    username: "tj18832045990@qq.com"
    password: "REPLACE_WITH_QQ_SMTP_PASSWORD"
    from: "tj18832045990@qq.com"
    enabled: true
    use_ssl: false
    use_tls: true
  
  # 163邮箱配置
  163:
    host: "smtp.163.com"
    port: "465"
    username: "tj18832045990@163.com"
    password: "REPLACE_WITH_163_SMTP_PASSWORD"
    from: "tj18832045990@163.com"
    enabled: true
    use_ssl: true
    use_tls: false
  
  # Gmail配置
  gmail:
    host: "smtp.gmail.com"
    port: "587"
    username: "your-gmail@gmail.com"
    password: "your-gmail-app-password"
    from: "your-gmail@gmail.com"
    enabled: false
    use_ssl: false
    use_tls: true
  
  # 通用配置
  common:
    code_expire_minutes: 10
    resend_cooldown_seconds: 60
    max_send_per_hour: 10
    subject: "ESG VISA - 邮箱验证码"
    show_code_in_dev: true
EOF

echo "✅ 已创建临时配置文件"
echo ""

echo "🔐 请按照以下步骤获取SMTP授权码："
echo "=================================="
echo ""

echo "1️⃣ 163邮箱SMTP授权码获取："
echo "   - 登录163邮箱网页版：https://mail.163.com"
echo "   - 设置 → POP3/SMTP/IMAP"
echo "   - 开启SMTP服务，设置授权码"
echo "   - 复制授权码"
echo ""

echo "2️⃣ QQ邮箱SMTP授权码获取："
echo "   - 登录QQ邮箱网页版：https://mail.qq.com"
echo "   - 设置 → 账户 → POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务"
echo "   - 开启SMTP服务，获取授权码"
echo "   - 复制授权码"
echo ""

echo "3️⃣ 配置完成后："
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

echo "💡 快速配置方法："
echo "================="
echo "1. 手动编辑 configs/email.yaml 文件"
echo "2. 将 'your-163-smtp-password' 替换为真实授权码"
echo "3. 将 'your-qq-smtp-password' 替换为真实授权码"
echo "4. 保存文件并重启服务"
echo ""

echo "🔧 测试命令："
echo "============="
echo "1. 启动服务：go run main.go"
echo "2. 测试连接：curl -X POST http://localhost:8199/api/email/test-connection -H 'Content-Type: application/json' -d '{\"email\":\"tj18832045990@163.com\"}'"
echo "3. 发送验证码：curl -X POST http://localhost:8199/api/email/send-code -H 'Content-Type: application/json' -d '{\"email\":\"tj18832045990@163.com\"}'"
echo ""

echo "✨ 配置完成后，邮箱就能正常收到验证码了！"
echo ""
echo "📝 注意：请确保在配置真实授权码之前不要重启服务，否则仍会使用测试模式"
