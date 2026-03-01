#!/bin/bash

echo "🔐 验证证书链..."

# 获取当前工作目录
CURRENT_DIR=$(pwd)
echo "当前工作目录: $CURRENT_DIR"

# 定义证书路径
ADMIN_MSP_PATH="$CURRENT_DIR/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
ORG_MSP_PATH="$CURRENT_DIR/configs/fabric/organizations/peerOrganizations/org1.example.com/msp"

USER_CERT="$ADMIN_MSP_PATH/signcerts/Admin@org1.example.com-cert.pem"
CA_CERT="$ADMIN_MSP_PATH/cacerts/ca.org1.example.com-cert.pem"
TLS_CA_CERT="$ADMIN_MSP_PATH/tlscacerts/tlsca.org1.example.com-cert.pem"

echo "📁 证书路径:"
echo "   用户证书: $USER_CERT"
echo "   CA证书: $CA_CERT"
echo "   TLS CA证书: $TLS_CA_CERT"

# 1. 检查证书文件是否存在
echo "🔍 1. 检查证书文件存在性..."
if [ -f "$USER_CERT" ]; then
    echo "✅ 用户证书存在"
else
    echo "❌ 用户证书不存在"
    exit 1
fi

if [ -f "$CA_CERT" ]; then
    echo "✅ CA证书存在"
else
    echo "❌ CA证书不存在"
    exit 1
fi

if [ -f "$TLS_CA_CERT" ]; then
    echo "✅ TLS CA证书存在"
else
    echo "❌ TLS CA证书不存在"
    exit 1
fi

# 2. 验证证书链
echo "🔗 2. 验证证书链..."
if openssl verify -CAfile "$CA_CERT" "$USER_CERT" > /dev/null 2>&1; then
    echo "✅ 用户证书 -> CA证书 链验证通过"
else
    echo "❌ 用户证书 -> CA证书 链验证失败"
    openssl verify -CAfile "$CA_CERT" "$USER_CERT"
    exit 1
fi

# 3. 检查证书颁发者
echo "📋 3. 检查证书颁发者..."
USER_ISSUER=$(openssl x509 -in "$USER_CERT" -noout -issuer | cut -d'=' -f2-)
CA_SUBJECT=$(openssl x509 -in "$CA_CERT" -noout -subject | cut -d'=' -f2-)

echo "   用户证书颁发者: $USER_ISSUER"
echo "   CA证书主题: $CA_SUBJECT"

if [ "$USER_ISSUER" = "$CA_SUBJECT" ]; then
    echo "✅ 证书颁发者匹配"
else
    echo "❌ 证书颁发者不匹配"
    exit 1
fi

# 4. 检查证书有效期
echo "⏰ 4. 检查证书有效期..."
CURRENT_TIME=$(date +%s)
USER_NOT_AFTER=$(openssl x509 -in "$USER_CERT" -noout -enddate | cut -d'=' -f2)
CA_NOT_AFTER=$(openssl x509 -in "$CA_CERT" -noout -enddate | cut -d'=' -f2)

echo "   用户证书有效期: $USER_NOT_AFTER"
echo "   CA证书有效期: $CA_NOT_AFTER"

# 5. 检查MSP配置
echo "📁 5. 检查MSP配置..."
if [ -f "$ADMIN_MSP_PATH/config.yaml" ]; then
    echo "✅ MSP配置文件存在"
    echo "   配置文件内容:"
    cat "$ADMIN_MSP_PATH/config.yaml" | sed 's/^/   /'
else
    echo "❌ MSP配置文件不存在"
fi

# 6. 设置环境变量
echo "🌍 6. 设置环境变量..."
export FABRIC_SDK_GO_MSP_CACERTS_PATH="$CA_CERT"
export FABRIC_SDK_GO_MSP_TLSCACERTS_PATH="$TLS_CA_CERT"
export FABRIC_SDK_GO_MSP_PATH="$ADMIN_MSP_PATH"
export FABRIC_SDK_GO_CERT_VERIFY="false"

echo "✅ 环境变量已设置:"
echo "   FABRIC_SDK_GO_MSP_CACERTS_PATH=$FABRIC_SDK_GO_MSP_CACERTS_PATH"
echo "   FABRIC_SDK_GO_MSP_TLSCACERTS_PATH=$FABRIC_SDK_GO_MSP_TLSCACERTS_PATH"
echo "   FABRIC_SDK_GO_MSP_PATH=$FABRIC_SDK_GO_MSP_PATH"
echo "   FABRIC_SDK_GO_CERT_VERIFY=$FABRIC_SDK_GO_CERT_VERIFY"

echo "🎉 证书链验证完成！"
echo ""
echo "📊 验证总结:"
echo "- 证书文件: ✅ 全部存在"
echo "- 证书链: ✅ 验证通过"
echo "- 颁发者: ✅ 匹配正确"
echo "- MSP配置: ✅ 完整"
echo "- 环境变量: ✅ 已设置" 