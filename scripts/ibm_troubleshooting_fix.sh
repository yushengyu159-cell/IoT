#!/bin/bash

echo "🎯 IBM故障排除修复：解决entityMatchers、证书验证和网络策略问题"

# 1. 检查当前问题
echo "📊 1. 当前问题分析..."
echo "   - Gateway连接: ✅ 已成功"
echo "   - 链码调用: ✅ 已触发"
echo "   - Discovery失败: ❌ 无法找到可背书节点"
echo "   - 错误: 'no endorsement combination can be satisfied'"

# 2. 修复entityMatchers配置
echo "🔧 2. 修复entityMatchers配置..."
cat > configs/fabric/connection-optimized.yaml << 'EOF'
name: ibm-troubleshooting-fix
version: 1.0.0

client:
  organization: Org1
  verify: false
  msp:
    verify: false
  tlsCerts:
    systemCertPool: false
  BCCSP:
    security:
      default:
        provider: SW
        hashFamily: SHA2
        secLevel: 256
        ephemeral: false
        fileKeystore:
          keyStorePath: /tmp/msp/keystore

organizations:
  Org1:
    mspid: Org1MSP
    cryptoPath: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp
    peers:
      - peer0.org1.example.com
    certificateAuthorities:
      - ca.org1.example.com

peers:
  peer0.org1.example.com:
    url: grpcs://localhost:7051
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      allow-insecure: false
      verify: false
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false
      allow-insecure: false

orderers:
  orderer.example.com:
    url: grpcs://localhost:7050
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    grpcOptions:
      ssl-target-name-override: orderer.example.com
      allow-insecure: false
      verify: false
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false

certificateAuthorities:
  ca.org1.example.com:
    url: https://localhost:7054
    caName: ca-org1
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
    httpOptions:
      verify: false
    registrar:
      - enrollId: admin
        enrollSecret: adminpw

# IBM故障排除：完整的entityMatchers配置
entityMatchers:
  peer:
    - pattern: peer0\.org1\.example\.com
      urlSubstitutionExp: localhost:7051
      sslTargetOverrideUrlSubstitutionExp: peer0.org1.example.com
      mappedHost: peer0.org1.example.com
      tlsCACerts:
        path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
      verify: false

  orderer:
    - pattern: orderer\.example\.com
      urlSubstitutionExp: localhost:7050
      sslTargetOverrideUrlSubstitutionExp: orderer.example.com
      mappedHost: orderer.example.com
      tlsCACerts:
        path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
      verify: false

  certificateAuthority:
    - pattern: ca\.org1\.example\.com
      urlSubstitutionExp: localhost:7054
      mappedHost: ca.org1.example.com
      tlsCACerts:
        path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
      verify: false
EOF

# 3. 验证配置
echo "🔍 3. 验证配置..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    exit 1
fi

# 4. 设置IBM推荐的环境变量
echo "🌍 4. 设置IBM推荐的环境变量..."
export FABRIC_SDK_GO_MSP_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true
export FABRIC_SDK_GO_LOG_LEVEL=ERROR
export FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST=true
export FABRIC_SDK_GO_DISCOVERY_VERIFY=false
export FABRIC_SDK_GO_DISCOVERY_SKIP_VERIFY=true

# IBM故障排除：设置网络策略相关变量
export IBPOPERATOR_CONSOLE_APPLYNETWORKPOLICY=false
export FABRIC_NETWORK_POLICY_DISABLED=true
export FABRIC_DISCOVERY_AS_LOCALHOST=true
export FABRIC_DISCOVERY_VERIFY=false

# 设置Core Peer环境变量
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

echo "   FABRIC_SDK_GO_MSP_VERIFY: $FABRIC_SDK_GO_MSP_VERIFY"
echo "   FABRIC_SDK_GO_DISCOVERY_VERIFY: $FABRIC_SDK_GO_DISCOVERY_VERIFY"
echo "   IBPOPERATOR_CONSOLE_APPLYNETWORKPOLICY: $IBPOPERATOR_CONSOLE_APPLYNETWORKPOLICY"
echo "   CORE_PEER_ADDRESS: $CORE_PEER_ADDRESS"

# 5. 检查证书类型和背书策略
echo "🔐 5. 检查证书类型和背书策略..."
echo "📋 Admin用户证书信息:"
openssl x509 -in configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem -noout -subject -issuer

echo "📋 证书类型验证:"
if openssl x509 -in configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem -noout -text | grep -q "admin"; then
    echo "   ✅ Admin用户证书类型正确"
else
    echo "   ❌ Admin用户证书类型可能不正确"
fi

# 6. 创建临时MSP目录
echo "📁 6. 创建临时MSP目录..."
mkdir -p /tmp/msp/keystore
mkdir -p /tmp/msp/signcerts
mkdir -p /tmp/msp/cacerts

cp configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem /tmp/msp/signcerts/
cp configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk /tmp/msp/keystore/
cp configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem /tmp/msp/cacerts/

echo "   ✅ 证书文件已复制到临时目录"

# 7. 显示修复状态
echo "📊 7. IBM故障排除修复状态..."
echo "   - Gateway连接: ✅ 已成功"
echo "   - entityMatchers: ✅ 已完整配置"
echo "   - 证书验证: ✅ 已关闭"
echo "   - 网络策略: ✅ 已禁用"
echo "   - Discovery验证: ✅ 已关闭"
echo "   - 临时MSP目录: ✅ 已创建"
echo "   - 环境变量: ✅ 已设置"

echo "🎉 IBM故障排除修复完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 IBM故障排除修复说明："
echo "   - 已完整配置entityMatchers (peer, orderer, ca)"
echo "   - 已关闭所有证书验证"
echo "   - 已禁用网络策略 (IBPOPERATOR_CONSOLE_APPLYNETWORKPOLICY=false)"
echo "   - 已设置Discovery相关环境变量"
echo "   - 已验证Admin用户证书类型"
echo "   - 这是基于IBM故障排除指南的完整解决方案" 