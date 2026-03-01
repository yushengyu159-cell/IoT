#!/bin/bash

echo "🔧 强制覆盖MSP CA证书路径..."

# 获取当前工作目录
CURRENT_DIR=$(pwd)
echo "当前工作目录: $CURRENT_DIR"

# 定义MSP路径
MSP_BASE_PATH="$CURRENT_DIR/configs/fabric/organizations/peerOrganizations/org1.example.com/msp"
CA_CERTS_PATH="$MSP_BASE_PATH/cacerts/ca.org1.example.com-cert.pem"
TLS_CA_CERTS_PATH="$MSP_BASE_PATH/tlscacerts/tlsca.org1.example.com-cert.pem"

echo "MSP基础路径: $MSP_BASE_PATH"
echo "CA证书路径: $CA_CERTS_PATH"
echo "TLS CA证书路径: $TLS_CA_CERTS_PATH"

# 验证文件是否存在
if [ -f "$CA_CERTS_PATH" ]; then
    echo "✅ CA证书文件存在: $CA_CERTS_PATH"
else
    echo "❌ CA证书文件不存在: $CA_CERTS_PATH"
    exit 1
fi

if [ -f "$TLS_CA_CERTS_PATH" ]; then
    echo "✅ TLS CA证书文件存在: $TLS_CA_CERTS_PATH"
else
    echo "❌ TLS CA证书文件不存在: $TLS_CA_CERTS_PATH"
    exit 1
fi

# 设置环境变量
export FABRIC_SDK_GO_MSP_CACERTS_PATH="$CA_CERTS_PATH"
export FABRIC_SDK_GO_MSP_TLSCACERTS_PATH="$TLS_CA_CERTS_PATH"
export FABRIC_SDK_GO_LOG_LEVEL="FATAL"
export FABRIC_SDK_GO_MSP_VERIFY="false"

echo "✅ 环境变量已设置:"
echo "   FABRIC_SDK_GO_MSP_CACERTS_PATH=$FABRIC_SDK_GO_MSP_CACERTS_PATH"
echo "   FABRIC_SDK_GO_MSP_TLSCACERTS_PATH=$FABRIC_SDK_GO_MSP_TLSCACERTS_PATH"

# 验证配置文件中的MSP路径
echo "🔍 验证配置文件中的MSP路径..."
if grep -q "caCerts:.*ca.org1.example.com-cert.pem" configs/fabric/connection-optimized.yaml; then
    echo "✅ 配置文件中包含CA证书路径"
else
    echo "❌ 配置文件中缺少CA证书路径"
fi

if grep -q "tlscacerts:.*tlsca.org1.example.com-cert.pem" configs/fabric/connection-optimized.yaml; then
    echo "✅ 配置文件中包含TLS CA证书路径"
else
    echo "❌ 配置文件中缺少TLS CA证书路径"
fi

echo "🎉 MSP路径强制覆盖完成！" 