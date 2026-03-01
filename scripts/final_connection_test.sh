#!/bin/bash

echo "🎯 最终连接测试：验证SDK配置和网络连接"

# 设置环境变量
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true

echo "🌍 环境变量已设置："
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"

# 1. 检查配置文件
echo "📋 1. 检查配置文件..."
if [ -f "configs/fabric/connection-optimized.yaml" ]; then
    echo "   ✅ 配置文件存在"
else
    echo "   ❌ 配置文件不存在"
    exit 1
fi

# 2. 检查配置文件语法
echo "🔍 2. 检查配置文件语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    exit 1
fi

# 3. 检查关键配置项
echo "⚙️ 3. 检查关键配置项..."
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

if grep -q "systemCertPool: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 系统证书池已禁用"
else
    echo "   ❌ 系统证书池未禁用"
    exit 1
fi

# 4. 检查网络连接
echo "🌐 4. 检查网络连接..."
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

# 5. 检查证书文件
echo "🔐 5. 检查证书文件..."
CERT_FILES=(
    "configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
    "configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
    "configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"
    "configs/fabric/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem"
)

for cert_file in "${CERT_FILES[@]}"; do
    if [ -f "$cert_file" ]; then
        echo "   ✅ $cert_file 存在"
    else
        echo "   ❌ $cert_file 不存在"
        exit 1
    fi
done

# 6. 显示最终状态
echo "📊 6. 最终连接测试状态..."
echo "   - 配置文件: ✅ 存在且语法正确"
echo "   - 证书验证: ✅ 已关闭"
echo "   - 不安全连接: ✅ 已允许"
echo "   - 系统证书池: ✅ 已禁用"
echo "   - 网络连接: ✅ 所有节点可连接"
echo "   - 证书文件: ✅ 所有必要文件存在"
echo "   - 环境变量: ✅ 已设置"

echo "🎉 最终连接测试完成！"
echo ""
echo "🚀 现在可以重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 如果仍有问题，请检查："
echo "   1. SDK版本兼容性"
echo "   2. Go模块依赖"
echo "   3. 链码背书策略" 