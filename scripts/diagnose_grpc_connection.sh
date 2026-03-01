#!/bin/bash

echo "🔍 gRPC连接诊断和修复 - 解决TRANSIENT_FAILURE"
echo "=========================================="

echo "📋 1. 检查peer容器状态..."
echo "   🔍 peer0.org1.example.com状态:"
docker ps | grep "peer0.org1.example.com" || echo "   ❌ peer0.org1.example.com未运行"

echo ""
echo "📋 2. 检查peer容器日志..."
echo "   📄 最近20行peer日志:"
docker logs --tail 20 peer0.org1.example.com 2>/dev/null | grep -E "(ERROR|WARN|FATAL|gRPC|grpc)" || echo "   无相关日志"

echo ""
echo "📋 3. 检查peer端口监听..."
echo "   🔍 检查端口7051监听状态:"
if ss -tlnp 2>/dev/null | grep ":7051"; then
    echo "   ✅ 端口7051正在监听"
else
    echo "   ❌ 端口7051未监听"
    echo "   🔍 检查Docker端口映射:"
    docker port peer0.org1.example.com | grep 7051 || echo "   端口映射异常"
fi

echo ""
echo "📋 4. 测试gRPC连接..."
echo "   🔍 使用grpcurl测试gRPC连接:"
if command -v grpcurl >/dev/null 2>&1; then
    echo "   📡 测试peer gRPC服务:"
    timeout 10 grpcurl -plaintext localhost:7051 list 2>/dev/null && echo "   ✅ gRPC连接成功" || echo "   ❌ gRPC连接失败"
else
    echo "   ⚠️ grpcurl未安装，使用telnet测试:"
    timeout 5 telnet localhost 7051 2>/dev/null && echo "   ✅ TCP连接成功" || echo "   ❌ TCP连接失败"
fi

echo ""
echo "📋 5. 检查TLS证书..."
echo "   🔍 检查peer TLS证书:"
TLS_CERT="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
if [ -f "$TLS_CERT" ]; then
    echo "   ✅ TLS证书存在: $TLS_CERT"
    echo "   📄 证书信息:"
    openssl x509 -in "$TLS_CERT" -text -noout | grep -E "(Subject:|Issuer:|Not Before|Not After)" | head -4
else
    echo "   ❌ TLS证书不存在: $TLS_CERT"
    echo "   🔍 查找其他TLS证书:"
    find /home/ubuntu/go/fabric-sdk/configs/fabric -name "*.crt" -type f 2>/dev/null | head -3
fi

echo ""
echo "📋 6. 检查网络配置..."
echo "   🔍 检查Docker网络:"
docker network ls | grep fabric || echo "   ⚠️ 未找到fabric网络"
echo "   🔍 检查peer容器网络配置:"
docker inspect peer0.org1.example.com | grep -A 10 "NetworkSettings" | head -10

echo ""
echo "📋 7. 重启peer容器..."
echo "   🔄 重启peer0.org1.example.com:"
docker restart peer0.org1.example.com
echo "   等待peer启动..."
sleep 15

echo "   🔍 检查重启后状态:"
docker ps | grep "peer0.org1.example.com" || echo "   ❌ peer重启失败"

echo ""
echo "📋 8. 验证修复效果..."
echo "   🔍 再次测试连接:"
sleep 5
if timeout 5 telnet localhost 7051 2>/dev/null; then
    echo "   ✅ peer端口现在可连接"
else
    echo "   ❌ peer端口仍然不可连接"
fi

echo ""
echo "📋 9. 检查gRPC服务..."
echo "   🔍 检查peer gRPC服务状态:"
docker exec peer0.org1.example.com ps aux | grep peer 2>/dev/null || echo "   无法检查peer进程"

echo ""
echo "📋 10. 生成修复建议..."
echo "   💡 gRPC连接修复建议:"
echo "   1. 检查peer容器是否正常运行"
echo "   2. 验证TLS证书配置"
echo "   3. 检查网络端口映射"
echo "   4. 重启peer容器"
echo "   5. 更新SDK配置中的gRPC选项"

echo ""
echo "📊 gRPC连接诊断总结"
echo "=================="
echo "✅ 检查了peer容器状态"
echo "✅ 分析了peer日志"
echo "✅ 验证了端口监听"
echo "✅ 测试了gRPC连接"
echo "✅ 检查了TLS证书"
echo "✅ 分析了网络配置"
echo "✅ 重启了peer容器"
echo "✅ 验证了修复效果"

echo ""
echo "🎯 关键发现:"
echo "   - TRANSIENT_FAILURE通常表示gRPC连接问题"
echo "   - 可能是TLS证书或网络配置问题"
echo "   - 重启peer容器可能解决问题"

echo ""
echo "💡 下一步建议:"
echo "   1. 如果peer重启成功，重新测试SDK连接"
echo "   2. 如果问题持续，检查TLS证书配置"
echo "   3. 考虑使用非TLS连接进行测试" 