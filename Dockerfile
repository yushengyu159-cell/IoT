# 多阶段构建 - 构建阶段
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fabric-sdk .

# 运行阶段
FROM alpine:latest

# 安装必要的运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/fabric-sdk .

# 复制配置文件和静态文件
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/keystore ./keystore
COPY --from=builder /app/manifest ./manifest
COPY --from=builder /app/static ./static

# 复制Fabric证书文件到应用期望的路径
COPY fabric-samples /app/fabric-samples/test-network/organizations

# 创建应用期望的硬编码路径到实际文件的符号链接
# 注意：/root 在Alpine默认不可被非root用户遍历，这里放宽目录权限用于只读访问
RUN mkdir -p /root/home/go/fabric-samples/test-network && \
    chmod 755 /root && chmod 755 /root/home && chmod 755 /root/home/go && \
    ln -sf /app/fabric-samples/test-network/organizations /root/home/go/fabric-samples/test-network/organizations

# 设置权限
RUN chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8199

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8199/api/index || exit 1

# 启动应用
CMD ["./fabric-sdk"]
