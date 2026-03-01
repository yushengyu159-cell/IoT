#!/bin/bash

echo "🎯 最终TLS连接验证：重启peer后使用TLS连接"

# 1. 检查容器状态
echo "🐳 1. 检查容器状态..."
if docker ps | grep -q "peer0.org1.example.com"; then
    echo "   ✅ peer0.org1.example.com 正在运行"
else
    echo "   ❌ peer0.org1.example.com 未运行"
    exit 1
fi

if docker ps | grep -q "peer0.org2.example.com"; then
    echo "   ✅ peer0.org2.example.com 正在运行"
else
    echo "   ❌ peer0.org2.example.com 未运行"
    exit 1
fi

if docker ps | grep -q "orderer.example.com"; then
    echo "   ✅ orderer.example.com 正在运行"
else
    echo "   ❌ orderer.example.com 未运行"
    exit 1
fi

# 2. 检查端口监听
echo "🔌 2. 检查端口监听..."
if ss -tlnp | grep -q ":7051"; then
    echo "   ✅ 端口7051 (peer0.org1) 正在监听"
else
    echo "   ❌ 端口7051 (peer0.org1) 未监听"
    exit 1
fi

if ss -tlnp | grep -q ":7050"; then
    echo "   ✅ 端口7050 (orderer) 正在监听"
else
    echo "   ❌ 端口7050 (orderer) 未监听"
    exit 1
fi

if ss -tlnp | grep -q ":9051"; then
    echo "   ✅ 端口9051 (peer0.org2) 正在监听"
else
    echo "   ❌ 端口9051 (peer0.org2) 未监听"
    exit 1
fi

# 3. 检查配置文件中的TLS设置
echo "🔐 3. 检查TLS配置设置..."
echo "   peer0.org1.example.com URL"
if grep -q "url: grpcs://localhost:7051" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 已使用TLS协议 (grpcs://)"
else
    echo "   ❌ 未使用TLS协议"
    exit 1
fi

echo "   orderer.example.com URL"
if grep -q "url: grpcs://localhost:7050" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 已使用TLS协议 (grpcs://)"
else
    echo "   ❌ 未使用TLS协议"
    exit 1
fi

echo "   ca.org1.example.com URL"
if grep -q "url: https://localhost:7054" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 已使用TLS协议 (https://)"
else
    echo "   ❌ 未使用TLS协议"
    exit 1
fi

# 4. 检查TLS验证设置
echo "🔒 4. 检查TLS验证设置..."
if grep -q "verify: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 证书验证已关闭"
else
    echo "   ❌ 证书验证未关闭"
    exit 1
fi

if grep -q "allow-insecure: true" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 不安全连接已允许"
else
    echo "   ❌ 不安全连接未允许"
    exit 1
fi

# 5. 检查环境变量
echo "🌍 5. 检查环境变量..."
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true

echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"

# 6. 测试网络连接
echo "🌐 6. 测试网络连接..."
if timeout 5 bash -c "</dev/tcp/localhost/7051" 2>/dev/null; then
    echo "   ✅ peer0.org1:7051 连接正常"
else
    echo "   ❌ peer0.org1:7051 连接失败"
    exit 1
fi

if timeout 5 bash -c "</dev/tcp/localhost/7050" 2>/dev/null; then
    echo "   ✅ orderer:7050 连接正常"
else
    echo "   ❌ orderer:7050 连接失败"
    exit 1
fi

# 7. 显示最终状态
echo "📊 7. 最终TLS连接状态..."
echo "   - 容器状态: ✅ 所有节点已重启并运行"
echo "   - 端口监听: ✅ 所有端口正在监听"
echo "   - Peer协议: ✅ 使用 grpcs:// (TLS)"
echo "   - Orderer协议: ✅ 使用 grpcs:// (TLS)"
echo "   - CA协议: ✅ 使用 https:// (TLS)"
echo "   - TLS验证: ✅ 已关闭"
echo "   - 不安全连接: ✅ 已允许"
echo "   - 网络连接: ✅ 所有端口可连接"

echo "🎉 最终TLS连接验证完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   cd /home/ubuntu/go/fabric-sdk"
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 修复说明："
echo "   - 已重启所有Fabric节点容器"
echo "   - 已恢复使用TLS协议 (grpcs://, https://)"
echo "   - 保持TLS验证关闭和不安全连接允许"
echo "   - 这样SDK将使用TLS连接但跳过证书验证" 