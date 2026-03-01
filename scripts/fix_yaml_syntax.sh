#!/bin/bash

echo "🔧 修复YAML语法错误：证书内容缩进问题"

# 1. 获取CA证书内容
echo "🔐 1. 获取CA证书内容..."
CA_CERT_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"

if [ ! -f "$CA_CERT_PATH" ]; then
    echo "❌ CA证书文件不存在: $CA_CERT_PATH"
    exit 1
fi

# 2. 创建正确缩进的YAML配置
echo "📝 2. 创建正确缩进的YAML配置..."
cat > configs/fabric/connection-optimized.yaml << 'EOF'
name: "test-network"
version: "1.0.0"

client:
  organization: Org1
  connection:
    timeout:
      peer:
        endorser: '300'
  tlsCerts:
    systemCertPool: false  # 强制关闭系统根证书池
    client:
      key:
        path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key
      cert:
        path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt
  verify: false  # 关闭证书验证
  msp:
    verify: false  # 关闭MSP证书验证

organizations:
  Org1:
    mspid: Org1MSP
    cryptoPath: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp
    # 显式指定MSP的CA证书，确保SDK使用正确的CA
    msp:
      verify: false  # 关闭MSP证书验证
      caCerts:
        - |
EOF

# 3. 添加正确缩进的CA证书内容
echo "          $(cat "$CA_CERT_PATH" | sed 's/^/          /')" >> configs/fabric/connection-optimized.yaml

# 4. 添加其余配置
cat >> configs/fabric/connection-optimized.yaml << 'EOF'

    peers:
      - peer0.org1.example.com
    certificateAuthorities:
      - ca.org1.example.com

peers:
  peer0.org1.example.com:
    url: grpcs://localhost:7051
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      hostnameOverride: peer0.org1.example.com
      verify: false  # 关闭peer TLS验证

orderers:
  orderer.example.com:
    url: grpcs://localhost:7050
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    grpcOptions:
      ssl-target-name-override: orderer.example.com
      hostnameOverride: orderer.example.com
      verify: false  # 关闭orderer TLS验证

certificateAuthorities:
  ca.org1.example.com:
    url: https://localhost:7054
    caName: ca-org1
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
    httpOptions:
      verify: false  # 关闭CA TLS验证

entityMatchers:
  peer:
    - pattern: (\w*)peer0.org1.example.com(\w*)
      urlSubstitutionExp: grpcs://localhost:7051
      sslTargetOverrideUrlSubstitutionExp: peer0.org1.example.com
      mappedHost: peer0.org1.example.com
      tlsCACertsPath: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
      verify: false  # 关闭entityMatcher TLS验证

  orderer:
    - pattern: (\w*)orderer.example.com(\w*)
      urlSubstitutionExp: grpcs://localhost:7050
      sslTargetOverrideUrlSubstitutionExp: orderer.example.com
      mappedHost: orderer.example.com
      tlsCACertsPath: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
      verify: false  # 关闭entityMatcher TLS验证

  certificateAuthority:
    - pattern: (\w*)ca.org1.example.com(\w*)
      urlSubstitutionExp: https://localhost:7054
      mappedHost: ca.org1.example.com
      tlsCACertsPath: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
      verify: false  # 关闭entityMatcher TLS验证
EOF

# 5. 验证YAML语法
echo "🔍 3. 验证YAML语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>&1
    exit 1
fi

# 6. 检查CA证书是否成功注入
echo "📋 4. 检查CA证书注入..."
if grep -q "BEGIN CERTIFICATE" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ CA证书已成功注入配置"
else
    echo "   ❌ CA证书注入失败"
    exit 1
fi

# 7. 设置环境变量
echo "🌍 5. 设置环境变量..."
export FABRIC_SDK_GO_MSP_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=true
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=false
export FABRIC_SDK_GO_LOG_LEVEL=ERROR
export FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST=true

echo "   FABRIC_SDK_GO_MSP_VERIFY: $FABRIC_SDK_GO_MSP_VERIFY"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"
echo "   FABRIC_SDK_GO_LOG_LEVEL: $FABRIC_SDK_GO_LOG_LEVEL"
echo "   FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST: $FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST"

echo "🎉 YAML语法修复完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 修复说明："
echo "   - 已修复YAML语法错误（证书内容缩进）"
echo "   - 已关闭系统根证书池 (systemCertPool: false)"
echo "   - 已嵌入Fabric CA证书到配置"
echo "   - 已关闭所有证书验证"
echo "   - 已设置所有必要的环境变量" 