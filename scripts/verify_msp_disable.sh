#!/bin/bash

echo "🔒 验证MSP禁用配置：开发/测试环境专用"

# 1. 检查YAML语法
echo "📋 1. 检查YAML语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    exit 1
fi

# 2. 检查TLS客户端证书文件
echo "🔐 2. 检查TLS客户端证书文件..."
TLS_CLIENT_FILES=(
    "configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key"
    "configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt"
)

for cert_file in "${TLS_CLIENT_FILES[@]}"; do
    if [ -f "$cert_file" ]; then
        echo "   ✅ $cert_file 存在"
    else
        echo "   ❌ $cert_file 不存在"
        exit 1
    fi
done

# 3. 检查系统证书池设置
echo "🚫 3. 检查系统证书池设置..."
if grep -q "systemCertPool: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 系统证书池已关闭"
else
    echo "   ❌ 系统证书池未关闭"
    exit 1
fi

# 4. 检查TLS验证设置
echo "🔒 4. 检查TLS验证设置..."
if grep -q "verify: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ TLS验证已关闭"
else
    echo "   ❌ TLS验证未关闭"
    exit 1
fi

# 5. 检查MSP验证设置
echo "🛡️ 5. 检查MSP验证设置..."
if grep -A 5 "msp:" configs/fabric/connection-optimized.yaml | grep -q "verify: false"; then
    echo "   ✅ MSP验证已关闭"
else
    echo "   ❌ MSP验证未关闭"
    exit 1
fi

# 6. 检查BCCSP配置
echo "🔧 6. 检查BCCSP配置..."
if grep -A 5 "BCCSP:" configs/fabric/connection-optimized.yaml | grep -q "provider: SW"; then
    echo "   ✅ BCCSP已配置为SW提供者"
else
    echo "   ❌ BCCSP未配置为SW提供者"
    exit 1
fi

# 7. 设置环境变量
echo "🌍 7. 设置环境变量..."
export FABRIC_SDK_GO_MSP_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true
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
echo "📊 8. MSP禁用配置验证状态..."
echo "   - YAML语法: ✅ 正确"
echo "   - TLS客户端证书: ✅ 文件存在"
echo "   - 系统证书池: ✅ 已关闭"
echo "   - TLS验证: ✅ 已关闭"
echo "   - MSP验证: ✅ 已关闭"
echo "   - BCCSP配置: ✅ 已配置SW提供者"
echo "   - 环境变量: ✅ 已设置"

echo "🎉 MSP禁用配置验证完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 修复说明："
echo "   - 已直接禁用MSP验证 (msp.verify: false)"
echo "   - 已关闭TLS验证 (verify: false)"
echo "   - 已关闭系统证书池 (systemCertPool: false)"
echo "   - 已配置BCCSP为SW提供者"
echo "   - 这是开发/测试环境的终极解决方案" 