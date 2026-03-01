#!/bin/bash

echo "🔍 查询链码背书策略 - ESG vs Basic"
echo "================================="

# 设置环境变量
export FABRIC_CFG_PATH=/home/ubuntu/go/fabric-samples/config
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051

echo "📋 1. 设置环境变量..."
echo "   - FABRIC_CFG_PATH: $FABRIC_CFG_PATH"
echo "   - CORE_PEER_LOCALMSPID: $CORE_PEER_LOCALMSPID"
echo "   - CORE_PEER_MSPCONFIGPATH: $CORE_PEER_MSPCONFIGPATH"
echo "   - CORE_PEER_ADDRESS: $CORE_PEER_ADDRESS"

echo ""
echo "📋 2. 查询ESG链码信息..."
echo "   🔍 查询ESG链码的背书策略:"
if peer lifecycle chaincode querycommitted -C mychannel --name esg --output json 2>/dev/null; then
    echo "   ✅ ESG链码查询成功"
else
    echo "   ❌ ESG链码查询失败"
    echo "   🔧 尝试查询所有已提交的链码:"
    peer lifecycle chaincode querycommitted -C mychannel --output json 2>/dev/null || echo "   查询失败"
fi

echo ""
echo "📋 3. 查询Basic链码信息..."
echo "   🔍 查询Basic链码的背书策略:"
if peer lifecycle chaincode querycommitted -C mychannel --name basic --output json 2>/dev/null; then
    echo "   ✅ Basic链码查询成功"
else
    echo "   ❌ Basic链码查询失败"
fi

echo ""
echo "📋 4. 检查通道配置..."
echo "   🔍 获取通道配置:"
if peer channel fetch config -c mychannel -o localhost:7050 --tls --cafile /home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem 2>/dev/null; then
    echo "   ✅ 通道配置获取成功"
else
    echo "   ❌ 通道配置获取失败"
fi

echo ""
echo "📋 5. 检查链码容器日志..."
echo "   🔍 ESG链码容器日志:"
docker logs --tail 10 dev-peer0.org1.example.com-esg_1.0-1caa6eb8c6120a023954e1baad60d50c54f0edfb7a6f3d95b4b340b1d44a70f1 2>/dev/null | grep -E "(ERROR|WARN|endorsement)" || echo "   无相关日志"

echo ""
echo "📋 6. 分析背书策略问题..."
echo "   💡 可能的背书策略问题:"
echo "   1. ESG链码的背书策略可能要求特定的组织组合"
echo "   2. 当前用户可能没有足够的权限"
echo "   3. 链码可能需要特定的背书策略配置"

echo ""
echo "📋 7. 生成修复建议..."
echo "   🔧 修复建议:"
echo "   1. 检查ESG链码的背书策略定义"
echo "      peer lifecycle chaincode querycommitted -C mychannel --name esg"
echo "   2. 如果需要，更新背书策略"
echo "      peer lifecycle chaincode approveformyorg -C mychannel --name esg --version 1.0 --package-id <package-id> --sequence 1 --signature-policy \"OR('Org1MSP.peer','Org2MSP.peer')\""
echo "   3. 提交链码定义更新"
echo "      peer lifecycle chaincode commit -C mychannel --name esg --version 1.0 --sequence 1 --signature-policy \"OR('Org1MSP.peer','Org2MSP.peer')\""

echo ""
echo "📊 链码背书策略查询总结"
echo "====================="
echo "✅ 设置了正确的环境变量"
echo "✅ 查询了ESG链码信息"
echo "✅ 查询了Basic链码信息"
echo "✅ 检查了通道配置"
echo "✅ 分析了链码容器日志"

echo ""
echo "🎯 关键发现:"
echo "   - ESG链码的背书策略可能是连接失败的原因"
echo "   - 需要检查ESG链码的具体背书策略要求"
echo "   - 可能需要更新背书策略以匹配当前配置"

echo ""
echo "💡 下一步建议:"
echo "   1. 检查ESG链码的背书策略定义"
echo "   2. 确认当前用户是否有足够的权限"
echo "   3. 如果需要，更新链码的背书策略" 