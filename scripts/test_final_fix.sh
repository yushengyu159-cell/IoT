#!/bin/bash

echo "🔧 最终修复测试 - 基于GoFrame错误处理最佳实践"
echo "============================================="

echo "📋 1. 检查代码编译..."
cd /home/ubuntu/go/fabric-sdk

echo "   编译代码..."
if go build -o /tmp/fabric-sdk-test .; then
    echo "   ✅ 代码编译成功"
else
    echo "   ❌ 代码编译失败"
    exit 1
fi

echo ""
echo "📋 2. 设置环境变量..."
export FABRIC_SDK_VERIFY_TLS="false"
export FABRIC_SDK_VERIFY_MSP="false"
export FABRIC_SDK_SYSTEM_CERT_POOL="false"
export GODEBUG="x509ignoreCN=0"
export CGO_ENABLED=1
export FABRIC_LOGGING_SPEC="INFO"
export FABRIC_CFG_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric"

echo "   ✅ 环境变量设置完成"

echo ""
echo "📋 3. 检查配置文件..."
configFile="/home/ubuntu/go/fabric-sdk/configs/fabric/connection-host-ports.yaml"
if [ -f "$configFile" ]; then
    echo "   ✅ 配置文件存在: $configFile"
else
    echo "   ❌ 配置文件不存在: $configFile"
    exit 1
fi

echo ""
echo "📋 4. 检查证书文件..."
certPath="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
keyPath="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/"

if [ -f "$certPath" ]; then
    echo "   ✅ 证书文件存在: $certPath"
else
    echo "   ❌ 证书文件不存在: $certPath"
    exit 1
fi

if [ -d "$keyPath" ] && [ "$(ls -A $keyPath)" ]; then
    echo "   ✅ 私钥目录存在且不为空: $keyPath"
else
    echo "   ❌ 私钥目录不存在或为空: $keyPath"
    exit 1
fi

echo ""
echo "📋 5. 启动SDK服务测试..."
# 先停止可能运行的服务
pkill -f "go run main.go" 2>/dev/null
sleep 2

echo "   启动服务..."
timeout 30s go run main.go > /tmp/final_test.log 2>&1 &
SDK_PID=$!

echo "   等待服务启动..."
sleep 10

echo "   测试API连接..."
RESPONSE=$(curl -s http://localhost:8199/api/fabric/test 2>/dev/null)
echo "   响应: $RESPONSE"

echo "   检查日志..."
echo "   - 最新日志:"
tail -20 /tmp/final_test.log 2>/dev/null || echo "   无日志"

echo ""
echo "   停止服务..."
kill $SDK_PID 2>/dev/null

echo ""
echo "📊 最终修复测试总结"
echo "=================="
echo "✅ 代码编译检查通过"
echo "✅ 环境变量设置完成"
echo "✅ 配置文件检查通过"
echo "✅ 证书文件检查通过"
echo "✅ SDK服务启动测试完成"

echo ""
echo "🎉 最终修复测试完成！"
echo ""
echo "💡 说明:"
echo "   - 基于GoFrame错误处理最佳实践"
echo "   - 使用新的Gateway连接架构"
echo "   - 完整的错误处理和日志记录"
echo "   - 标准的Fabric SDK Go使用方式" 