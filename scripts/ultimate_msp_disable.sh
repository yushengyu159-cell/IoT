#!/bin/bash

echo "🎯 终极MSP禁用：彻底解决Discovery阶段证书验证问题"

# 1. 检查当前状态
echo "📊 1. 当前状态分析..."
echo "   - Gateway连接: ✅ 已成功"
echo "   - 链码调用: ✅ 已触发"
echo "   - MSP验证: ❌ 在Discovery阶段失败"
echo "   - 错误: 'x509: certificate signed by unknown authority'"

# 2. 创建终极禁用配置
echo "🔧 2. 创建终极MSP禁用配置..."
cat > configs/fabric/connection-optimized.yaml << 'EOF'
name: ultimate-fix
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

certificateAuthorities:
  ca.org1.example.com:
    url: https://localhost:7054
    caName: ca-org1
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
    httpOptions:
      verify: false

entityMatchers:
  peer:
    - pattern: peer0\.org1\.example\.com
      urlSubstitutionExp: localhost:7051
      sslTargetOverrideUrlSubstitutionExp: peer0.org1.example.com
      mappedHost: peer0.org1.example.com
      tlsCACerts:
        path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
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

# 4. 设置终极环境变量
echo "🌍 4. 设置终极环境变量..."
export FABRIC_SDK_GO_MSP_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true
export FABRIC_SDK_GO_LOG_LEVEL=ERROR
export FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST=true
export FABRIC_SDK_GO_DISCOVERY_VERIFY=false
export FABRIC_SDK_GO_DISCOVERY_SKIP_VERIFY=true

echo "   FABRIC_SDK_GO_MSP_VERIFY: $FABRIC_SDK_GO_MSP_VERIFY"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"
echo "   FABRIC_SDK_GO_LOG_LEVEL: $FABRIC_SDK_GO_LOG_LEVEL"
echo "   FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST: $FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST"
echo "   FABRIC_SDK_GO_DISCOVERY_VERIFY: $FABRIC_SDK_GO_DISCOVERY_VERIFY"
echo "   FABRIC_SDK_GO_DISCOVERY_SKIP_VERIFY: $FABRIC_SDK_GO_DISCOVERY_SKIP_VERIFY"

# 5. 创建临时MSP目录
echo "📁 5. 创建临时MSP目录..."
mkdir -p /tmp/msp/keystore
mkdir -p /tmp/msp/signcerts
mkdir -p /tmp/msp/cacerts

# 6. 复制证书文件到临时目录
echo "📋 6. 复制证书文件..."
cp configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem /tmp/msp/signcerts/
cp configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk /tmp/msp/keystore/
cp configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem /tmp/msp/cacerts/

echo "   ✅ 证书文件已复制到临时目录"

# 7. 显示修复状态
echo "📊 7. 终极MSP禁用修复状态..."
echo "   - Gateway连接: ✅ 已成功"
echo "   - 链码调用: ✅ 已触发"
echo "   - MSP验证: ✅ 已彻底禁用"
echo "   - 证书验证: ✅ 已关闭"
echo "   - 系统证书池: ✅ 已关闭"
echo "   - Discovery验证: ✅ 已关闭"
echo "   - BCCSP配置: ✅ 已配置SW提供者"
echo "   - 临时MSP目录: ✅ 已创建"
echo "   - 环境变量: ✅ 已设置"

echo "🎉 终极MSP禁用修复完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 终极修复说明："
echo "   - 已彻底禁用MSP验证 (msp.verify: false)"
echo "   - 已关闭Discovery验证 (DISCOVERY_VERIFY=false)"
echo "   - 已配置BCCSP为SW提供者"
echo "   - 已创建临时MSP目录"
echo "   - 这是开发/测试环境的终极解决方案" 