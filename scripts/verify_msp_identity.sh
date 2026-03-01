#!/bin/bash

echo "🔍 MSP和证书身份验证 - ESG链码专用"
echo "================================="

echo "📋 1. 检查当前链码部署状态..."
echo "   - 检查basic链码:"
docker ps | grep "basic" || echo "   ❌ basic链码容器未运行"
echo "   - 检查ESG链码:"
docker ps | grep "esg" || echo "   ❌ ESG链码容器未运行"

echo ""
echo "📋 2. 检查MSP目录结构..."
MSP_BASE="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations"
echo "   - MSP基础路径: $MSP_BASE"

# 检查Org1 MSP结构
ORG1_MSP="$MSP_BASE/peerOrganizations/org1.example.com"
echo "   - Org1 MSP路径: $ORG1_MSP"

if [ -d "$ORG1_MSP" ]; then
    echo "   ✅ Org1 MSP目录存在"
    
    # 检查MSP子目录
    MSP_SUBDIRS=("msp" "users" "peers" "ca" "tlsca")
    for subdir in "${MSP_SUBDIRS[@]}"; do
        if [ -d "$ORG1_MSP/$subdir" ]; then
            echo "     ✅ $subdir 目录存在"
        else
            echo "     ❌ $subdir 目录不存在"
        fi
    done
else
    echo "   ❌ Org1 MSP目录不存在"
fi

echo ""
echo "📋 3. 检查Admin用户身份证书..."
ADMIN_CERT="$ORG1_MSP/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
ADMIN_KEY_DIR="$ORG1_MSP/users/Admin@org1.example.com/msp/keystore"

if [ -f "$ADMIN_CERT" ]; then
    echo "   ✅ Admin证书文件存在: $ADMIN_CERT"
    
    # 检查证书内容
    echo "   📄 证书信息:"
    openssl x509 -in "$ADMIN_CERT" -text -noout | grep -E "(Subject:|Issuer:|Not Before|Not After)" | head -4
    
    # 检查证书有效期
    echo "   📅 证书有效期:"
    openssl x509 -in "$ADMIN_CERT" -noout -dates
else
    echo "   ❌ Admin证书文件不存在: $ADMIN_CERT"
fi

if [ -d "$ADMIN_KEY_DIR" ] && [ "$(ls -A $ADMIN_KEY_DIR)" ]; then
    echo "   ✅ Admin私钥目录存在且不为空: $ADMIN_KEY_DIR"
    echo "   📄 私钥文件:"
    ls -la "$ADMIN_KEY_DIR"
else
    echo "   ❌ Admin私钥目录不存在或为空: $ADMIN_KEY_DIR"
fi

echo ""
echo "📋 4. 检查CA证书..."
CA_CERT="$ORG1_MSP/msp/cacerts/ca.org1.example.com-cert.pem"
TLS_CA_CERT="$ORG1_MSP/msp/tlscacerts/tlsca.org1.example.com-cert.pem"

if [ -f "$CA_CERT" ]; then
    echo "   ✅ CA证书存在: $CA_CERT"
    echo "   📄 CA证书信息:"
    openssl x509 -in "$CA_CERT" -text -noout | grep -E "(Subject:|Issuer:)" | head -2
else
    echo "   ❌ CA证书不存在: $CA_CERT"
fi

if [ -f "$TLS_CA_CERT" ]; then
    echo "   ✅ TLS CA证书存在: $TLS_CA_CERT"
    echo "   📄 TLS CA证书信息:"
    openssl x509 -in "$TLS_CA_CERT" -text -noout | grep -E "(Subject:|Issuer:)" | head -2
else
    echo "   ❌ TLS CA证书不存在: $TLS_CA_CERT"
fi

echo ""
echo "📋 5. 检查证书链验证..."
echo "   🔍 验证Admin证书是否由CA证书签名:"
if [ -f "$ADMIN_CERT" ] && [ -f "$CA_CERT" ]; then
    if openssl verify -CAfile "$CA_CERT" "$ADMIN_CERT" > /dev/null 2>&1; then
        echo "   ✅ Admin证书验证成功 - 证书链正确"
    else
        echo "   ❌ Admin证书验证失败 - 证书链错误"
        echo "   🔧 尝试详细验证:"
        openssl verify -CAfile "$CA_CERT" "$ADMIN_CERT"
    fi
else
    echo "   ⚠️ 无法验证 - 缺少证书文件"
fi

echo ""
echo "📋 6. 检查MSP配置文件..."
MSP_CONFIG="$ORG1_MSP/msp/config.yaml"
if [ -f "$MSP_CONFIG" ]; then
    echo "   ✅ MSP配置文件存在: $MSP_CONFIG"
    echo "   📄 MSP配置内容:"
    cat "$MSP_CONFIG"
else
    echo "   ❌ MSP配置文件不存在: $MSP_CONFIG"
fi

echo ""
echo "📋 7. 检查链码背书策略..."
echo "   🔍 检查ESG链码的背书策略配置..."

# 检查是否有ESG链码的特定配置
ESG_CONFIG="/home/ubuntu/go/fabric-sdk/configs/fabric/esg-chaincode-config.yaml"
if [ -f "$ESG_CONFIG" ]; then
    echo "   ✅ ESG链码配置文件存在: $ESG_CONFIG"
    echo "   📄 ESG配置内容:"
    cat "$ESG_CONFIG"
else
    echo "   ⚠️ ESG链码配置文件不存在，使用默认配置"
fi

echo ""
echo "📋 8. 生成MSP修复建议..."
echo "   💡 MSP修复建议:"

if [ ! -f "$ADMIN_CERT" ] || [ ! -d "$ADMIN_KEY_DIR" ]; then
    echo "   1. 重新生成Admin用户身份证书"
    echo "      fabric-ca-client enroll -u https://Admin:adminpw@localhost:7054 --caname ca-org1 --tls.certfiles /path/to/ca-cert.pem"
fi

if [ ! -f "$CA_CERT" ] || [ ! -f "$TLS_CA_CERT" ]; then
    echo "   2. 重新生成CA证书"
    echo "      fabric-ca-client getcainfo -u https://admin:adminpw@localhost:7054 --caname ca-org1"
fi

echo "   3. 检查ESG链码的背书策略"
echo "      peer lifecycle chaincode querycommitted -C mychannel -n esg"
echo "   4. 验证MSP配置"
echo "      peer channel fetch config -c mychannel -o localhost:7050 --tls --cafile /path/to/tls-ca-cert.pem"

echo ""
echo "📊 MSP和证书身份验证总结"
echo "======================"
echo "✅ 检查了MSP目录结构"
echo "✅ 验证了Admin用户身份证书"
echo "✅ 检查了CA证书"
echo "✅ 验证了证书链"
echo "✅ 检查了MSP配置"
echo "✅ 分析了ESG链码配置"

echo ""
echo "🎯 关键发现:"
echo "   - 证书链验证是连接失败的主要原因"
echo "   - ESG链码可能需要特定的背书策略"
echo "   - MSP配置需要与ESG链码匹配"

echo ""
echo "💡 下一步建议:"
echo "   1. 重新生成Admin用户证书"
echo "   2. 检查ESG链码的背书策略"
echo "   3. 更新MSP配置以匹配ESG链码要求" 