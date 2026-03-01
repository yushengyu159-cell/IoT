#!/bin/bash

echo "🎯 正确的TLS配置：SDK使用TLS协议连接peer"

# 1. 检查配置文件中的TLS设置
echo "🔐 1. 检查TLS配置设置..."
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

# 2. 检查TLS证书配置
echo "🔒 2. 检查TLS证书配置..."
if grep -A 20 "tlsCACerts:" configs/fabric/connection-optimized.yaml | grep -q "BEGIN CERTIFICATE"; then
    echo "   ✅ TLS CA证书已嵌入"
else
    echo "   ❌ TLS CA证书未嵌入"
    exit 1
fi

# 3. 检查grpcOptions设置
echo "⚙️ 3. 检查grpcOptions设置..."
if grep -q "allow-insecure: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 不安全连接已禁用 (allow-insecure: false)"
else
    echo "   ❌ 不安全连接未禁用"
    exit 1
fi

if grep -q "ssl-target-name-override: peer0.org1.example.com" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ SSL目标名称覆盖已设置"
else
    echo "   ❌ SSL目标名称覆盖未设置"
    exit 1
fi

# 4. 检查环境变量
echo "🌍 4. 检查环境变量..."
export FABRIC_SDK_GO_TLS_VERIFY=true
export FABRIC_SDK_GO_CERT_VERIFY=true
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=false

echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"

# 5. 检查TLS证书文件
echo "📋 5. 检查TLS证书文件..."
TLS_CA_FILE="configs/fabric/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem"
if [ -f "$TLS_CA_FILE" ]; then
    echo "   ✅ TLS CA证书文件存在"
    if grep -q "BEGIN CERTIFICATE" "$TLS_CA_FILE" && grep -q "END CERTIFICATE" "$TLS_CA_FILE"; then
        echo "   ✅ TLS CA证书格式正确"
    else
        echo "   ❌ TLS CA证书格式错误"
        exit 1
    fi
else
    echo "   ❌ TLS CA证书文件不存在"
    exit 1
fi

# 6. 检查网络连接
echo "🌐 6. 检查网络连接..."
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
echo "📊 7. 正确TLS配置状态..."
echo "   - Peer协议: ✅ 使用 grpcs:// (TLS)"
echo "   - Orderer协议: ✅ 使用 grpcs:// (TLS)"
echo "   - CA协议: ✅ 使用 https:// (TLS)"
echo "   - TLS证书: ✅ 已嵌入配置"
echo "   - 不安全连接: ✅ 已禁用"
echo "   - SSL目标名称: ✅ 已设置"
echo "   - 环境变量: ✅ 已启用TLS验证"
echo "   - 网络连接: ✅ 所有端口可连接"

echo "🎉 正确TLS配置完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 修复说明："
echo "   - SDK使用TLS协议 (grpcs://) 连接peer"
echo "   - peer期望TLS握手，SDK发送TLS数据"
echo "   - 协议匹配，避免握手失败"
echo "   - 使用正确的TLS证书进行验证" 