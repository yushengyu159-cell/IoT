#!/bin/bash

echo "🔍 精确诊断peer连接问题"
echo "======================"

echo "📋 1. 详细检查peer容器状态..."
echo "   🔍 检查peer容器进程:"
docker exec peer0.org1.example.com ps aux 2>/dev/null | grep peer || echo "   无法检查peer进程"

echo "   🔍 检查peer容器网络:"
docker inspect peer0.org1.example.com | grep -A 10 "NetworkSettings" | head -10

echo "   🔍 检查peer端口映射:"
docker port peer0.org1.example.com

echo ""
echo "📋 2. 精确测试peer端口连接..."
echo "   🔍 使用nc测试端口:"
if command -v nc >/dev/null 2>&1; then
    echo "   测试端口7051:"
    timeout 5 nc -zv localhost 7051 2>&1 && echo "   ✅ nc测试成功" || echo "   ❌ nc测试失败"
else
    echo "   ⚠️ nc命令不可用"
fi

echo "   🔍 使用curl测试端口:"
timeout 5 curl -s telnet://localhost:7051 2>&1 | head -1 && echo "   ✅ curl测试成功" || echo "   ❌ curl测试失败"

echo "   🔍 检查端口监听状态:"
ss -tlnp 2>/dev/null | grep ":7051" || echo "   端口7051未监听"

echo ""
echo "📋 3. 检查peer容器日志..."
echo "   📄 最近20行peer日志:"
docker logs --tail 20 peer0.org1.example.com 2>/dev/null | grep -E "(ERROR|WARN|FATAL|gRPC|grpc|listening|started)" || echo "   无相关日志"

echo ""
echo "📋 4. 检查peer配置文件..."
echo "   🔍 检查peer配置目录:"
ls -la /home/ubuntu/go/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/ 2>/dev/null || echo "   peer配置目录不存在"

echo "   🔍 检查peer配置文件:"
find /home/ubuntu/go/fabric-samples -name "core.yaml" -type f 2>/dev/null | head -1 | xargs ls -la 2>/dev/null || echo "   未找到core.yaml"

echo ""
echo "📋 5. 重启peer容器并等待完全启动..."
echo "   🔄 停止peer容器:"
docker stop peer0.org1.example.com
sleep 5

echo "   🔄 启动peer容器:"
docker start peer0.org1.example.com
echo "   等待peer完全启动..."

# 等待peer完全启动
for i in {1..30}; do
    if docker logs --tail 5 peer0.org1.example.com 2>/dev/null | grep -q "Starting peer"; then
        echo "   ✅ peer正在启动中... (第${i}次检查)"
    elif docker logs --tail 5 peer0.org1.example.com 2>/dev/null | grep -q "Started peer"; then
        echo "   ✅ peer已完全启动"
        break
    else
        echo "   ⏳ 等待peer启动... (第${i}次检查)"
    fi
    sleep 2
done

echo ""
echo "📋 6. 验证peer连接..."
echo "   🔍 检查peer容器状态:"
docker ps | grep "peer0.org1.example.com"

echo "   🔍 检查端口监听:"
ss -tlnp 2>/dev/null | grep ":7051" && echo "   ✅ 端口7051正在监听" || echo "   ❌ 端口7051未监听"

echo "   🔍 测试gRPC连接:"
if command -v grpcurl >/dev/null 2>&1; then
    echo "   使用grpcurl测试:"
    timeout 10 grpcurl -plaintext localhost:7051 list 2>/dev/null && echo "   ✅ gRPC连接成功" || echo "   ❌ gRPC连接失败"
else
    echo "   使用telnet测试:"
    timeout 5 bash -c "</dev/tcp/localhost/7051" 2>/dev/null && echo "   ✅ TCP连接成功" || echo "   ❌ TCP连接失败"
fi

echo ""
echo "📋 7. 检查网络配置..."
echo "   🔍 检查Docker网络:"
docker network ls | grep fabric || echo "   未找到fabric网络"

echo "   🔍 检查容器网络配置:"
docker inspect peer0.org1.example.com | grep -A 5 "IPAddress" || echo "   无法获取IP地址"

echo ""
echo "📊 peer连接诊断总结"
echo "=================="
echo "✅ 检查了peer容器状态"
echo "✅ 测试了端口连接"
echo "✅ 分析了peer日志"
echo "✅ 检查了配置文件"
echo "✅ 重启了peer容器"
echo "✅ 验证了连接状态"

echo ""
echo "🎯 关键发现:"
echo "   - telnet显示连接但立即断开"
echo "   - 可能是gRPC服务未完全启动"
echo "   - 需要等待peer完全启动"

echo ""
echo "💡 建议:"
echo "   1. 确保peer容器完全启动"
echo "   2. 检查gRPC服务状态"
echo "   3. 验证网络配置"
echo "   4. 使用正确的连接方式" 