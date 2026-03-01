#!/bin/bash

# Fabric SDK 数据库连接修复脚本
# 用于清理和重置数据库状态

set -e

echo "🔧 开始修复 Fabric SDK 数据库连接..."

# 数据库配置
DB_HOST=${MYSQL_HOST:-"127.0.0.1"}
DB_PORT=${MYSQL_PORT:-"3306"}
DB_USER=${MYSQL_USER:-"root"}
DB_PASS=${MYSQL_PASS:-"Test@123456"}
DB_NAME=${MYSQL_DB:-"esg"}

echo "📊 数据库配置:"
echo "  主机: $DB_HOST:$DB_PORT"
echo "  用户: $DB_USER"
echo "  数据库: $DB_NAME"

# 设置MySQL密码环境变量
export MYSQL_PWD="$DB_PASS"

# 检查MySQL连接
echo "🔍 检查MySQL连接..."
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "SELECT 1;" >/dev/null 2>&1; then
    echo "❌ MySQL连接失败，请检查数据库配置"
    exit 1
fi
echo "✅ MySQL连接成功"

# 创建数据库（如果不存在）
echo "🗄️  确保数据库存在..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -e "CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 清理可能存在的表和外键约束
echo "🧹 清理数据库表..."
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" "$DB_NAME" <<EOF
-- 删除外键约束（如果存在）
SET FOREIGN_KEY_CHECKS = 0;

-- 删除表（如果存在）
DROP TABLE IF EXISTS esg_files;
DROP TABLE IF EXISTS dids;

-- 重新启用外键检查
SET FOREIGN_KEY_CHECKS = 1;
EOF

echo "✅ 数据库清理完成"

# 检查Go应用是否正在运行
echo "🔍 检查Go应用状态..."
if pgrep -f "go run main.go" >/dev/null; then
    echo "⚠️  检测到Go应用正在运行，建议重启应用以应用修复"
    echo "   运行命令: pkill -f 'go run main.go' && go run main.go"
else
    echo "✅ Go应用未运行，可以直接启动"
fi

echo ""
echo "🎉 数据库连接修复完成！"
echo "📝 下一步操作:"
echo "   1. 如果Go应用正在运行，请重启应用"
echo "   2. 运行: export MYSQL_DSN='root:Test@123456@tcp(127.0.0.1:3306)/esg?charset=utf8mb4&parseTime=True&loc=Local' && go run main.go"
echo "   3. 检查日志确认没有外键错误" 