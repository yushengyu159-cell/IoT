#!/bin/bash

# Fabric SDK 启动脚本

echo "启动 Fabric SDK with GoFrame Framework..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到Go环境，请先安装Go"
    exit 1
fi

# 检查Go版本
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Go版本: $GO_VERSION"

# 安装依赖
echo "安装依赖..."
go mod tidy

# 构建项目
echo "构建项目..."
go build -o fabric-sdk main.go

# 启动服务
echo "启动HTTP服务..."
echo "服务地址: http://localhost:8199"
echo "API文档: http://localhost:8199/swagger"
echo "按 Ctrl+C 停止服务"

./fabric-sdk 