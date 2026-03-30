#!/bin/bash

# Fabric SDK 香港服务器优化部署脚本
# 使用官方Go代理，适配香港网络环境

echo '========================================'
echo 'Fabric SDK 香港服务器部署脚本'
echo '========================================'
echo ''

# 配置
MAX_RETRIES=2
RETRY_DELAY=5
USE_CHINA_PROXY=false  # 香港服务器不需要中国代理

echo "配置: 重试=${MAX_RETRIES}次, 优化代理=${USE_CHINA_PROXY}"
echo ''

# 设置Go代理 - 香港使用官方代理
if [ "${USE_CHINA_PROXY}" = "false" ]; then
    echo "[优化] 使用官方Go代理 (适配香港网络)..."
    export GOPROXY=https://proxy.golang.org,direct
    export GOSUMDB=sum.golang.org
fi

# 部署重试函数
deploy_with_retry() {
    local retry_count=0
    
    while [ $retry_count -lt $MAX_RETRIES ]; do
        echo ""
        echo "========================================"
        echo "部署尝试 $((retry_count + 1))/${MAX_RETRIES}"
        echo "========================================"
        echo ""
        
        # 调用原始deploy.sh
        if ./deploy.sh; then
            echo ""
            echo "========================================"
            echo "✓ 部署成功！"
            echo "========================================"
            return 0
        else
            retry_count=$((retry_count + 1))
            
            if [ $retry_count -lt $MAX_RETRIES ]; then
                echo ""
                echo "部署失败，${RETRY_DELAY}秒后重试..."
                sleep ${RETRY_DELAY}
            else
                echo ""
                echo "========================================"
                echo "✗ 部署失败，已达到最大重试次数"
                echo "========================================"
                return 1
            fi
        fi
    done
}

# 清理旧容器和镜像
echo "[1/3] 清理旧容器和镜像..."
docker-compose -f docker-compose.yml down 2>/dev/null || true
docker rm -f fabric-sdk-app 2>/dev/null || true
docker rmi fabric-sdk_fabric-sdk:latest 2>/dev/null || true
docker rmi fabricsdk_fabric-sdk:latest 2>/dev/null || true
echo "✓ 清理完成"
echo ""

# 执行部署
deploy_with_retry

# 显示结果
if [ $? -eq 0 ]; then
    echo ""
    echo "========================================"
    echo "部署完成！"
    echo "========================================"
    echo ""
    echo "容器状态："
    docker ps --filter "name=fabric-sdk-app" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    echo "服务访问地址："
    echo "  - API: http://localhost:8199/api/index"
    echo "  - Swagger: http://localhost:8199/swagger"
    echo "  - 登录页面: http://localhost:8199/static/index.html#login"
    echo ""
    echo "查看日志："
    echo "  docker logs -f fabric-sdk-app"
else
    echo ""
    echo "部署失败，请检查错误信息"
    exit 1
fi
