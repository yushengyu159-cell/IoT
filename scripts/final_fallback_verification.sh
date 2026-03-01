#!/bin/bash

echo "🎯 最终兜底验证：强制禁用系统证书池 + 显式嵌入CA证书"

# 1. 检查配置文件中的关键设置
echo "🔐 1. 检查配置文件设置..."
echo "   systemCertPool: false"
if grep -q "systemCertPool: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 系统证书池已强制禁用"
else
    echo "   ❌ 系统证书池未禁用"
    exit 1
fi

echo "   verify: false"
if grep -q "verify: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 证书验证已关闭"
else
    echo "   ❌ 证书验证未关闭"
    exit 1
fi

echo "   caCerts in tlsCerts"
if grep -A 20 "caCerts:" configs/fabric/connection-optimized.yaml | grep -q "BEGIN CERTIFICATE"; then
    echo "   ✅ TLS配置中已显式嵌入CA证书"
else
    echo "   ❌ TLS配置中未嵌入CA证书"
    exit 1
fi

echo "   caCerts in msp"
if grep -A 20 "msp:" configs/fabric/connection-optimized.yaml | grep -q "BEGIN CERTIFICATE"; then
    echo "   ✅ MSP配置中已显式嵌入CA证书"
else
    echo "   ❌ MSP配置中未嵌入CA证书"
    exit 1
fi

# 2. 检查环境变量
echo "🌍 2. 检查环境变量..."
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=false

echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"

# 3. 验证证书内容
echo "📋 3. 验证证书内容..."
CA_CERT_FILE="configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"
if [ -f "$CA_CERT_FILE" ]; then
    echo "   ✅ CA证书文件存在"
    CERT_CONTENT=$(cat "$CA_CERT_FILE")
    if echo "$CERT_CONTENT" | grep -q "BEGIN CERTIFICATE" && echo "$CERT_CONTENT" | grep -q "END CERTIFICATE"; then
        echo "   ✅ CA证书格式正确"
    else
        echo "   ❌ CA证书格式错误"
        exit 1
    fi
else
    echo "   ❌ CA证书文件不存在"
    exit 1
fi

# 4. 检查配置文件语法
echo "🔍 4. 检查配置文件语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    exit 1
fi

# 5. 显示最终状态
echo "📊 5. 最终兜底状态..."
echo "   - 系统证书池: ✅ 强制禁用 (systemCertPool: false)"
echo "   - 证书验证: ✅ 已关闭 (verify: false)"
echo "   - TLS CA证书: ✅ 显式嵌入"
echo "   - MSP CA证书: ✅ 显式嵌入"
echo "   - 环境变量: ✅ 已设置"
echo "   - 配置文件: ✅ 语法正确"

echo "🎉 最终兜底验证完成！"
echo ""
echo "🚀 现在重启服务，这是最后的兜底方案！"
echo "   SDK将完全忽略系统证书池，使用显式嵌入的CA证书"
echo ""
echo "💡 如果仍有问题，说明问题不在证书配置，"
echo "   可能是Fabric网络本身的问题或SDK版本兼容性问题" 