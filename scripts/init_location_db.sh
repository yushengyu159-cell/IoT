#!/bin/bash

# 位置信息数据库初始化脚本
# 创建时间: 2025-09-17

echo "开始初始化位置信息数据库..."

# 数据库连接信息
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD="Test@123456"
DB_NAME="esg"

# 执行SQL文件
mysql -h$DB_HOST -P$DB_PORT -u$DB_USER -p$DB_PASSWORD $DB_NAME < sql/location_tables.sql

if [ $? -eq 0 ]; then
    echo "✅ 位置信息数据库初始化成功"
else
    echo "❌ 位置信息数据库初始化失败"
    exit 1
fi

echo "数据库表创建完成："
echo "- provinces: 省份表"
echo "- cities: 城市表" 
echo "- districts: 区县表"
echo "- user_locations: 用户位置信息表"
echo "- pois: 兴趣点表"

echo "可以开始使用位置信息API了！"
