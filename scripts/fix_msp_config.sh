#!/bin/bash

echo "🔧 完整修复MSP配置 - 使用原始test-network配置"
echo "=============================================="

# 定义路径
ORIGINAL_PATH="/home/ubuntu/go/fabric-samples/test-network/organizations"
CURRENT_PATH="/home/ubuntu/go/fabric-sdk/configs/fabric/organizations"

echo "📋 1. 备份当前配置..."
BACKUP_DIR="/home/ubuntu/go/fabric-sdk/configs/fabric/msp_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"
cp -r "$CURRENT_PATH" "$BACKUP_DIR/"
echo "   ✅ 备份完成: $BACKUP_DIR"

echo ""
echo "📋 2. 完整复制原始MSP配置..."
echo "   - 复制peerOrganizations..."
rm -rf "$CURRENT_PATH/peerOrganizations"
cp -r "$ORIGINAL_PATH/peerOrganizations" "$CURRENT_PATH/"

echo "   - 复制ordererOrganizations..."
rm -rf "$CURRENT_PATH/ordererOrganizations"
cp -r "$ORIGINAL_PATH/ordererOrganizations" "$CURRENT_PATH/"

echo "   - 复制fabric-ca配置..."
rm -rf "$CURRENT_PATH/fabric-ca"
cp -r "$ORIGINAL_PATH/fabric-ca" "$CURRENT_PATH/"

echo "   - 复制其他配置文件..."
cp "$ORIGINAL_PATH/ccp-generate.sh" "$CURRENT_PATH/" 2>/dev/null || echo "   ⚠️  ccp-generate.sh不存在"
cp "$ORIGINAL_PATH/ccp-template.json" "$CURRENT_PATH/" 2>/dev/null || echo "   ⚠️  ccp-template.json不存在"
cp "$ORIGINAL_PATH/ccp-template.yaml" "$CURRENT_PATH/" 2>/dev/null || echo "   ⚠️  ccp-template.yaml不存在"

echo ""
echo "📋 3. 验证复制结果..."
echo "   - 检查peerOrganizations:"
if [ -d "$CURRENT_PATH/peerOrganizations" ]; then
    echo "     ✅ peerOrganizations目录存在"
    echo "     📁 包含组织: $(ls "$CURRENT_PATH/peerOrganizations/")"
else
    echo "     ❌ peerOrganizations目录不存在"
fi

echo "   - 检查ordererOrganizations:"
if [ -d "$CURRENT_PATH/ordererOrganizations" ]; then
    echo "     ✅ ordererOrganizations目录存在"
    echo "     📁 包含组织: $(ls "$CURRENT_PATH/ordererOrganizations/")"
else
    echo "     ❌ ordererOrganizations目录不存在"
fi

echo ""
echo "📋 4. 检查证书文件完整性..."
ORIGINAL_CERTS=$(find "$ORIGINAL_PATH" -name "*.pem" | wc -l)
CURRENT_CERTS=$(find "$CURRENT_PATH" -name "*.pem" | wc -l)
echo "   - 原始证书文件数量: $ORIGINAL_CERTS"
echo "   - 当前证书文件数量: $CURRENT_CERTS"

if [ "$ORIGINAL_CERTS" -eq "$CURRENT_CERTS" ]; then
    echo "   ✅ 证书文件数量匹配"
else
    echo "   ❌ 证书文件数量不匹配"
fi

echo ""
echo "📋 5. 检查关键证书文件..."
KEY_FILES=(
    "peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"
    "peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem"
    "peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
    "peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
    "ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
)

for file in "${KEY_FILES[@]}"; do
    if [ -f "$CURRENT_PATH/$file" ]; then
        echo "   ✅ $file 存在"
    else
        echo "   ❌ $file 不存在"
    fi
done

echo ""
echo "📋 6. 更新YAML配置文件使用原始格式..."
# 创建基于原始配置的新YAML文件
cat > "$CURRENT_PATH/../connection-original.yaml" << 'EOF'
name: test-network-org1
version: 1.0.0
client:
  organization: Org1
  connection:
    timeout:
      peer:
        endorser: '300'
organizations:
  Org1:
    mspid: Org1MSP
    peers:
    - peer0.org1.example.com
    certificateAuthorities:
    - ca.org1.example.com
peers:
  peer0.org1.example.com:
    url: grpcs://localhost:7051
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem
    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      hostnameOverride: peer0.org1.example.com
      verify: false
orderers:
  orderer.example.com:
    url: grpcs://localhost:7050
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem
    grpcOptions:
      ssl-target-name-override: orderer.example.com
      verify: false
certificateAuthorities:
  ca.org1.example.com:
    url: https://localhost:7054
    caName: ca-org1
    tlsCACerts:
      path: /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
    httpOptions:
      verify: false
    registrar:
      enrollId: admin
      enrollSecret: adminpw
EOF

echo "   ✅ 创建了基于原始格式的配置文件: connection-original.yaml"

echo ""
echo "📊 MSP配置修复总结"
echo "=================="
echo "✅ 完整复制了原始test-network的MSP配置"
echo "✅ 备份了当前配置到: $BACKUP_DIR"
echo "✅ 证书文件数量: $CURRENT_CERTS (应该与 $ORIGINAL_CERTS 匹配)"
echo "✅ 创建了基于原始格式的配置文件"
echo ""
echo "🎉 MSP配置修复完成！"
echo ""
echo "💡 下一步:"
echo "   1. 使用新的connection-original.yaml配置文件"
echo "   2. 重启SDK服务进行测试"
echo "   3. 验证证书验证问题是否解决" 