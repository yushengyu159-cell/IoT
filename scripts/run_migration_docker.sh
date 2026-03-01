#!/bin/bash

# DID表结构迁移执行脚本（Docker版本）
# 为支持前端注册功能，扩展DID表结构

echo "🚀 开始执行DID表结构迁移（Docker版本）..."

# 检查Docker MySQL容器
echo "🐳 检查Docker MySQL容器..."
MYSQL_CONTAINER=$(docker ps --format "table {{.Names}}\t{{.Ports}}" | grep ":3306" | awk '{print $1}' | head -1)

if [ -z "$MYSQL_CONTAINER" ]; then
    echo "❌ 未找到运行中的MySQL Docker容器"
    echo "请确保MySQL容器正在运行："
    echo "docker ps | grep mysql"
    exit 1
fi

echo "✅ 找到MySQL容器: $MYSQL_CONTAINER"

# 检查MySQL连接
echo "📡 检查MySQL连接..."
if ! docker exec -i $MYSQL_CONTAINER mysql -u root -p"Test@123456" -e "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ MySQL连接失败，请检查密码配置"
    echo "尝试使用默认密码..."
    if ! docker exec -i $MYSQL_CONTAINER mysql -u root -e "SELECT 1;" > /dev/null 2>&1; then
        echo "❌ 默认密码也无法连接，请检查容器配置"
        exit 1
    else
        echo "✅ 使用默认密码连接成功"
        MYSQL_PASSWORD=""
    fi
else
    echo "✅ MySQL连接成功"
    MYSQL_PASSWORD="Test@123456"
fi

# 获取当前工作目录的数据库名称
echo "🔍 获取数据库名称..."
CURRENT_DIR=$(pwd)
DB_NAME=$(basename $CURRENT_DIR)

# 如果数据库名称是fabric-sdk，使用默认的test数据库
if [ "$DB_NAME" = "fabric-sdk" ]; then
    DB_NAME="test"
fi

# 如果数据库名称是scripts，使用fabric-sdk数据库
if [ "$DB_NAME" = "scripts" ]; then
    DB_NAME="fabric-sdk"
fi

echo "📊 使用数据库: $DB_NAME"

# 检查数据库是否存在，如果不存在则创建
echo "🔍 检查数据库是否存在..."
if [ -n "$MYSQL_PASSWORD" ]; then
    DB_EXISTS=$(docker exec -i $MYSQL_CONTAINER mysql -u root -p"$MYSQL_PASSWORD" -e "SHOW DATABASES LIKE '$DB_NAME';" 2>/dev/null | grep -c "$DB_NAME")
else
    DB_EXISTS=$(docker exec -i $MYSQL_CONTAINER mysql -u root -e "SHOW DATABASES LIKE '$DB_NAME';" 2>/dev/null | grep -c "$DB_NAME")
fi

if [ "$DB_EXISTS" -eq 0 ]; then
    echo "📝 数据库 $DB_NAME 不存在，正在创建..."
    if [ -n "$MYSQL_PASSWORD" ]; then
        docker exec -i $MYSQL_CONTAINER mysql -u root -p"$MYSQL_PASSWORD" -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    else
        docker exec -i $MYSQL_CONTAINER mysql -u root -e "CREATE DATABASE IF NOT EXISTS \`$DB_NAME\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
    fi
    echo "✅ 数据库 $DB_NAME 创建成功"
else
    echo "✅ 数据库 $DB_NAME 已存在"
fi

# 执行迁移脚本
echo "📝 执行数据库迁移..."
if [ -n "$MYSQL_PASSWORD" ]; then
    docker exec -i $MYSQL_CONTAINER mysql -u root -p"$MYSQL_PASSWORD" -D"$DB_NAME" < /root/home/go/fabric-sdk/scripts/migrate_did_table.sql
else
    docker exec -i $MYSQL_CONTAINER mysql -u root -D"$DB_NAME" < /root/home/go/fabric-sdk/scripts/migrate_did_table.sql
fi

if [ $? -eq 0 ]; then
    echo "✅ 数据库迁移成功完成"
    
    # 验证表结构
    echo "🔍 验证表结构..."
    if [ -n "$MYSQL_PASSWORD" ]; then
        docker exec -i $MYSQL_CONTAINER mysql -u root -p"$MYSQL_PASSWORD" -e "DESCRIBE dids;"
    else
        docker exec -i $MYSQL_CONTAINER mysql -u root -e "DESCRIBE dids;"
    fi
    
    echo "🔍 验证索引..."
    if [ -n "$MYSQL_PASSWORD" ]; then
        docker exec -i $MYSQL_CONTAINER mysql -u root -p"$MYSQL_PASSWORD" -e "SHOW INDEX FROM dids;"
    else
        docker exec -i $MYSQL_CONTAINER mysql -u root -e "SHOW INDEX FROM dids;"
    fi
    
    echo "🎉 迁移完成！DID表已成功扩展，支持前端注册的所有字段"
else
    echo "❌ 数据库迁移失败"
    exit 1
fi
