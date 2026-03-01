#!/bin/bash

echo "🔧 修复peer网络问题"
echo "=================="

echo "📋 1. 检查当前peer状态..."
echo "   🔍 检查peer容器状态:"
docker ps | grep "peer0.org1.example.com" || echo "   ❌ peer容器未运行"

echo "   🔍 检查peer进程:"
docker exec peer0.org1.example.com ps aux 2>/dev/null | grep peer || echo "   ❌ 无法检查peer进程"

echo "   🔍 检查peer端口:"
netstat -tlnp 2>/dev/null | grep ":7051" || echo "   ❌ 端口7051未监听"

echo ""
echo "📋 2. 检查Fabric网络状态..."
echo "   🔍 检查Fabric网络:"
cd /home/ubuntu/go/fabric-samples/test-network
./network.sh status 2>/dev/null || echo "   ❌ 无法检查网络状态"

echo "   🔍 检查网络配置:"
ls -la organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/ 2>/dev/null || echo "   ❌ peer配置目录不存在"

echo ""
echo "📋 3. 重启Fabric网络..."
echo "   🔄 停止网络:"
cd /home/ubuntu/go/fabric-samples/test-network
./network.sh down 2>/dev/null || echo "   ⚠️ 网络停止失败"

echo "   🔄 清理容器:"
docker system prune -f 2>/dev/null || echo "   ⚠️ 清理失败"

echo "   🔄 启动网络:"
./network.sh up -ca 2>/dev/null || echo "   ❌ 网络启动失败"

echo "   🔄 创建通道:"
./network.sh createChannel -c mychannel 2>/dev/null || echo "   ❌ 通道创建失败"

echo "   🔄 部署链码:"
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go 2>/dev/null || echo "   ❌ 链码部署失败"

echo ""
echo "📋 4. 等待peer完全启动..."
echo "   等待peer启动完成..."
for i in {1..60}; do
    if docker logs --tail 10 peer0.org1.example.com 2>/dev/null | grep -q "Started peer"; then
        echo "   ✅ peer已完全启动 (第${i}次检查)"
        break
    elif docker logs --tail 10 peer0.org1.example.com 2>/dev/null | grep -q "Starting peer"; then
        echo "   ⏳ peer正在启动中... (第${i}次检查)"
    else
        echo "   ⏳ 等待peer启动... (第${i}次检查)"
    fi
    sleep 3
done

echo ""
echo "📋 5. 验证peer连接..."
echo "   🔍 检查peer容器状态:"
docker ps | grep "peer0.org1.example.com"

echo "   🔍 检查端口监听:"
netstat -tlnp 2>/dev/null | grep ":7051" && echo "   ✅ 端口7051正在监听" || echo "   ❌ 端口7051未监听"

echo "   🔍 测试TCP连接:"
timeout 5 bash -c "</dev/tcp/localhost/7051" 2>/dev/null && echo "   ✅ TCP连接成功" || echo "   ❌ TCP连接失败"

echo "   🔍 测试gRPC连接:"
if command -v grpcurl >/dev/null 2>&1; then
    timeout 10 grpcurl -plaintext localhost:7051 list 2>/dev/null && echo "   ✅ gRPC连接成功" || echo "   ❌ gRPC连接失败"
else
    echo "   ⚠️ grpcurl不可用，跳过gRPC测试"
fi

echo ""
echo "📋 6. 检查ESG链码..."
echo "   🔍 检查ESG链码状态:"
docker ps | grep "dev-peer0.org1.example.com-esg" || echo "   ❌ ESG链码容器未运行"

echo "   🔍 部署ESG链码:"
cd /home/ubuntu/go/fabric-samples/test-network
./network.sh deployCC -ccn esg -ccp /home/ubuntu/go/fabric-sdk/chaincode -ccl go 2>/dev/null || echo "   ❌ ESG链码部署失败"

echo ""
echo "📋 7. 最终验证..."
echo "   🔍 检查所有容器:"
docker ps | grep -E "(peer|orderer|ca)" | head -10

echo "   🔍 检查网络状态:"
./network.sh status 2>/dev/null | head -20

echo ""
echo "📊 peer网络修复总结"
echo "=================="
echo "✅ 检查了peer状态"
echo "✅ 重启了Fabric网络"
echo "✅ 等待了peer启动"
echo "✅ 验证了peer连接"
echo "✅ 检查了ESG链码"

echo ""
echo "🎯 修复结果:"
if docker ps | grep -q "peer0.org1.example.com" && netstat -tlnp 2>/dev/null | grep -q ":7051"; then
    echo "   ✅ peer网络已修复"
    echo "   ✅ peer容器正常运行"
    echo "   ✅ 端口7051正在监听"
else
    echo "   ❌ peer网络修复失败"
    echo "   ❌ 需要进一步诊断"
fi

echo ""
echo "💡 建议:"
echo "   1. 如果修复成功，重新测试ESG功能"
echo "   2. 如果仍有问题，检查Docker和网络配置"
echo "   3. 确保Fabric网络配置正确" 