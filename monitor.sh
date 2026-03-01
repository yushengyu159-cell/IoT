#!/bin/bash

# ESG VISA Fabric SDK 监控脚本
# 自动检测和重启掉线的服务

echo "🔍 ESG VISA Fabric SDK 监控脚本启动..."

# 检查服务是否运行
check_service() {
    if [ -f "sdk.pid" ]; then
        PID=$(cat sdk.pid)
        if ps -p $PID > /dev/null 2>&1; then
            return 0  # 服务运行中
        else
            return 1  # 服务已停止
        fi
    else
        return 1  # PID文件不存在
    fi
}

# 检查端口是否监听
check_port() {
    if lsof -Pi :8199 -sTCP:LISTEN -t >/dev/null 2>&1; then
        return 0  # 端口正在监听
    else
        return 1  # 端口未监听
    fi
}

# 检查Fabric网络
check_fabric() {
    if docker ps | grep -q "peer0.org1.example.com" && docker ps | grep -q "dev-peer.*basic"; then
        return 0  # Fabric网络正常
    else
        return 1  # Fabric网络异常
    fi
}

# 重启服务
restart_service() {
    echo "🔄 检测到服务异常，正在重启..."
    
    # 清理旧进程
    if [ -f "sdk.pid" ]; then
        OLD_PID=$(cat sdk.pid)
        kill -9 $OLD_PID 2>/dev/null || true
        rm -f sdk.pid
    fi
    
    # 清理端口
    if lsof -Pi :8199 -sTCP:LISTEN -t >/dev/null 2>&1; then
        kill -9 $(lsof -Pi :8199 -sTCP:LISTEN -t) 2>/dev/null || true
    fi
    
    # 重启Fabric网络（如果需要）
    if ! check_fabric; then
        echo "🔧 重启Fabric网络..."
        cd /root/home/go/fabric-samples/test-network
        ./network.sh down
        sleep 3
        ./network.sh up createChannel
        sleep 10
        ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go -ccv 1.0 -c mychannel
        sleep 15
        cd /root/home/go/fabric-sdk
    fi
    
    # 重启fabric-sdk服务
    echo "🚀 重启fabric-sdk服务..."
    export FABRIC_SDK_GO_LOG_LEVEL=FATAL
    export FABRIC_SDK_GO_MSP_VERIFY=false
    export GF_GERROR_BRIEF=true
    export GF_GERROR_STACK=false
    
    nohup ./fabric-sdk > sdk.log 2>&1 &
    PID=$!
    echo $PID > sdk.pid
    
    sleep 5
    
    if ps -p $PID > /dev/null; then
        echo "✅ 服务重启成功 (PID: $PID)"
    else
        echo "❌ 服务重启失败"
        tail -20 sdk.log
    fi
}

# 主监控循环
while true; do
    echo "$(date): 检查服务状态..."
    
    if ! check_service || ! check_port; then
        echo "⚠️  服务异常，正在处理..."
        restart_service
    else
        echo "✅ 服务正常运行"
    fi
    
    # 每30秒检查一次
    sleep 30
done

