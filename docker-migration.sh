#!/bin/bash

# Docker迁移脚本
echo "🚀 开始Docker迁移..."

# 1. 停止当前服务
echo "📦 停止当前服务..."
pkill -f fabric-sdk || true
docker-compose down || true

# 2. 备份数据
echo "💾 备份数据..."
mkdir -p ./backup
cp -r ./configs ./backup/ 2>/dev/null || true
cp -r ./keystore ./backup/ 2>/dev/null || true
cp -r ./manifest ./backup/ 2>/dev/null || true

# 3. 创建日志目录
echo "📝 创建日志目录..."
mkdir -p ./logs

# 4. 构建Docker镜像
echo "🔨 构建Docker镜像..."
docker-compose build

# 5. 启动服务
echo "🚀 启动Docker服务..."
docker-compose up -d

# 6. 等待服务启动
echo "⏳ 等待服务启动..."
sleep 30

# 7. 检查服务状态
echo "🔍 检查服务状态..."
docker-compose ps

# 8. 测试API
echo "🧪 测试API..."
curl -s http://localhost:8199/api/index || echo "API测试失败"

echo "✅ Docker迁移完成！"
echo "📊 服务状态："
echo "  - Fabric SDK: http://localhost:8199"
echo "  - phpMyAdmin: http://localhost:8080"
echo "  - IPFS Gateway: http://localhost:8081"
echo "  - IPFS API: http://localhost:5001"
