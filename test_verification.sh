#!/bin/bash

echo "🧪 ESG VISA 邮箱验证码功能测试脚本"
echo "=================================="
echo ""

# 测试邮箱
EMAIL="tj18832045990@163.com"

echo "📧 测试邮箱: $EMAIL"
echo ""

echo "1️⃣ 发送验证码..."
echo "=================="
RESPONSE=$(curl -s -X POST http://localhost:8199/api/email/send-code \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$EMAIL\"}")

echo "响应: $RESPONSE"
echo ""

# 提取验证码（如果存在）
CODE=$(echo $RESPONSE | grep -o '"code":"[^"]*"' | cut -d'"' -f4)
if [ -n "$CODE" ]; then
    echo "✅ 验证码: $CODE"
    echo ""
    
    echo "2️⃣ 验证验证码..."
    echo "=================="
    VERIFY_RESPONSE=$(curl -s -X POST http://localhost:8199/api/email/verify-code \
      -H 'Content-Type: application/json' \
      -d "{\"email\":\"$EMAIL\",\"code\":\"$CODE\"}")
    
    echo "验证响应: $VERIFY_RESPONSE"
    echo ""
    
    echo "3️⃣ 检查验证状态..."
    echo "=================="
    STATUS_RESPONSE=$(curl -s -X GET "http://localhost:8199/api/email/status?email=$EMAIL")
    echo "状态响应: $STATUS_RESPONSE"
    
else
    echo "❌ 未获取到验证码"
fi

echo ""
echo "🎯 测试完成！"
echo "============="
echo "如果验证成功，应该看到 'valid\": true' 的响应"
echo "如果验证失败，请检查日志和配置"
