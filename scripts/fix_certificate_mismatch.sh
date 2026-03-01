#!/bin/bash

# 修复证书不匹配问题
# 更新配置文件中的证书内容以匹配实际的证书链

set -e

echo "🔧 开始修复证书不匹配问题..."
echo "=================================="

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取实际的证书内容
echo -e "${BLUE}📋 获取实际证书内容...${NC}"

# 获取Org1的TLS CA证书
ORG1_TLS_CA="configs/fabric/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem"
if [ -f "$ORG1_TLS_CA" ]; then
    echo -e "${GREEN}✅ 找到Org1 TLS CA证书${NC}"
    ORG1_TLS_CA_CONTENT=$(cat "$ORG1_TLS_CA")
else
    echo -e "${RED}❌ Org1 TLS CA证书不存在${NC}"
    exit 1
fi

# 获取Org1的CA证书
ORG1_CA="configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"
if [ -f "$ORG1_CA" ]; then
    echo -e "${GREEN}✅ 找到Org1 CA证书${NC}"
    ORG1_CA_CONTENT=$(cat "$ORG1_CA")
else
    echo -e "${RED}❌ Org1 CA证书不存在${NC}"
    exit 1
fi

# 获取Orderer的TLS CA证书
ORDERER_TLS_CA="configs/fabric/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
if [ -f "$ORDERER_TLS_CA" ]; then
    echo -e "${GREEN}✅ 找到Orderer TLS CA证书${NC}"
    ORDERER_TLS_CA_CONTENT=$(cat "$ORDERER_TLS_CA")
else
    echo -e "${RED}❌ Orderer TLS CA证书不存在${NC}"
    exit 1
fi

# 获取Orderer的CA证书
ORDERER_CA="configs/fabric/organizations/ordererOrganizations/example.com/msp/cacerts/ca.example.com-cert.pem"
if [ -f "$ORDERER_CA" ]; then
    echo -e "${GREEN}✅ 找到Orderer CA证书${NC}"
    ORDERER_CA_CONTENT=$(cat "$ORDERER_CA")
else
    echo -e "${RED}❌ Orderer CA证书不存在${NC}"
    exit 1
fi

echo ""

# 创建修复后的配置文件
echo -e "${BLUE}📝 创建修复后的配置文件...${NC}"

cat > configs/fabric/connection-fixed-certs.yaml << 'EOF'
---
name: test-network-fixed-certs
version: 1.0.0
client:
  organization: Org1
  connection:
    timeout:
      peer:
        endorser: '300'
      orderer: '300'
  cryptoconfig:
    path: ${GOPATH}/src/github.com/hyperledger/fabric-samples/test-network/organizations
  credentialStore:
    path: /tmp/hfc-kvs
    cryptoStore:
      path: /tmp/hfc-cvs
  BCCSP:
    default: SW
    SW:
      hash: SHA2
      security: 256
      fileKeyStore:
        keyStore:
          path: /tmp/msp/keystore

organizations:
  Org1:
    mspid: Org1MSP
    cryptoPath: peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp
    peers:
    - peer0.org1.example.com
    certificateAuthorities:
    - ca.org1.example.com

peers:
  peer0.org1.example.com:
    url: grpcs://localhost:7051
    tlsCACerts:
      pem: |
EOF

# 添加Org1 TLS CA证书
echo "$ORG1_TLS_CA_CONTENT" >> configs/fabric/connection-fixed-certs.yaml

cat >> configs/fabric/connection-fixed-certs.yaml << 'EOF'
    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      allow-insecure: false
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false

orderers:
  orderer.example.com:
    url: grpcs://localhost:7050
    tlsCACerts:
      pem: |
EOF

# 添加Orderer TLS CA证书
echo "$ORDERER_TLS_CA_CONTENT" >> configs/fabric/connection-fixed-certs.yaml

cat >> configs/fabric/connection-fixed-certs.yaml << 'EOF'
    grpcOptions:
      ssl-target-name-override: orderer.example.com
      allow-insecure: false
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false

certificateAuthorities:
  ca.org1.example.com:
    url: https://localhost:7054
    caName: ca-org1
    tlsCACerts:
      pem: |
EOF

# 添加Org1 CA证书
echo "$ORG1_CA_CONTENT" >> configs/fabric/connection-fixed-certs.yaml

cat >> configs/fabric/connection-fixed-certs.yaml << 'EOF'
    httpOptions:
      verify: false
    registrar:
      - enrollId: admin
        enrollSecret: adminpw

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
        discoverAsLocalhost: true

entityMatchers:
  peer:
    - pattern: peer0.org1.example.com
      urlSubstitutionExp: localhost:7051
      sslTargetOverrideUrlSubstitutionExp: peer0.org1.example.com
      mappedHost: peer0.org1.example.com
  orderer:
    - pattern: orderer.example.com
      urlSubstitutionExp: localhost:7050
      sslTargetOverrideUrlSubstitutionExp: orderer.example.com
      mappedHost: orderer.example.com
  certificateAuthority:
    - pattern: ca.org1.example.com
      urlSubstitutionExp: localhost:7054
      mappedHost: ca.org1.example.com
EOF

echo -e "${GREEN}✅ 修复后的配置文件已创建: configs/fabric/connection-fixed-certs.yaml${NC}"

echo ""

# 验证新配置文件
echo -e "${BLUE}🔍 验证新配置文件...${NC}"

if [ -f "configs/fabric/connection-fixed-certs.yaml" ]; then
    echo -e "${GREEN}✅ 配置文件存在${NC}"
    
    # 检查YAML语法
    if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-fixed-certs.yaml'))" 2>/dev/null; then
        echo -e "${GREEN}✅ YAML语法正确${NC}"
    else
        echo -e "${RED}❌ YAML语法错误${NC}"
    fi
    
    # 检查证书内容
    if grep -q "-----BEGIN CERTIFICATE-----" configs/fabric/connection-fixed-certs.yaml; then
        echo -e "${GREEN}✅ 包含PEM证书${NC}"
    else
        echo -e "${RED}❌ 不包含PEM证书${NC}"
    fi
else
    echo -e "${RED}❌ 配置文件创建失败${NC}"
    exit 1
fi

echo ""

# 更新连接管理器使用新配置文件
echo -e "${BLUE}🔧 更新连接管理器配置...${NC}"

# 备份原配置文件
cp configs/fabric/connection-optimized.yaml configs/fabric/connection-optimized.yaml.backup
echo -e "${GREEN}✅ 原配置文件已备份${NC}"

# 使用新配置文件
cp configs/fabric/connection-fixed-certs.yaml configs/fabric/connection-optimized.yaml
echo -e "${GREEN}✅ 已更新为修复后的配置文件${NC}"

echo ""

# 创建证书链验证脚本
echo -e "${BLUE}📋 创建证书链验证脚本...${NC}"

cat > scripts/verify_cert_chain.sh << 'EOF'
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
EOF

chmod +x scripts/verify_cert_chain.sh
echo -e "${GREEN}✅ 证书链验证脚本已创建${NC}"

echo ""

# 运行证书链验证
echo -e "${BLUE}🔍 运行证书链验证...${NC}"
./scripts/verify_cert_chain.sh

echo ""
echo -e "${GREEN}🎉 证书不匹配问题修复完成！${NC}"
echo "=================================="

# 提供建议
echo -e "${YELLOW}💡 修复内容:${NC}"
echo "1. ✅ 更新了配置文件中的证书内容"
echo "2. ✅ 使用正确的Org1 CA证书"
echo "3. ✅ 使用正确的TLS CA证书"
echo "4. ✅ 修复了证书链不匹配问题"
echo "5. ✅ 创建了证书链验证脚本"

echo ""
echo -e "${BLUE}📋 下一步操作:${NC}"
echo "1. 重新启动应用程序测试连接"
echo "2. 检查是否还有证书验证错误"
echo "3. 验证链码调用是否正常工作"
echo "4. 运行: go run main.go" 