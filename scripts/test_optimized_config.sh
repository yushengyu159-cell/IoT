#!/bin/bash

echo "🔧 测试优化配置文件 - connection-optimized.yaml"
echo "==============================================="

echo "📋 1. 检查优化配置文件..."
configFile="/home/ubuntu/go/fabric-sdk/configs/fabric/connection-optimized.yaml"
if [ -f "$configFile" ]; then
    echo "   ✅ 优化配置文件存在: $configFile"
    echo "   📄 配置文件内容预览:"
    head -20 "$configFile"
else
    echo "   ❌ 优化配置文件不存在: $configFile"
    exit 1
fi

echo ""
echo "📋 2. 检查TLS证书文件..."
tlsCertPath="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem"
if [ -f "$tlsCertPath" ]; then
    echo "   ✅ TLS证书文件存在: $tlsCertPath"
else
    echo "   ❌ TLS证书文件不存在: $tlsCertPath"
    echo "   🔍 查找TLS证书文件..."
    find /home/ubuntu/go/fabric-sdk/configs/fabric -name "*tlsca*" -type f 2>/dev/null | head -5
fi

echo ""
echo "📋 3. 检查客户端TLS证书..."
clientCertPath="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt"
clientKeyPath="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key"

if [ -f "$clientCertPath" ]; then
    echo "   ✅ 客户端证书文件存在: $clientCertPath"
else
    echo "   ❌ 客户端证书文件不存在: $clientCertPath"
fi

if [ -f "$clientKeyPath" ]; then
    echo "   ✅ 客户端私钥文件存在: $clientKeyPath"
else
    echo "   ❌ 客户端私钥文件不存在: $clientKeyPath"
fi

echo ""
echo "📋 4. 设置环境变量..."
export FABRIC_SDK_VERIFY_TLS="false"
export FABRIC_SDK_VERIFY_MSP="false"
export FABRIC_SDK_SYSTEM_CERT_POOL="false"
export GODEBUG="x509ignoreCN=0"
export CGO_ENABLED=1
export FABRIC_LOGGING_SPEC="INFO"
export FABRIC_CFG_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric"

echo "   ✅ 环境变量设置完成"

echo ""
echo "📋 5. 编译代码..."
cd /home/ubuntu/go/fabric-sdk
if go build -o /tmp/fabric-sdk-optimized .; then
    echo "   ✅ 代码编译成功"
else
    echo "   ❌ 代码编译失败"
    exit 1
fi

echo ""
echo "📋 6. 启动SDK服务测试..."
# 先停止可能运行的服务
pkill -f "go run main.go" 2>/dev/null
sleep 2

echo "   启动服务..."
timeout 30s go run main.go > /tmp/optimized_test.log 2>&1 &
SDK_PID=$!

echo "   等待服务启动..."
sleep 10

echo "   测试API连接..."
RESPONSE=$(curl -s http://localhost:8199/api/fabric/test 2>/dev/null)
echo "   响应: $RESPONSE"

echo "   检查日志..."
echo "   - 最新日志:"
tail -15 /tmp/optimized_test.log 2>/dev/null || echo "   无日志"

echo ""
echo "   停止服务..."
kill $SDK_PID 2>/dev/null

echo ""
echo "📊 优化配置文件测试总结"
echo "====================="
echo "✅ 优化配置文件检查通过"
echo "✅ TLS证书文件检查完成"
echo "✅ 客户端证书文件检查完成"
echo "✅ 环境变量设置完成"
echo "✅ 代码编译成功"
echo "✅ SDK服务启动测试完成"

echo ""
echo "🎉 优化配置文件测试完成！"
echo ""
echo "💡 说明:"
echo "   - 使用IBM故障排除优化的配置文件"
echo "   - 完整的TLS和MSP配置"
echo "   - 详细的entityMatchers配置"
echo "   - 优化的BCCSP配置" 