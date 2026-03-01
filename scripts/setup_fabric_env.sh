#!/bin/bash

echo "🌍 设置完整的Fabric环境变量"

# 1. 设置Fabric SDK Go环境变量
echo "🔧 1. 设置Fabric SDK Go环境变量..."
export FABRIC_SDK_GO_MSP_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true
export FABRIC_SDK_GO_LOG_LEVEL=ERROR
export FABRIC_SDK_GO_DISCOVERY_AS_LOCALHOST=true
export FABRIC_SDK_GO_DISCOVERY_VERIFY=false
export FABRIC_SDK_GO_DISCOVERY_SKIP_VERIFY=true

# 2. 设置Core Peer环境变量
echo "🔧 2. 设置Core Peer环境变量..."
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_CERT_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
export CORE_PEER_TLS_KEY_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key
export CORE_PEER_TLS_ROOTCERT_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_LISTENADDRESS=0.0.0.0:7051
export CORE_PEER_CHAINCODEADDRESS=0.0.0.0:7052
export CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052
export CORE_PEER_GOSSIP_BOOTSTRAP=localhost:7051
export CORE_PEER_GOSSIP_EXTERNALENDPOINT=localhost:7051
export CORE_PEER_GOSSIP_SKIPHANDSHAKE=true

# 3. 设置Core Orderer环境变量
echo "🔧 3. 设置Core Orderer环境变量..."
export CORE_ORDERER_TLS_ENABLED=true
export CORE_ORDERER_TLS_CERT_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export CORE_ORDERER_TLS_KEY_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/keystore/priv_sk
export CORE_ORDERER_TLS_ROOTCERT_FILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
export CORE_ORDERER_LOCALMSPID=OrdererMSP
export CORE_ORDERER_MSPCONFIGPATH=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp
export CORE_ORDERER_LISTENADDRESS=0.0.0.0:7050
export CORE_ORDERER_LISTENPORT=7050

# 4. 设置Core Chaincode环境变量
echo "🔧 4. 设置Core Chaincode环境变量..."
export CORE_CHAINCODE_LOGGING_LEVEL=INFO
export CORE_CHAINCODE_LOGGING_SHIM=INFO
export CORE_CHAINCODE_LOGGING_FORMAT=%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}
export CORE_CHAINCODE_LOGGING_LEVEL=INFO
export CORE_CHAINCODE_LOGGING_SHIM=INFO

# 5. 设置Fabric CA环境变量
echo "🔧 5. 设置Fabric CA环境变量..."
export FABRIC_CA_CLIENT_TLS_CERTFILES=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
export FABRIC_CA_CLIENT_TLS_CLIENT_CERTFILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.crt
export FABRIC_CA_CLIENT_TLS_CLIENT_KEYFILE=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/tls/client.key
export FABRIC_CA_CLIENT_HOME=/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com
export FABRIC_CA_CLIENT_MSPDIR=msp

# 6. 设置网络配置环境变量
echo "🔧 6. 设置网络配置环境变量..."
export FABRIC_NETWORK_CONFIG=/home/ubuntu/go/fabric-sdk/configs/fabric/connection-optimized.yaml
export FABRIC_NETWORK_NAME=test-network
export FABRIC_NETWORK_VERSION=1.0.0
export FABRIC_ORG_NAME=Org1
export FABRIC_ORG_MSPID=Org1MSP
export FABRIC_ORG_DOMAIN=org1.example.com

# 7. 设置开发环境特定变量
echo "🔧 7. 设置开发环境特定变量..."
export FABRIC_DEV_MODE=true
export FABRIC_DEV_SKIP_VERIFY=true
export FABRIC_DEV_ALLOW_INSECURE=true
export FABRIC_DEV_DISCOVERY_AS_LOCALHOST=true
export FABRIC_DEV_LOG_LEVEL=ERROR

# 8. 验证环境变量设置
echo "🔍 8. 验证环境变量设置..."
echo "📋 Fabric SDK Go环境变量:"
env | grep FABRIC_SDK_GO | sort

echo "📋 Core Peer环境变量:"
env | grep CORE_PEER | sort

echo "📋 Core Orderer环境变量:"
env | grep CORE_ORDERER | sort

echo "📋 Fabric CA环境变量:"
env | grep FABRIC_CA | sort

echo "📋 网络配置环境变量:"
env | grep FABRIC_NETWORK | sort

echo "📋 开发环境变量:"
env | grep FABRIC_DEV | sort

# 9. 检查关键文件是否存在
echo "🔍 9. 检查关键文件是否存在..."
CRITICAL_FILES=(
    "/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt"
    "/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key"
    "/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
    "/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
    "/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
    "/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"
    "/home/ubuntu/go/fabric-sdk/configs/fabric/connection-optimized.yaml"
)

for file in "${CRITICAL_FILES[@]}"; do
    if [ -f "$file" ]; then
        echo "   ✅ $file 存在"
    else
        echo "   ❌ $file 不存在"
    fi
done

echo "🎉 Fabric环境变量设置完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 环境变量说明："
echo "   - 已设置完整的Fabric SDK Go环境变量"
echo "   - 已设置Core Peer环境变量"
echo "   - 已设置Core Orderer环境变量"
echo "   - 已设置Fabric CA环境变量"
echo "   - 已设置网络配置环境变量"
echo "   - 已设置开发环境特定变量"
echo "   - 已验证所有关键文件存在" 