#!/bin/bash

# 验证证书池脚本
echo "🔍 验证证书池..."

# 创建临时证书池文件
CERT_POOL_FILE="/tmp/fabric_cert_pool.pem"

# 清空证书池
> "$CERT_POOL_FILE"

# 添加Org1 CA证书
if [[ -f "organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem" ]]; then
    cat "organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem" >> "$CERT_POOL_FILE"
    echo "✅ 已添加Org1 CA证书到证书池"
fi

# 添加Org1 TLS CA证书
if [[ -f "organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem" ]]; then
    cat "organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem" >> "$CERT_POOL_FILE"
    echo "✅ 已添加Org1 TLS CA证书到证书池"
fi

# 添加Orderer CA证书
if [[ -f "organizations/ordererOrganizations/example.com/msp/cacerts/ca.example.com-cert.pem" ]]; then
    cat "organizations/ordererOrganizations/example.com/msp/cacerts/ca.example.com-cert.pem" >> "$CERT_POOL_FILE"
    echo "✅ 已添加Orderer CA证书到证书池"
fi

# 添加Orderer TLS CA证书
if [[ -f "organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem" ]]; then
    cat "organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem" >> "$CERT_POOL_FILE"
    echo "✅ 已添加Orderer TLS CA证书到证书池"
fi

# 验证Admin证书
if [[ -f "organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem" ]]; then
    if openssl verify -CAfile "$CERT_POOL_FILE" "organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem" >/dev/null 2>&1; then
        echo "✅ Admin证书验证通过"
    else
        echo "❌ Admin证书验证失败"
    fi
fi

echo "📋 证书池文件位置: $CERT_POOL_FILE"
echo "📋 证书池内容:"
openssl crl2pkcs7 -nocrl -certfile "$CERT_POOL_FILE" | openssl pkcs7 -print_certs -noout | grep "Subject:"
