# Docker部署fabric-sdk项目实现过程

## 项目概述
将fabric-sdk Go应用容器化，确保与现有MySQL、IPFS和Hyperledger Fabric 2.5.5网络正常连接。

## 1. 初始问题分析
- 系统不稳定，经常掉线
- 需要将fabric-sdk应用改为Docker容器启动
- 保持与现有MySQL、IPFS、Fabric网络的连接

## 2. Dockerfile配置

### 多阶段构建设计
```dockerfile
# 构建阶段
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fabric-sdk .

# 运行阶段
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo Asia/Shanghai > /etc/timezone
RUN addgroup -g 1001 -S appgroup && adduser -u 1001 -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/fabric-sdk .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/keystore ./keystore
COPY --from=builder /app/manifest ./manifest
COPY --from=builder /app/static ./static
COPY fabric-samples /app/fabric-samples/test-network/organizations
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8199
CMD ["./fabric-sdk"]
```

### 关键配置点
- **证书文件处理**: 直接复制fabric-samples到应用期望路径
- **权限管理**: 创建非root用户运行应用
- **时区设置**: 设置为Asia/Shanghai
- **静态文件**: 确保static目录正确复制

## 3. docker-compose.yml配置

### 服务定义
```yaml
version: '3.8'
services:
  fabric-sdk:
    build: .
    container_name: fabric-sdk-app
    restart: unless-stopped
    ports:
      - "8199:8199"
    environment:
      - GF_GERROR_BRIEF=true
      - GF_GERROR_STACK=false
      - FABRIC_SDK_GO_LOG_LEVEL=FATAL
      - FABRIC_SDK_GO_MSP_VERIFY=false
      - MYSQL_DSN=root:Test@123456@tcp(fabric-sdk-mysql:3306)/esg?charset=utf8mb4&parseTime=True&loc=Local
    volumes:
      - ./configs:/app/configs:ro
      - ./keystore:/app/keystore:ro
      - ./manifest:/app/manifest:ro
      - ./logs:/app/logs
      - ./manifest/config/config.yaml:/app/manifest/config/config.yaml:ro
    networks:
      - fabric-network
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
        reservations:
          memory: 256M
          cpus: '0.25'
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8199/api/index"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

networks:
  fabric-network:
    external: true
    name: fabric-sdk_fabric-network
```

### 关键配置点
- **外部网络**: 连接到现有的fabric-network
- **环境变量**: 使用MYSQL_DSN覆盖硬编码数据库连接
- **资源限制**: 设置内存和CPU限制
- **健康检查**: 定期检查服务可用性

## 4. 问题解决过程

### 4.1 容器名称冲突
**问题**: `Conflict. The container name "/fabric-sdk-ipfs" is already in use`
**解决**: 
```bash
docker-compose down --remove-orphans
docker rm -f fabric-sdk-ipfs fabric-sdk-mysql fabric-sdk-phpmyadmin
```

### 4.2 静态文件缺失
**问题**: `AddStaticPath failed: cannot find "static"`
**解决**: 在Dockerfile中添加
```dockerfile
COPY --from=builder /app/static ./static
```

### 4.3 端口冲突
**问题**: `failed to bind host port for 0.0.0.0:8199: address already in use`
**解决**: 
```bash
lsof -i :8199
kill -9 <PID>
```

### 4.4 数据库连接问题
**问题**: `dial tcp 127.0.0.1:3306: connect: connection refused`
**解决**: 
- 修改docker-compose.yml中的MYSQL_DSN环境变量
- 使用Docker服务名`fabric-sdk-mysql`替代`127.0.0.1`

### 4.5 证书文件权限问题
**问题**: `permission denied` 访问Fabric证书文件
**解决**: 
- 直接在Dockerfile中复制证书文件到应用期望路径
- 避免使用volumes挂载，减少权限问题

## 5. 管理员账号恢复

### 5.1 账号信息
- **邮箱**: esgvisa@gmail.com
- **密码**: 12345678
- **角色**: owner
- **权限**: 硬编码管理员权限

### 5.2 注册流程
```bash
# Step 1: 基础信息
curl -X POST http://localhost:8199/api/register/step1 \
  -H "Content-Type: application/json" \
  -d '{"email":"esgvisa@gmail.com","password":"12345678","fullName":"Admin"}'

# Step 2: 角色选择
curl -X POST http://localhost:8199/api/register/step2 \
  -H "Content-Type: application/json" \
  -d '{"email":"esgvisa@gmail.com","password":"12345678","fullName":"Admin","role":"owner","phone":"1234567890"}'

# Step 3: 建筑信息
curl -X POST http://localhost:8199/api/register/step3 \
  -H "Content-Type: application/json" \
  -d '{"email":"esgvisa@gmail.com","password":"12345678","fullName":"Admin","role":"owner","phone":"1234567890","buildingName":"Admin Building","buildingAddr":"Admin Address","buildingType":"residential"}'
```

## 6. 最终部署结果

### 6.1 服务状态
- **fabric-sdk**: Up (healthy) - 端口8199
- **MySQL**: Up - 端口3306
- **IPFS**: Up (healthy) - 端口4001, 5001, 8081
- **phpMyAdmin**: Up - 端口8080

### 6.2 网络连接
- 所有服务连接到`fabric-sdk_fabric-network`外部网络
- fabric-sdk容器可以访问MySQL和IPFS服务
- 证书文件正确配置，支持Fabric 2.5.5

### 6.3 功能验证
- API接口正常响应
- 管理员登录功能正常
- 权限管理界面可访问
- 数据库连接稳定

## 7. 部署命令

### 7.1 一键部署脚本（推荐）
使用自动化部署脚本，会自动检查并删除旧镜像和容器，然后重新部署：

```bash
cd /root/home/go/fabric-sdk
./deploy.sh
```

**脚本功能**：
- ✅ 自动检查并停止/删除旧容器 `fabric-sdk-app`
- ✅ 自动检查并删除旧镜像（项目名_fabric-sdk:latest）
- ✅ 自动清理dangling镜像
- ✅ 重新构建新镜像（使用--no-cache确保全新构建）
- ✅ 启动新容器
- ✅ 等待服务就绪并检查健康状态
- ✅ 显示部署结果和访问地址

**脚本执行流程**：
1. 检查并删除旧容器
2. 检查并删除旧镜像
3. 清理dangling镜像
4. 构建新镜像
5. 启动服务
6. 检查服务状态

### 7.2 手动构建和启动
如果需要手动控制部署过程：

```bash
cd /root/home/go/fabric-sdk
docker-compose down
docker-compose build --no-cache fabric-sdk
docker-compose up -d fabric-sdk
```

### 7.3 状态检查
```bash
docker ps | grep fabric-sdk
docker logs fabric-sdk-app --tail=10
curl -s http://localhost:8199/api/index
```

### 7.4 服务访问
- **API地址**: http://localhost:8199/api/index
- **Swagger文档**: http://localhost:8199/swagger
- **登录页面**: http://localhost:8199/static/index.html#login

## 8. 技术要点总结

### 8.1 Docker最佳实践
- 使用多阶段构建减小镜像大小
- 创建非root用户提高安全性
- 设置资源限制防止资源滥用
- 配置健康检查确保服务可用性

### 8.2 网络配置
- 使用外部网络连接现有服务
- 避免重复创建MySQL和IPFS容器
- 正确配置服务间通信

### 8.3 证书处理
- 直接复制证书文件到容器内
- 避免复杂的volumes挂载
- 确保文件权限正确

## 9. 故障排除

### 9.1 常见问题
1. **容器重启**: 检查日志 `docker logs fabric-sdk-app`
2. **数据库连接失败**: 检查MYSQL_DSN环境变量
3. **证书文件问题**: 确认fabric-samples目录存在
4. **端口冲突**: 检查端口占用情况

### 9.2 调试命令
```bash
# 查看容器状态
docker ps -a

# 查看日志
docker logs fabric-sdk-app

# 进入容器
docker exec -it fabric-sdk-app sh

# 检查网络
docker network ls
docker network inspect fabric-sdk_fabric-network
```

## 10. 部署成功标志
- ✅ 容器状态: Up (healthy)
- ✅ 端口映射: 8199:8199
- ✅ API响应: 正常返回JSON
- ✅ 数据库连接: MySQL连接成功
- ✅ 管理员登录: esgvisa@gmail.com可正常登录
- ✅ 权限管理: 管理员界面可访问

## 11. 自动化部署脚本说明

### 11.1 脚本位置
- **脚本路径**: `/root/home/go/fabric-sdk/deploy.sh`
- **执行权限**: 已设置为可执行（chmod +x）

### 11.2 脚本特性
- **智能检测**: 自动检测旧容器和镜像是否存在
- **安全删除**: 仅在确认存在时才删除，避免误操作
- **完整清理**: 同时清理dangling镜像，释放磁盘空间
- **状态监控**: 部署后自动检查服务健康状态
- **友好输出**: 使用颜色标识，清晰显示执行步骤和结果

### 11.3 使用场景
- **首次部署**: 脚本会跳过删除步骤，直接构建和启动
- **重新部署**: 自动清理旧资源，确保全新部署
- **版本更新**: 更新代码后运行脚本，自动完成部署流程

### 11.4 注意事项
- 脚本会删除旧镜像，确保有足够的磁盘空间
- 部署过程中会短暂停止服务，建议在维护窗口执行
- 如果部署失败，可以查看脚本输出的错误信息

---
**部署时间**: 2025-11-19 23:00
**部署状态**: 成功
**服务版本**: fabric-sdk v1.0.0
**Docker版本**: 支持多阶段构建
**Fabric版本**: 2.5.5
**新增功能**: 自动化部署脚本（deploy.sh）