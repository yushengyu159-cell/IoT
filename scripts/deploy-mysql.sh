#!/bin/bash

echo "🚀 开始部署MySQL服务..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker未运行，请先启动Docker服务"
    exit 1
fi

# 停止并删除现有容器
echo "🔄 清理现有MySQL容器..."
docker-compose -f docker-compose.mysql.yml down -v

# 拉取MySQL镜像
echo "📥 拉取MySQL 8.0镜像..."
docker pull mysql:8.0

# 拉取phpMyAdmin镜像
echo "📥 拉取phpMyAdmin镜像..."
docker pull phpmyadmin/phpmyadmin

# 启动MySQL服务
echo "🚀 启动MySQL服务..."
docker-compose -f docker-compose.mysql.yml up -d

# 等待MySQL启动
echo "⏳ 等待MySQL服务启动..."
sleep 30

# 检查MySQL连接
echo "🔍 检查MySQL连接..."
if docker exec fabric-sdk-mysql mysql -uroot -pTest@123456 -e "SELECT 1;" > /dev/null 2>&1; then
    echo "✅ MySQL服务启动成功！"
    echo "📊 数据库信息："
    echo "   - 主机: localhost"
    echo "   - 端口: 3306"
    echo "   - 数据库: esg"
    echo "   - 用户名: root"
    echo "   - 密码: Test@123456"
    echo "   - 用户名: fabric"
    echo "   - 密码: Test@123456"
    echo ""
    echo "🌐 phpMyAdmin管理界面: http://localhost:8080"
    echo "   - 用户名: root"
    echo "   - 密码: Test@123456"
else
    echo "❌ MySQL连接失败，请检查日志："
    docker logs fabric-sdk-mysql
    exit 1
fi

echo ""
echo "🎉 MySQL部署完成！"
echo "💡 现在可以重新启动fabric-sdk服务，启用MySQL功能"
