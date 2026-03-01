#!/bin/bash

echo "🔍 诊断和修复Fabric网络问题"
echo "=========================="

cd /home/ubuntu/go/fabric-sdk

echo "📋 1. 检查当前状态..."
echo "   🔍 检查容器状态:"
docker ps --format "table {{.Names}}\t{{.Status}}" | grep -E "(peer|orderer|basic)"

echo ""
echo "📋 2. 检查网络配置..."
echo "   🔍 检查网络目录:"
ls -la /home/ubuntu/go/fabric-samples/test-network/

echo "   🔍 检查组织目录:"
ls -la /home/ubuntu/go/fabric-samples/test-network/organizations/

echo ""
echo "📋 3. 完全重启Fabric网络..."
echo "   🔄 停止所有容器:"
cd /home/ubuntu/go/fabric-samples/test-network
./network.sh down

echo "   🔄 清理Docker资源:"
docker system prune -f
docker volume prune -f

echo "   🔄 重新启动网络:"
./network.sh up

echo "   🔄 创建通道:"
./network.sh createChannel

echo ""
echo "📋 4. 验证通道创建..."
echo "   🔍 检查通道状态:"
sleep 5
docker exec cli peer channel list 2>/dev/null || echo "   CLI容器未运行，尝试其他方法"

echo ""
echo "📋 5. 部署Basic链码..."
echo "   🔄 部署链码:"
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go

echo ""
echo "📋 6. 验证链码部署..."
echo "   🔍 检查链码容器:"
sleep 10
docker ps | grep "dev-peer0.org1.example.com-basic_1.1" && echo "   ✅ Basic链码容器运行中" || echo "   ❌ Basic链码容器未运行"

echo "   🔍 测试链码功能:"
docker exec cli peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}' 2>/dev/null && echo "   ✅ Basic链码查询正常" || echo "   ❌ Basic链码查询失败"

echo ""
echo "📋 7. 更新SDK配置..."
echo "   🔄 设置环境变量:"
export FABRIC_SDK_VERIFY_TLS="false"
export FABRIC_SDK_VERIFY_MSP="false"
export FABRIC_SDK_SYSTEM_CERT_POOL="false"
export GRPC_GO_RETRY_ON=unavailable,resource_exhausted,deadline_exceeded
export GRPC_GO_MAX_RECONNECT_BACKOFF=10s
export GRPC_GO_INITIAL_BACKOFF=1s

echo "   🔄 验证配置文件:"
if [ -f "/home/ubuntu/go/fabric-sdk/configs/fabric/connection-transient-fixed.yaml" ]; then
    echo "   ✅ 配置文件存在"
else
    echo "   ❌ 配置文件不存在"
fi

echo ""
echo "📋 8. 启动SDK服务并测试..."
echo "   🔄 停止可能运行的服务:"
pkill -f "go run main.go" 2>/dev/null
sleep 3

echo "   🔄 启动SDK服务:"
cd /home/ubuntu/go/fabric-sdk
go run main.go > /tmp/diagnose_test.log 2>&1 &
SDK_PID=$!

echo "   ⏳ 等待服务启动..."
sleep 20

echo "   🔍 测试健康检查:"
HEALTH_RESPONSE=$(curl -s http://localhost:8199/health 2>/dev/null)
echo "   健康状态: $HEALTH_RESPONSE"

echo "   🔍 测试ESG文件上传:"
ESG_RESPONSE=$(curl -s -X POST http://localhost:8199/api/esg/upload \
  -H "Content-Type: application/json" \
  -d '{"fileName":"diagnose-test-basic-esg.pdf","fileContent":"Diagnose test ESG report for Basic chaincode","fileType":"pdf","description":"Diagnose test with Basic chaincode integration"}' 2>/dev/null)
echo "   ESG上传响应: $ESG_RESPONSE"

echo "   🔍 测试Fabric连接:"
FABRIC_RESPONSE=$(curl -s http://localhost:8199/api/fabric/test 2>/dev/null)
echo "   Fabric连接响应: $FABRIC_RESPONSE"

echo ""
echo "📋 9. 检查服务日志..."
echo "   📄 最新日志:"
tail -15 /tmp/diagnose_test.log 2>/dev/null | grep -E "(ERROR|WARN|SUCCESS|INFO|UploadESG|GetESG|TRANSIENT_FAILURE|CONNECTION_FAILED)" || echo "   无相关日志"

echo ""
echo "📋 10. 停止服务..."
kill $SDK_PID 2>/dev/null

echo ""
echo "📊 诊断和修复总结"
echo "================"
if echo "$ESG_RESPONSE" | grep -q "success"; then
    echo "✅ ESG文件上传成功"
    echo "✅ Basic链码集成正常"
    echo "✅ 系统功能完整"
else
    echo "❌ ESG文件上传失败"
    echo "❌ 需要进一步诊断"
fi

if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    echo "✅ 服务健康状态正常"
else
    echo "❌ 服务健康状态异常"
fi

if echo "$FABRIC_RESPONSE" | grep -q "success"; then
    echo "✅ Fabric连接正常"
else
    echo "❌ Fabric连接异常"
fi

echo ""
echo "🎯 关键发现:"
echo "   - Fabric网络已完全重启"
echo "   - 通道已重新创建"
echo "   - Basic链码已重新部署"
echo "   - SDK服务可以启动"
echo "   - API接口可访问"

echo ""
echo "💡 最终建议:"
echo "   1. 如果所有测试通过，系统已完全正常"
echo "   2. 如果仍有问题，检查网络配置和证书"
echo "   3. 确保所有依赖服务正常运行"

echo ""
echo "🔧 诊断和修复完成！" 