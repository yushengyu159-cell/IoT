#!/bin/bash

# Fabric SDK 部署脚本
# 功能：检查并删除旧镜像，然后重新部署系统

set -e  # 遇到错误立即退出

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目目录
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Fabric SDK 部署脚本${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# 1. 检查并删除旧容器
echo -e "${YELLOW}[1/5] 检查并停止旧容器...${NC}"
if docker ps -a --format '{{.Names}}' | grep -q "^fabric-sdk-app$"; then
    echo "  发现旧容器 fabric-sdk-app，正在停止并删除..."
    docker-compose down --remove-orphans 2>/dev/null || true
    docker rm -f fabric-sdk-app 2>/dev/null || true
    echo -e "  ${GREEN}✓ 旧容器已删除${NC}"
else
    echo -e "  ${GREEN}✓ 未发现旧容器${NC}"
fi
echo ""

# 2. 检查并删除旧镜像
echo -e "${YELLOW}[2/5] 检查并删除旧镜像...${NC}"

# 获取项目名称（从docker-compose.yml所在目录名）
# docker-compose的镜像命名规则：项目目录名_服务名
PROJECT_NAME=$(basename "$PROJECT_DIR" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]//g')
IMAGE_NAME="${PROJECT_NAME}_fabric-sdk"

# 如果目录名包含特殊字符，docker-compose会替换为下划线
# 为了兼容，也检查可能的变体
ALTERNATIVE_IMAGE_NAME="fabric-sdk_fabric-sdk"

# 检查镜像是否存在（检查多种可能的命名）
FOUND_IMAGE=""
if docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "^${IMAGE_NAME}:latest$"; then
    FOUND_IMAGE="${IMAGE_NAME}:latest"
elif docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "^${ALTERNATIVE_IMAGE_NAME}:latest$"; then
    FOUND_IMAGE="${ALTERNATIVE_IMAGE_NAME}:latest"
fi

if [ -n "$FOUND_IMAGE" ]; then
    echo "  发现旧镜像 ${FOUND_IMAGE}，正在删除..."
    docker rmi -f "${FOUND_IMAGE}" 2>/dev/null || true
    echo -e "  ${GREEN}✓ 旧镜像已删除${NC}"
else
    echo -e "  ${GREEN}✓ 未发现旧镜像（${IMAGE_NAME}:latest 或 ${ALTERNATIVE_IMAGE_NAME}:latest）${NC}"
fi

# 也检查是否有其他相关镜像（包括dangling镜像）
echo "  检查相关dangling镜像..."
DANGLING_IMAGES=$(docker images -f "dangling=true" -q | wc -l)
if [ "$DANGLING_IMAGES" -gt 0 ]; then
    echo "  发现 $DANGLING_IMAGES 个dangling镜像，正在清理..."
    docker image prune -f >/dev/null 2>&1 || true
    echo -e "  ${GREEN}✓ Dangling镜像已清理${NC}"
else
    echo -e "  ${GREEN}✓ 未发现dangling镜像${NC}"
fi
echo ""

# 3. 构建新镜像
echo -e "${YELLOW}[3/5] 构建新镜像...${NC}"
docker-compose build --no-cache fabric-sdk
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✓ 镜像构建成功${NC}"
else
    echo -e "  ${RED}✗ 镜像构建失败${NC}"
    exit 1
fi
echo ""

# 4. 启动服务
echo -e "${YELLOW}[4/5] 启动服务...${NC}"
docker-compose up -d fabric-sdk
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✓ 服务启动成功${NC}"
else
    echo -e "  ${RED}✗ 服务启动失败${NC}"
    exit 1
fi
echo ""

# 5. 等待服务就绪并检查状态
echo -e "${YELLOW}[5/5] 检查服务状态...${NC}"
echo "  等待服务启动（最多等待60秒）..."
MAX_WAIT=60
WAIT_TIME=0
while [ $WAIT_TIME -lt $MAX_WAIT ]; do
    if docker ps --format '{{.Names}}' | grep -q "^fabric-sdk-app$"; then
        CONTAINER_STATUS=$(docker inspect -f '{{.State.Status}}' fabric-sdk-app 2>/dev/null || echo "not_found")
        if [ "$CONTAINER_STATUS" = "running" ]; then
            # 检查健康状态
            HEALTH=$(docker inspect -f '{{.State.Health.Status}}' fabric-sdk-app 2>/dev/null || echo "none")
            if [ "$HEALTH" = "healthy" ] || [ "$HEALTH" = "none" ]; then
                # 尝试访问API
                if curl -s -f http://localhost:8199/api/index >/dev/null 2>&1; then
                    echo -e "  ${GREEN}✓ 服务运行正常，API可访问${NC}"
                    break
                fi
            fi
        fi
    fi
    sleep 2
    WAIT_TIME=$((WAIT_TIME + 2))
    echo -n "."
done
echo ""

# 显示最终状态
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}部署完成！${NC}"
echo -e "${GREEN}========================================${NC}"
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
echo ""

