#!/bin/bash

echo "🔍 检查Fabric网络状态"
echo "=================="

echo "📋 1. 检查容器状态..."
echo "   🔍 所有Fabric容器:"
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" | grep -E "(peer|orderer|ca)" || echo "   无Fabric容器运行"

echo ""
echo "📋 2. 检查端口监听..."
echo "   🔍 端口7051 (peer0.org1):"
ss -tlnp 2>/dev/null | grep ":7051" && echo "   ✅ peer0.org1端口正常" || echo "   ❌ peer0.org1端口未监听"

echo "   🔍 端口9051 (peer0.org2):"
ss -tlnp 2>/dev/null | grep ":9051" && echo "   ✅ peer0.org2端口正常" || echo "   ❌ peer0.org2端口未监听"

echo "   🔍 端口7050 (orderer):"
ss -tlnp 2>/dev/null | grep ":7050" && echo "   ✅ orderer端口正常" || echo "   ❌ orderer端口未监听"

echo ""
echo "📋 3. 检查网络状态..."
cd /home/ubuntu/go/fabric-samples/test-network
if [ -f "network.sh" ]; then
    echo "   🔍 网络状态:"
    ./network.sh status 2>/dev/null | head -20 || echo "   无法获取网络状态"
else
    echo "   ❌ network.sh不存在"
fi

echo ""
echo "📋 4. 测试peer连接..."
echo "   🔍 测试peer0.org1连接:"
timeout 5 bash -c "</dev/tcp/localhost/7051" 2>/dev/null && echo "   ✅ peer0.org1连接成功" || echo "   ❌ peer0.org1连接失败"

echo "   🔍 测试peer0.org2连接:"
timeout 5 bash -c "</dev/tcp/localhost/9051" 2>/dev/null && echo "   ✅ peer0.org2连接成功" || echo "   ❌ peer0.org2连接失败"

echo ""
echo "📊 网络状态总结"
echo "=============="
PEER1_UP=$(docker ps | grep -c "peer0.org1.example.com")
PEER2_UP=$(docker ps | grep -c "peer0.org2.example.com")
ORDERER_UP=$(docker ps | grep -c "orderer.example.com")
CA_UP=$(docker ps | grep -c "ca_")

echo "   peer0.org1: $([ $PEER1_UP -gt 0 ] && echo "✅ 运行中" || echo "❌ 未运行")"
echo "   peer0.org2: $([ $PEER2_UP -gt 0 ] && echo "✅ 运行中" || echo "❌ 未运行")"
echo "   orderer: $([ $ORDERER_UP -gt 0 ] && echo "✅ 运行中" || echo "❌ 未运行")"
echo "   CA服务: $([ $CA_UP -gt 0 ] && echo "✅ 运行中" || echo "❌ 未运行")"

echo ""
if [ $PEER1_UP -gt 0 ] && [ $PEER2_UP -gt 0 ] && [ $ORDERER_UP -gt 0 ]; then
    echo "🎉 Fabric网络已完全启动！"
    echo "✅ 所有节点正常运行"
    echo "✅ 可以开始测试ESG功能"
else
    echo "⏳ Fabric网络正在启动中..."
    echo "💡 请等待几分钟后再次检查"
fi 