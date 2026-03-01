#!/bin/bash

# Fabric证书信任链修复脚本
# 解决 x509: certificate signed by unknown authority 问题

set -e

echo "🔍 开始检查Fabric证书信任链..."

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置路径
CONFIG_BASE="go/fabric-sdk/configs/fabric"
CURRENT_DIR=$(pwd)

# 检查函数
check_file() {
    local file="$1"
    local description="$2"
    
    if [[ -f "$file" ]]; then
        echo -e "${GREEN}✅ $description: $file${NC}"
        return 0
    else
        echo -e "${RED}❌ $description: $file (文件不存在)${NC}"
        return 1
    fi
}

check_cert_validity() {
    local cert_file="$1"
    local description="$2"
    
    if [[ -f "$cert_file" ]]; then
        if openssl x509 -in "$cert_file" -noout -text >/dev/null 2>&1; then
            echo -e "${GREEN}✅ $description: 证书格式有效${NC}"
            
            # 检查证书有效期
            local not_after=$(openssl x509 -in "$cert_file" -noout -enddate | cut -d= -f2)
            local expiry_date=$(date -d "$not_after" +%s)
            local current_date=$(date +%s)
            
            if [[ $expiry_date -gt $current_date ]]; then
                echo -e "${GREEN}✅ $description: 证书未过期 (过期时间: $not_after)${NC}"
            else
                echo -e "${RED}❌ $description: 证书已过期 (过期时间: $not_after)${NC}"
            fi
        else
            echo -e "${RED}❌ $description: 证书格式无效${NC}"
        fi
    else
        echo -e "${RED}❌ $description: 证书文件不存在${NC}"
    fi
}

verify_cert_chain() {
    local cert_file="$1"
    local ca_file="$2"
    local description="$3"
    
    if [[ -f "$cert_file" && -f "$ca_file" ]]; then
        if openssl verify -CAfile "$ca_file" "$cert_file" >/dev/null 2>&1; then
            echo -e "${GREEN}✅ $description: 证书链验证通过${NC}"
        else
            echo -e "${RED}❌ $description: 证书链验证失败${NC}"
            echo -e "${YELLOW}   证书: $cert_file${NC}"
            echo -e "${YELLOW}   CA: $ca_file${NC}"
        fi
    else
        echo -e "${RED}❌ $description: 证书或CA文件不存在${NC}"
    fi
}

echo -e "${BLUE}📋 1. 检查关键证书文件...${NC}"

# 检查Org1的CA证书
check_cert_validity "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem" "Org1 CA证书"

# 检查Org1的TLS CA证书
check_cert_validity "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem" "Org1 TLS CA证书"

# 检查Admin证书
check_cert_validity "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem" "Admin签名证书"

# 检查Admin私钥
check_file "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk" "Admin私钥文件"

# 检查Orderer的CA证书
check_cert_validity "$CONFIG_BASE/organizations/ordererOrganizations/example.com/msp/cacerts/ca.example.com-cert.pem" "Orderer CA证书"

# 检查Orderer的TLS CA证书
check_cert_validity "$CONFIG_BASE/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem" "Orderer TLS CA证书"

echo -e "${BLUE}📋 2. 验证证书信任链...${NC}"

# 验证Admin证书是否由Org1 CA签发
verify_cert_chain \
    "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem" \
    "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem" \
    "Admin证书信任链"

# 验证Peer TLS证书是否由Org1 TLS CA签发
verify_cert_chain \
    "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt" \
    "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" \
    "Peer TLS证书信任链"

echo -e "${BLUE}📋 3. 检查配置文件中的证书路径...${NC}"

# 检查connection-minimal.yaml中的路径
echo -e "${YELLOW}检查 connection-minimal.yaml 中的证书路径:${NC}"
grep -n "path:" "$CONFIG_BASE/connection-minimal.yaml" | while read line; do
    echo -e "${YELLOW}  $line${NC}"
    path=$(echo "$line" | sed 's/.*path: *//')
    if [[ -f "$CONFIG_BASE/$path" ]]; then
        echo -e "${GREEN}    ✅ 文件存在${NC}"
    else
        echo -e "${RED}    ❌ 文件不存在${NC}"
    fi
done

echo -e "${BLUE}📋 4. 创建证书信任链修复配置...${NC}"

# 创建修复后的配置文件
cat > "$CONFIG_BASE/connection-fixed-trust.yaml" << 'EOF'
---
name: test-network-fixed-trust
version: 1.0.0
client:
  organization: Org1
  connection:
    timeout:
      peer:
        endorser: '300'
      orderer: '300'
organizations:
  Org1:
    mspid: Org1MSP
    peers:
    - peer0.org1.example.com
peers:
  peer0.org1.example.com:
    url: grpcs://localhost:7051
    tlsCACerts:
      pem: |
EOF

# 读取并添加Org1的TLS CA证书
if [[ -f "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem" ]]; then
    cat "$CONFIG_BASE/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem" >> "$CONFIG_BASE/connection-fixed-trust.yaml"
    echo -e "${GREEN}✅ 已添加Org1 TLS CA证书到配置${NC}"
fi

cat >> "$CONFIG_BASE/connection-fixed-trust.yaml" << 'EOF'
    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      allow-insecure: false
orderers:
  orderer.example.com:
    url: grpcs://localhost:7050
    tlsCACerts:
      pem: |
EOF

# 读取并添加Orderer的TLS CA证书
if [[ -f "$CONFIG_BASE/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem" ]]; then
    cat "$CONFIG_BASE/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem" >> "$CONFIG_BASE/connection-fixed-trust.yaml"
    echo -e "${GREEN}✅ 已添加Orderer TLS CA证书到配置${NC}"
fi

cat >> "$CONFIG_BASE/connection-fixed-trust.yaml" << 'EOF'
    grpcOptions:
      ssl-target-name-override: orderer.example.com
      allow-insecure: false
channels:
  mychannel:
    orderers:
      - orderer.example.com
    peers:
      peer0.org1.example.com:
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
EOF

echo -e "${GREEN}✅ 已创建修复后的配置文件: connection-fixed-trust.yaml${NC}"

echo -e "${BLUE}📋 5. 创建证书池验证脚本...${NC}"

# 创建证书池验证脚本
cat > "$CONFIG_BASE/verify_cert_pool.sh" << 'EOF'
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
EOF

chmod +x "$CONFIG_BASE/verify_cert_pool.sh"
echo -e "${GREEN}✅ 已创建证书池验证脚本: verify_cert_pool.sh${NC}"

echo -e "${BLUE}📋 6. 生成修复建议...${NC}"

echo -e "${YELLOW}🔧 修复建议:${NC}"
echo -e "${YELLOW}1. 使用新创建的 connection-fixed-trust.yaml 配置文件${NC}"
echo -e "${YELLOW}2. 在连接代码中使用以下配置:${NC}"
echo -e "${BLUE}   config.FromFile(\"connection-fixed-trust.yaml\")${NC}"
echo -e "${YELLOW}3. 如果仍有问题，可以临时设置 allow-insecure: true 进行调试${NC}"
echo -e "${YELLOW}4. 运行证书池验证脚本: cd $CONFIG_BASE && ./verify_cert_pool.sh${NC}"

echo -e "${GREEN}�� 证书信任链检查完成！${NC}" 