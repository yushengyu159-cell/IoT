#!/bin/bash

# 邮箱验证功能测试脚本
# 使用方法: ./scripts/test_email.sh

BASE_URL="http://localhost:8199"
TEST_EMAIL="test@example.com"

echo "🧪 开始测试邮箱验证功能..."
echo "=================================="

# 1. 发送验证码
echo "📧 1. 发送验证码到: $TEST_EMAIL"
SEND_RESPONSE=$(curl -s -X POST "$BASE_URL/api/email/send-code" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\"}")

echo "发送响应: $SEND_RESPONSE"
echo ""

# 2. 获取验证码状态
echo "📊 2. 获取验证码状态"
STATUS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/email/status?email=$TEST_EMAIL")
echo "状态响应: $STATUS_RESPONSE"
echo ""

# 3. 验证验证码（使用错误的验证码）
echo "❌ 3. 测试错误验证码"
WRONG_CODE="000000"
VERIFY_WRONG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/email/verify-code" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"code\":\"$WRONG_CODE\"}")

echo "错误验证码响应: $VERIFY_WRONG_RESPONSE"
echo ""

# 4. 重新发送验证码
echo "🔄 4. 重新发送验证码"
RESEND_RESPONSE=$(curl -s -X POST "$BASE_URL/api/email/resend-code" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\"}")

echo "重新发送响应: $RESEND_RESPONSE"
echo ""

echo "✅ 邮箱验证功能测试完成！"
echo "=================================="
echo ""
echo "📝 注意事项:"
echo "1. 确保后端服务已启动 (go run main.go)"
echo "2. 检查邮箱配置文件 configs/email.yaml"
echo "3. 生产环境需要配置真实的邮箱和SMTP密码"
echo "4. 验证码有效期为10分钟"
echo "5. 开发环境会显示验证码，生产环境应隐藏"
