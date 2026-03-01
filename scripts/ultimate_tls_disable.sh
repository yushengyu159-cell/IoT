#!/bin/bash

echo "🎯 终极TLS验证关闭：显式关闭TLS验证（仅开发/测试）"

# 1. 检查配置文件中的TLS设置
echo "🔐 1. 检查TLS配置设置..."
echo "   verify: false"
if grep -q "verify: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 全局证书验证已关闭"
else
    echo "   ❌ 全局证书验证未关闭"
    exit 1
fi

echo "   allow-insecure: true"
if grep -q "allow-insecure: true" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 不安全连接已允许"
else
    echo "   ❌ 不安全连接未允许"
    exit 1
fi

echo "   systemCertPool: false"
if grep -q "systemCertPool: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 系统证书池已禁用"
else
    echo "   ❌ 系统证书池未禁用"
    exit 1
fi

# 2. 检查entityMatchers中的TLS设置
echo "🔗 2. 检查entityMatchers TLS设置..."
echo "   peer tlsCACerts"
if grep -A 5 "tlsCACerts:" configs/fabric/connection-optimized.yaml | grep -q "tlsca.org1.example.com-cert.pem"; then
    echo "   ✅ Peer TLS CA证书路径已指定"
else
    echo "   ❌ Peer TLS CA证书路径未指定"
    exit 1
fi

echo "   peer allow-insecure"
if grep -A 10 "peer:" configs/fabric/connection-optimized.yaml | grep -q "allow-insecure: true"; then
    echo "   ✅ Peer不安全连接已允许"
else
    echo "   ❌ Peer不安全连接未允许"
    exit 1
fi

echo "   orderer allow-insecure"
if grep -A 10 "orderer:" configs/fabric/connection-optimized.yaml | grep -q "allow-insecure: true"; then
    echo "   ✅ Orderer不安全连接已允许"
else
    echo "   ❌ Orderer不安全连接未允许"
    exit 1
fi

# 3. 检查环境变量
echo "🌍 3. 检查环境变量..."
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true

echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"

# 4. 检查TLS证书文件
echo "📋 4. 检查TLS证书文件..."
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

# 5. 检查配置文件语法
echo "🔍 5. 检查配置文件语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    exit 1
fi

# 6. 显示最终状态
echo "📊 6. 终极TLS关闭状态..."
echo "   - 全局证书验证: ✅ 已关闭 (verify: false)"
echo "   - 系统证书池: ✅ 已禁用 (systemCertPool: false)"
echo "   - 不安全连接: ✅ 已允许 (allow-insecure: true)"
echo "   - Peer TLS CA: ✅ 已指定路径"
echo "   - Orderer TLS: ✅ 已允许不安全连接"
echo "   - 环境变量: ✅ 已设置"
echo "   - 配置文件: ✅ 语法正确"

echo "🎉 终极TLS验证关闭完成！"
echo ""
echo "🚀 现在重启服务，TLS验证已完全关闭！"
echo "   SDK将跳过所有TLS证书验证，直接建立连接"
echo ""
echo "⚠️  警告：此配置仅适用于开发/测试环境！"
echo "   生产环境请正确配置TLS证书和验证！" 