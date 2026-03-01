#!/bin/bash

echo "🚀 开始部署IPFS分布式存储服务..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker服务"
    exit 1
fi

# 停止并删除现有IPFS容器
echo "🔄 清理现有IPFS容器..."
docker-compose -f docker-compose.ipfs.yml down -v

# 清理可能存在的IPFS容器
echo "🧹 清理残留的IPFS容器..."
docker rm -f fabric-sdk-ipfs fabric-sdk-ipfs-cluster fabric-sdk-ipfs-webui 2>/dev/null || true

# 清理IPFS网络
echo "🌐 清理IPFS网络..."
docker network rm fabric-sdk_ipfs-network 2>/dev/null || true

# 拉取IPFS镜像
echo "📥 拉取IPFS镜像..."
docker pull ipfs/kubo:latest

# 检查镜像拉取是否成功
if [ $? -ne 0 ]; then
    echo "❌ IPFS镜像拉取失败，尝试使用国内镜像源..."
    # 可以在这里添加国内镜像源配置
    exit 1
fi

# 启动IPFS服务
echo "🚀 启动IPFS服务..."
docker-compose -f docker-compose.ipfs.yml up -d

# 检查服务启动状态
if [ $? -ne 0 ]; then
    echo "❌ IPFS服务启动失败，尝试使用简化配置..."
    echo "🔄 使用简化配置重新部署..."
    docker-compose -f docker-compose.ipfs-simple.yml up -d
    if [ $? -ne 0 ]; then
        echo "❌ 简化配置也启动失败，尝试稳定配置..."
        echo "🔄 使用稳定配置重新部署..."
        docker-compose -f docker-compose.ipfs-stable.yml up -d
        if [ $? -ne 0 ]; then
            echo "❌ 所有配置都启动失败"
            exit 1
        fi
    fi
fi

# 等待IPFS节点启动
echo "⏳ 等待IPFS节点启动..."
sleep 30

# 检查IPFS节点状态
echo "🔍 检查IPFS节点状态..."
if docker exec fabric-sdk-ipfs ipfs id > /dev/null 2>&1; then
    echo "✅ IPFS节点启动成功！"
    
    # 获取节点信息
    echo "📊 IPFS节点信息："
    docker exec fabric-sdk-ipfs ipfs id | grep -E "(ID|Addresses|ProtocolVersion)"
    
    echo ""
    echo "🌐 服务访问地址："
    echo "   - IPFS API: http://localhost:5001"
    echo "   - IPFS Gateway: http://localhost:8081"
    
    # 检查是否启动了完整版本
    if docker ps | grep -q "fabric-sdk-ipfs-webui"; then
        echo "   - IPFS Web UI: http://localhost:3001"
        echo "   - IPFS Cluster API: http://localhost:9094"
        echo "   - IPFS Cluster Proxy: http://localhost:9095"
    fi
    
    echo ""
    echo "🔗 测试IPFS连接..."
    if curl -s http://localhost:5001/api/v0/version > /dev/null; then
        echo "✅ IPFS API连接成功！"
    else
        echo "❌ IPFS API连接失败"
    fi
    
    # 初始化IPFS配置
    echo "⚙️ 初始化IPFS配置..."
    docker exec fabric-sdk-ipfs ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '["*"]'
    docker exec fabric-sdk-ipfs ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "POST", "GET"]'
    docker exec fabric-sdk-ipfs ipfs config --json API.HTTPHeaders.Access-Control-Allow-Headers '["Authorization"]'
    
    echo "✅ IPFS配置初始化完成！"
    
else
    echo "❌ IPFS节点启动失败，请检查日志："
    docker logs fabric-sdk-ipfs
    exit 1
fi

echo ""
echo "🎉 IPFS分布式存储部署完成！"
echo "💡 现在可以启动fabric-sdk服务，测试IPFS功能"
echo ""
echo "📝 测试命令："
echo "   # 测试IPFS API"
echo "   curl http://localhost:5001/api/v0/version"
echo ""
echo "   # 测试文件上传"
echo "   curl -X POST -F file=@test.txt http://localhost:5001/api/v0/add"
echo ""
     echo "   # 访问Web UI"
     echo "   open http://localhost:3001"
