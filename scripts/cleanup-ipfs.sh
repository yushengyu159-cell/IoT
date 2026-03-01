#!/bin/bash

echo "🧹 开始清理IPFS相关资源..."

# 停止所有IPFS相关容器
echo "🛑 停止IPFS容器..."
docker stop fabric-sdk-ipfs fabric-sdk-ipfs-cluster fabric-sdk-ipfs-webui 2>/dev/null || true

# 删除所有IPFS相关容器
echo "🗑️ 删除IPFS容器..."
docker rm -f fabric-sdk-ipfs fabric-sdk-ipfs-cluster fabric-sdk-ipfs-webui 2>/dev/null || true

# 清理IPFS网络
echo "🌐 清理IPFS网络..."
docker network rm fabric-sdk_ipfs-network 2>/dev/null || true

# 清理IPFS卷（可选，保留数据）
if [ "$1" = "--volumes" ]; then
    echo "💾 清理IPFS数据卷..."
    docker volume rm fabric-sdk_ipfs_data fabric-sdk_ipfs_staging fabric-sdk_ipfs_cluster_data 2>/dev/null || true
fi

# 清理Docker Compose
echo "🔄 清理Docker Compose..."
docker-compose -f docker-compose.ipfs.yml down -v 2>/dev/null || true
docker-compose -f docker-compose.ipfs-simple.yml down -v 2>/dev/null || true

echo "✅ IPFS资源清理完成！"
echo ""
echo "💡 现在可以重新部署IPFS服务："
echo "   ./scripts/deploy-ipfs.sh"




