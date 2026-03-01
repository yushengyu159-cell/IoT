#!/bin/bash

echo "🔍 验证entityMatchers配置：彻底关闭系统根证书池"

# 1. 检查YAML语法
echo "📋 1. 检查YAML语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    exit 1
fi

# 2. 检查entityMatchers配置
echo "🎯 2. 检查entityMatchers配置..."
if grep -A 20 "entityMatchers:" configs/fabric/connection-optimized.yaml | grep -q "tlsCACerts:"; then
    echo "   ✅ entityMatchers中已配置tlsCACerts"
else
    echo "   ❌ entityMatchers中未配置tlsCACerts"
    exit 1
fi

# 3. 检查TLS CA证书文件
echo "🔐 3. 检查TLS CA证书文件..."
TLS_CERT_FILES=(
    "configs/fabric/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem"
    "configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem"
    "configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
)

for cert_file in "${TLS_CERT_FILES[@]}"; do
    if [ -f "$cert_file" ]; then
        echo "   ✅ $cert_file 存在"
    else
        echo "   ❌ $cert_file 不存在"
        exit 1
    fi
done

# 4. 检查系统证书池设置
echo "🚫 4. 检查系统证书池设置..."
if grep -q "systemCertPool: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 系统证书池已关闭"
else
    echo "   ❌ 系统证书池未关闭"
    exit 1
fi

# 5. 检查证书验证设置
echo "🔒 5. 检查证书验证设置..."
if grep -q "verify: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 证书验证已关闭"
else
    echo "   ❌ 证书验证未关闭"
    exit 1
fi

# 6. 检查MSP CA证书注入
echo "📄 6. 检查MSP CA证书注入..."
if grep -A 20 "msp:" configs/fabric/connection-optimized.yaml | grep -q "BEGIN CERTIFICATE"; then
    echo "   ✅ MSP CA证书已注入配置"
else
    echo "   ❌ MSP CA证书未注入配置"
    exit 1
fi

# 7. 设置环境变量
echo "🌍 7. 设置环境变量..."
export FABRIC_SDK_GO_MSP_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=true
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=false
export FABRIC_SDK_GO_LOG_LEVEL=ERROR
export FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST=true

echo "   FABRIC_SDK_GO_MSP_VERIFY: $FABRIC_SDK_GO_MSP_VERIFY"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"
echo "   FABRIC_SDK_GO_LOG_LEVEL: $FABRIC_SDK_GO_LOG_LEVEL"
echo "   FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST: $FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST"

# 8. 显示修复状态
echo "📊 8. entityMatchers配置验证状态..."
echo "   - YAML语法: ✅ 正确"
echo "   - entityMatchers配置: ✅ 已配置tlsCACerts"
echo "   - TLS CA证书文件: ✅ 所有文件存在"
echo "   - 系统证书池: ✅ 已关闭"
echo "   - 证书验证: ✅ 已关闭"
echo "   - MSP CA证书: ✅ 已注入"
echo "   - 环境变量: ✅ 已设置"

echo "🎉 entityMatchers配置验证完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 修复说明："
echo "   - 已在entityMatchers中显式指定TLS CA证书路径"
echo "   - 已彻底关闭系统根证书池"
echo "   - 已强制使用Fabric自己的CA证书"
echo "   - 这样应该能彻底解决Discovery阶段证书验证问题" 