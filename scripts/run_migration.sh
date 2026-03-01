#!/bin/bash

# DID表结构迁移执行脚本
# 为支持前端注册功能，扩展DID表结构

echo "🚀 开始执行DID表结构迁移..."

# 检查MySQL连接
echo "📡 检查MySQL连接..."
if ! mysql -u root -p"Test@123456" -e "SELECT 1;" > /dev/null 2>&1; then
    echo "❌ MySQL连接失败，请检查配置"
    exit 1
fi

echo "✅ MySQL连接成功"

# 执行迁移脚本
echo "📝 执行数据库迁移..."
mysql -u root -p"Test@123456" < scripts/migrate_did_table.sql

if [ $? -eq 0 ]; then
    echo "✅ 数据库迁移成功完成"
    
    # 验证表结构
    echo "🔍 验证表结构..."
    mysql -u root -p"Test@123456" -e "DESCRIBE dids;"
    
    echo "🔍 验证索引..."
    mysql -u root -p"Test@123456" -e "SHOW INDEX FROM dids;"
    
    echo "🎉 迁移完成！DID表已成功扩展，支持前端注册的所有字段"
else
    echo "❌ 数据库迁移失败"
    exit 1
fi
