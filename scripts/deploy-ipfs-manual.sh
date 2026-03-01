#!/bin/bash

echo "🔧 手动部署IPFS服务（调试模式）..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker服务"
    exit 1
fi

# 清理现有容器
echo "🧹 清理现有IPFS容器..."
docker stop fabric-sdk-ipfs 2>/dev/null || true
docker rm -f fabric-sdk-ipfs 2>/dev/null || true

# 清理网络
echo "🌐 清理IPFS网络..."
docker network rm fabric-sdk_ipfs-network 2>/dev/null || true

# 创建网络
echo "🌐 创建IPFS网络..."
docker network create fabric-sdk_ipfs-network

# 创建数据卷
echo "💾 创建IPFS数据卷..."
docker volume create fabric-sdk_ipfs_data
docker volume create fabric-sdk_ipfs_staging

# 启动IPFS容器（基础模式）
echo "🚀 启动IPFS容器（基础模式）..."
docker run -d \
  --name fabric-sdk-ipfs \
  --network fabric-sdk_ipfs-network \
  -p 4001:4001 \
  -p 5001:5001 \
  -p 8081:8080 \
  -v fabric-sdk_ipfs_data:/data/ipfs \
  -v fabric-sdk_ipfs_staging:/export \
  -e IPFS_PROFILE=server \
  -e IPFS_PATH=/data/ipfs \
  ipfs/kubo:latest

# 检查容器状态
echo "🔍 检查容器状态..."
if docker ps | grep -q "fabric-sdk-ipfs"; then
    echo "✅ IPFS容器启动成功！"
    
    # 等待容器完全启动
    echo "⏳ 等待容器完全启动..."
    sleep 10
    
    # 检查IPFS进程
    echo "🔍 检查IPFS进程..."
    if docker exec fabric-sdk-ipfs ps aux | grep -q "ipfs daemon"; then
        echo "✅ IPFS守护进程运行正常！"
    else
        echo "⚠️ IPFS守护进程未运行，尝试手动启动..."
        docker exec fabric-sdk-ipfs ipfs daemon &
        sleep 5
    fi
    
    # 测试IPFS功能
    echo "🧪 测试IPFS功能..."
    if docker exec fabric-sdk-ipfs ipfs id > /dev/null 2>&1; then
        echo "✅ IPFS节点功能正常！"
        
        # 获取节点信息
        echo "📊 IPFS节点信息："
        docker exec fabric-sdk-ipfs ipfs id | grep -E "(ID|Addresses|ProtocolVersion)"
        
        # 配置CORS
        echo "⚙️ 配置IPFS CORS..."
        docker exec fabric-sdk-ipfs ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]'
        docker exec fabric-sdk-ipfs ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "POST", "GET"]'
        docker exec fabric-sdk-ipfs ipfs config --json API.HTTPHeaders.Access-Control-Allow-Headers '["Authorization"]'
        
        echo ""
        echo "🌐 服务访问地址："
        echo "   - IPFS API: http://localhost:5001"
        echo "   - IPFS Gateway: http://localhost:8081"
        echo "   - P2P端口: 4001"
        
        echo ""
        echo "🔗 测试API连接..."
        if curl -s http://localhost:5001/api/v0/version > /dev/null; then
            echo "✅ IPFS API连接成功！"
        else
            echo "❌ IPFS API连接失败"
        fi
        
    else
        echo "❌ IPFS节点功能异常"
        echo "📋 容器日志："
        docker logs fabric-sdk-ipfs
    fi
    
else
    echo "❌ IPFS容器启动失败"
    echo "📋 容器日志："
    docker logs fabric-sdk-ipfs
    exit 1
fi

echo ""
echo "🎉 手动部署完成！"
echo "💡 如果遇到问题，可以查看容器日志："
echo "   docker logs fabric-sdk-ipfs"




