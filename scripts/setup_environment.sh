#!/bin/bash

echo "🔧 设置Fabric SDK环境变量 - 最终修复版本"

# 强制设置MSP路径
export FABRIC_SDK_CRYPTO_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
export FABRIC_SDK_CONFIG_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric/connection-optimized.yaml"

# 强制设置证书路径
export FABRIC_SDK_CERT_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
export FABRIC_SDK_KEY_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
export FABRIC_SDK_CA_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"

# 禁用证书验证
export FABRIC_SDK_VERIFY_TLS="false"
export FABRIC_SDK_VERIFY_MSP="false"
export FABRIC_SDK_SYSTEM_CERT_POOL="false"

# 设置网络配置
export FABRIC_SDK_NETWORK_CONFIG="/home/ubuntu/go/fabric-sdk/configs/fabric/connection-optimized.yaml"
export FABRIC_SDK_ORG="Org1"
export FABRIC_SDK_MSPID="Org1MSP"

# 设置连接参数
export FABRIC_SDK_PEER_URL="grpcs://localhost:7051"
export FABRIC_SDK_ORDERER_URL="grpcs://localhost:7050"
export FABRIC_SDK_CA_URL="https://localhost:7054"

# 设置链码参数
export FABRIC_SDK_CHANNEL_NAME="mychannel"
export FABRIC_SDK_CHAINCODE_NAME="esg"

# 禁用系统证书池
export GODEBUG="x509ignoreCN=0"
export CGO_ENABLED=1

# 设置Go环境
export GOPATH="/home/ubuntu/go"
export GOROOT="/usr/local/go"
export PATH="$GOROOT/bin:$GOPATH/bin:$PATH"

echo "✅ 环境变量设置完成"
echo "📋 关键配置:"
echo "   - MSP路径: $FABRIC_SDK_CRYPTO_PATH"
echo "   - 配置文件: $FABRIC_SDK_CONFIG_PATH"
echo "   - TLS验证: $FABRIC_SDK_VERIFY_TLS"
echo "   - MSP验证: $FABRIC_SDK_VERIFY_MSP"
echo "   - 系统证书池: $FABRIC_SDK_SYSTEM_CERT_POOL"

# 验证关键文件存在
echo "📋 验证关键文件:"
if [ -f "$FABRIC_SDK_CERT_PATH" ]; then
    echo "   ✅ 证书文件存在: $FABRIC_SDK_CERT_PATH"
else
    echo "   ❌ 证书文件不存在: $FABRIC_SDK_CERT_PATH"
fi

if [ -f "$FABRIC_SDK_KEY_PATH" ]; then
    echo "   ✅ 私钥文件存在: $FABRIC_SDK_KEY_PATH"
else
    echo "   ❌ 私钥文件不存在: $FABRIC_SDK_KEY_PATH"
fi

if [ -f "$FABRIC_SDK_CA_PATH" ]; then
    echo "   ✅ CA证书文件存在: $FABRIC_SDK_CA_PATH"
else
    echo "   ❌ CA证书文件不存在: $FABRIC_SDK_CA_PATH"
fi

if [ -f "$FABRIC_SDK_CONFIG_PATH" ]; then
    echo "   ✅ 配置文件存在: $FABRIC_SDK_CONFIG_PATH"
else
    echo "   ❌ 配置文件不存在: $FABRIC_SDK_CONFIG_PATH"
fi

echo "🎉 环境变量设置完成，可以启动SDK服务" 