#!/bin/bash

# ESG VISA Fabric SDK 稳定启动脚本
# 解决掉线问题的专业方案

echo "🚀 启动 ESG VISA Fabric SDK 稳定版本..."

# 1. 检查并清理端口占用
echo "🔍 检查端口占用..."
if lsof -Pi :8199 -sTCP:LISTEN -t >/dev/null ; then
    echo "⚠️  端口8199被占用，正在清理..."
    kill -9 $(lsof -Pi :8199 -sTCP:LISTEN -t) 2>/dev/null || true
    sleep 2
fi

# 2. 检查Fabric网络状态
echo "🔍 检查Fabric网络状态..."
cd /root/home/go/fabric-samples/test-network

# 检查peer节点是否运行
if ! docker ps | grep -q "peer0.org1.example.com"; then
    echo "❌ Fabric网络未运行，正在启动..."
    ./network.sh up createChannel
    sleep 10
fi

# 检查链码是否部署
if ! docker ps | grep -q "dev-peer.*basic"; then
    echo "❌ 链码未部署，正在部署..."
    ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go -ccv 1.0 -c mychannel
    sleep 15
fi

# 3. 检查IPFS节点
echo "🔍 检查IPFS节点..."
if ! docker ps | grep -q "fabric-sdk-ipfs"; then
    echo "❌ IPFS节点未运行，正在启动..."
    docker start fabric-sdk-ipfs
    sleep 5
fi

# 4. 检查MySQL数据库
echo "🔍 检查MySQL数据库..."
if ! docker ps | grep -q "fabric-sdk-mysql"; then
    echo "❌ MySQL数据库未运行，正在启动..."
    docker start fabric-sdk-mysql
    sleep 5
fi

# 5. 设置环境变量
echo "🔧 设置环境变量..."
export FABRIC_SDK_GO_LOG_LEVEL=FATAL
export FABRIC_SDK_GO_MSP_VERIFY=false
export GF_GERROR_BRIEF=true
export GF_GERROR_STACK=false

# 6. 启动fabric-sdk服务
echo "🚀 启动fabric-sdk服务..."
cd /root/home/go/fabric-sdk

# 使用nohup后台运行，并重定向日志
nohup ./fabric-sdk > sdk.log 2>&1 &
PID=$!

# 等待服务启动
sleep 5

# 检查服务是否启动成功
if ps -p $PID > /dev/null; then
    echo "✅ fabric-sdk服务启动成功 (PID: $PID)"
    echo "📊 服务状态:"
    echo "   - 端口: 8199"
    echo "   - API文档: http://localhost:8199/swagger/"
    echo "   - 前端界面: http://localhost:8199/static/"
    echo "   - 日志文件: sdk.log"
    
    # 保存PID到文件
    echo $PID > sdk.pid
    echo "💾 PID已保存到 sdk.pid"
else
    echo "❌ fabric-sdk服务启动失败"
    echo "📋 查看日志:"
    tail -20 sdk.log
    exit 1
fi

echo "🎉 启动完成！"

