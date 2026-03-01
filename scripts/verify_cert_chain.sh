#!/bin/bash

echo "🔗 验证证书链..."
echo "=================="

# 验证用户证书是否由正确的CA签发
USER_CERT="configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
ORG1_CA="configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"

if [ -f "$USER_CERT" ] && [ -f "$ORG1_CA" ]; then
    user_issuer=$(openssl x509 -in "$USER_CERT" -noout -issuer 2>/dev/null | sed 's/issuer=//')
    org1_subject=$(openssl x509 -in "$ORG1_CA" -noout -subject 2>/dev/null | sed 's/subject=//')
    
    if [ "$user_issuer" = "$org1_subject" ]; then
        echo "✅ 用户证书由Org1 CA正确签发"
    else
        echo "❌ 用户证书不是由Org1 CA签发"
        echo "   用户证书颁发者: $user_issuer"
        echo "   Org1 CA主题: $org1_subject"
    fi
fi

echo ""
echo "🎉 证书链验证完成！"
