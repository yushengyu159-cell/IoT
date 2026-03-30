# 香港服务器部署方案记录

## 问题描述

在香港服务器上部署应用时，使用中国大陆的Go代理（goproxy.cn）导致部署问题。

## 根本原因

- 服务器位置：香港（HK）
- 网络环境：香港可以直连国际网络，无需使用中国代理
- 问题配置：之前使用了 GOPROXY=https://goproxy.cn,https://goproxy.io,direct

## 正确方案

### Go代理配置

对于香港服务器，应该使用官方Go代理：

```bash
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org
```

**配置说明：**
- https://proxy.golang.org - Go官方代理，全球可用
- direct - 代理不可用时直接从源码仓库下载
- sum.golang.org - Go官方校验和数据库

### 错误配置对比

| 服务器位置 | 错误配置 | 正确配置 |
|-----------|---------|---------|
| 中国大陆 | 无（应使用中国代理） | GOPROXY=https://goproxy.cn,direct |
| **香港** | GOPROXY=https://goproxy.cn,direct | GOPROXY=https://proxy.golang.org,direct ✅ |
| 美国/其他 | GOPROXY=https://goproxy.cn,direct | GOPROXY=https://proxy.golang.org,direct ✅ |

## 部署脚本

已创建专用部署脚本：deploy-hk.sh

### 脚本特点

1. 使用官方Go代理
2. 保留重试机制（最多2次）
3. 自动清理旧容器和镜像
4. 详细的部署日志

### 使用方法

```bash
cd /root/home/go/fabric-sdk
./deploy-hk.sh
```

## 部署结果

### 成功部署指标

- ✓ 镜像构建成功
- ✓ Go模块下载：117.8秒
- ✓ 代码编译：64.0秒
- ✓ 容器健康状态：healthy
- ✓ 服务端口：8199

### 服务访问地址

- API: http://localhost:8199/api/index
- Swagger: http://localhost:8199/swagger
- 登录页面: http://localhost:8199/static/index.html#login

## 脚本对比

| 脚本 | Go代理 | 适用场景 | 状态 |
|-----|--------|---------|------|
| deploy.sh | 默认（系统代理） | 通用 | ✅ 可用 |
| deploy-with-retry.sh | goproxy.cn | 中国大陆 | ❌ 香港不适用 |
| **deploy-hk.sh** | proxy.golang.org | **香港/国际** | ✅ **推荐** |

## 重要提示

1. **不要在香港服务器使用中国Go代理**
   - goproxy.cn 和 goproxy.io 专为国内网络优化
   - 香港网络可以直接访问国际服务
   - 使用中国代理可能增加延迟或导致连接问题

2. **网络诊断方法**

```bash
# 测试官方代理连接
curl -I https://proxy.golang.org

# 测试中国代理连接
curl -I https://goproxy.cn

# 查看当前Go代理配置
go env GOPROXY
```

3. **临时修改Go代理**

```bash
# 仅在当前shell设置
export GOPROXY=https://proxy.golang.org,direct

# 永久设置（写入 ~/.bashrc 或 ~/.profile）
echo 'export GOPROXY=https://proxy.golang.org,direct' >> ~/.bashrc
source ~/.bashrc
```

## 版本信息

- 创建时间：2026-03-25
- 服务器位置：香港（HK）
- 服务器IP：47.238.159.234
- 应用：fabric-sdk
- 部署方式：Docker Compose

## 相关文档

- FRONTEND_FUNCTIONALITY_TEST_REPORT.md - 前端功能测试报告
- OPERATION_RECORDS.md - Fabric网络操作记录
- DASHBOARD_PERFORMANCE_OPTIMIZATION.md - 性能优化记录

---
**记录人：** Claude Code Assistant  
**最后更新：** 2026-03-25
