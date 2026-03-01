# Fabric SDK with GoFrame Framework

基于GoFrame v2框架开发的Fabric SDK项目，集成Hyperledger Fabric区块链网络。

## 项目概述

本项目成功实现了GoFrame框架与Hyperledger Fabric的集成，提供了完整的区块链业务开发基础架构。通过RESTful API接口，可以方便地进行Fabric网络连接管理、状态监控和业务操作。

## 项目结构

```
fabric-sdk/
├── api/                    # API定义
├── cmd/                    # 命令行工具
├── configs/                # 配置文件
│   └── fabric/            # Fabric网络配置
├── docs/                   # 文档
├── internal/               # 内部包
│   ├── cmd/               # 命令行入口
│   ├── controller/        # 控制器层
│   ├── logic/             # 逻辑层
│   ├── model/             # 数据模型层
│   └── service/           # 服务层
├── manifest/              # 资源文件
│   ├── config/           # 配置文件
│   └── sql/              # SQL文件
├── pkg/                   # 可导出的包
│   └── fabric/           # Fabric连接管理
├── scripts/               # 脚本文件
├── test/                  # 测试文件
├── go.mod                 # Go模块文件
├── go.sum                 # 依赖校验文件
├── main.go                # 主入口文件
├── README.md              # 项目说明
├── PROJECT_SUMMARY.md     # 项目总结
└── FABRIC_INTEGRATION_SUMMARY.md  # Fabric集成总结
```

## 环境要求

- Go 1.22+
- GoFrame v2.9.0+
- Hyperledger Fabric 2.5.5
- Docker & Docker Compose

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 运行项目

```bash
go run main.go
```

### 3. 访问API

- 首页: http://localhost:8199/api
- Swagger文档: http://localhost:8199/swagger/
- OpenAPI规范: http://localhost:8199/api.json

## API接口文档

### 基础接口

#### 1. 系统信息接口
```bash
GET /api
```

**响应示例:**
```json
{
  "code": 0,
  "message": "Hello Fabric SDK with GoFrame!",
  "data": {
    "fabric_connection": {
      "message": "Fabric连接测试成功",
      "network": {
        "configPath": "configs/fabric",
        "connected": true,
        "lastTest": "2025-07-18 02:31:05",
        "network": {
          "chaincode": "basic",
          "channel": "mychannel",
          "orderer": "orderer.example.com:7050",
          "organizations": ["Org1", "Org2"],
          "peers": {
            "peer0.org1.example.com": "localhost:7051",
            "peer0.org2.example.com": "localhost:9051"
          }
        }
      },
      "status": "success"
    },
    "framework": "GoFrame v2",
    "goVersion": "1.22.2",
    "platform": "linux/amd64",
    "project": "fabric-sdk",
    "version": "1.0.0"
  }
}
```

### Fabric相关接口

#### 1. 获取连接状态
```bash
GET /api/fabric/status
```

#### 2. 获取网络信息
```bash
GET /api/fabric/info
```

#### 3. 测试连接
```bash
GET /api/fabric/test
```

#### 4. 初始化连接
```bash
GET /api/fabric/init
```

#### 5. 模拟资产操作
```bash
GET /api/fabric/simulate?operation=create&assetID=asset001
```

## 开发历程记录

### 阶段一：环境搭建
1. **Ubuntu 24.04.2 AWS EC2环境配置**
   - 安装Go 1.22.2
   - 配置Docker环境
   - 扩展EBS卷解决磁盘空间问题

2. **Hyperledger Fabric 2.5.5网络部署**
   - 下载fabric-samples
   - 启动test-network
   - 部署basic链码
   - 验证网络正常运行

### 阶段二：GoFrame项目创建
1. **项目初始化**
   ```bash
   mkdir -p go/fabric-sdk
   cd go/fabric-sdk
   go mod init fabric-sdk
   ```

2. **GoFrame框架集成**
   - 安装GoFrame v2依赖
   - 创建标准项目结构
   - 配置路由和中间件

### 阶段三：Fabric SDK集成
1. **连接管理模块开发**
   - 创建`pkg/fabric/connection.go`
   - 实现证书文件检查
   - 配置网络信息管理

2. **分层架构实现**
   - **Controller层**: 处理HTTP请求和响应
   - **Logic层**: 业务逻辑处理
   - **Service层**: 服务层抽象
   - **Model层**: 数据模型定义

3. **API接口开发**
   - Hello接口：系统信息和Fabric连接状态
   - Fabric状态接口：连接状态监控
   - 网络信息接口：Fabric网络配置
   - 测试接口：连接测试功能
   - 模拟操作接口：业务操作演示

### 阶段四：问题解决与优化
1. **证书文件路径问题**
   - 问题：期望`cert.pem`，实际为`Admin@org1.example.com-cert.pem`
   - 解决：更新连接管理器中的证书文件路径

2. **连接实例统一问题**
   - 问题：Hello接口和Fabric接口使用不同的连接实例
   - 解决：统一使用服务层的连接实例，实现自动初始化

3. **路由配置优化**
   - 问题：API路径不匹配
   - 解决：正确配置GoFrame路由绑定

## 技术实现要点

### 1. 分层架构设计
```
Controller -> Logic -> Service -> Model
    ↓           ↓        ↓        ↓
  HTTP处理   业务逻辑   服务抽象   数据模型
```

### 2. Fabric连接管理
- **证书验证**: 自动检查TLS证书文件
- **网络配置**: 支持多组织、多节点配置
- **状态监控**: 实时连接状态跟踪
- **错误处理**: 完善的错误处理机制

### 3. GoFrame特性应用
- **路由管理**: 分组路由配置
- **中间件**: 请求拦截和处理
- **配置管理**: YAML配置文件支持
- **日志系统**: 结构化日志记录
- **API文档**: 自动生成Swagger文档

### 4. 错误处理策略
- **统一响应格式**: 使用`ghttp.DefaultHandlerResponse`
- **错误码管理**: 标准化的错误码定义
- **日志记录**: 详细的错误日志记录
- **用户友好**: 清晰的错误消息提示

## 配置说明

### 项目配置
- `manifest/config/config.yaml` - 主配置文件
- `configs/fabric/` - Fabric网络配置目录

### 环境变量
- `FABRIC_CONFIG_PATH`: Fabric配置路径
- `SERVER_PORT`: 服务端口（默认8199）

## 开发规范

1. **代码组织**: 严格按照分层架构组织代码
2. **命名规范**: 使用Go语言标准命名规范
3. **错误处理**: 统一使用GoFrame错误处理机制
4. **日志记录**: 使用GoFrame日志系统记录关键操作
5. **API文档**: 通过注释自动生成API文档
6. **测试覆盖**: 编写单元测试和集成测试

## 部署说明

### 开发环境
```bash
# 启动Fabric网络
cd /home/ubuntu/go/fabric-samples/test-network
./network.sh up

# 启动GoFrame应用
cd /home/ubuntu/go/fabric-sdk
go run main.go
```

### 生产环境
```bash
# 构建可执行文件
go build -o fabric-sdk main.go

# 运行服务
./fabric-sdk
```

## 监控和维护

### 健康检查
- 访问 `/api` 接口检查系统状态
- 访问 `/api/fabric/status` 检查Fabric连接状态

### 日志查看
- 应用日志：控制台输出
- 错误日志：自动记录到日志文件

### 性能监控
- 连接状态监控
- API响应时间监控
- 资源使用情况监控

## 扩展功能

### 已实现功能
- ✅ Fabric网络连接管理
- ✅ 连接状态监控
- ✅ 网络信息查询
- ✅ 基础业务操作模拟

### 待扩展功能
- 🔄 链码调用接口
- 🔄 交易查询功能
- 🔄 用户身份管理
- 🔄 数据持久化
- 🔄 微服务架构
- 🔄 负载均衡
- 🔄 监控告警

## 技术栈

- **后端框架**: GoFrame v2.9.0
- **区块链**: Hyperledger Fabric 2.5.5
- **编程语言**: Go 1.22.2
- **容器化**: Docker & Docker Compose
- **API文档**: Swagger/OpenAPI
- **配置管理**: YAML
- **日志系统**: GoFrame Logger
- **错误处理**: GoFrame Error Handling

## 项目总结

### 成功实现的功能
1. **完整的Fabric集成**: 成功将Hyperledger Fabric与GoFrame框架集成
2. **RESTful API**: 提供完整的RESTful API接口
3. **自动文档**: 自动生成Swagger API文档
4. **连接管理**: 完善的Fabric连接管理机制
5. **错误处理**: 统一的错误处理和响应格式
6. **分层架构**: 清晰的分层架构设计

### 技术亮点
1. **模块化设计**: 高度模块化的代码组织
2. **配置驱动**: 灵活的配置管理
3. **自动初始化**: 智能的连接初始化机制
4. **状态监控**: 实时的连接状态监控
5. **开发友好**: 完善的开发工具和文档

### 项目价值
1. **开发效率**: 大幅提升区块链应用开发效率
2. **维护性**: 良好的代码组织和错误处理
3. **扩展性**: 易于扩展新功能
4. **稳定性**: 完善的错误处理和状态管理
5. **文档化**: 自动生成的API文档

## 许可证

MIT License

## 贡献指南

欢迎提交Issue和Pull Request来改进项目。

## 联系方式

如有问题，请通过以下方式联系：
- 提交GitHub Issue
- 发送邮件至项目维护者

---

**项目状态**: ✅ 已完成基础功能开发  
**最后更新**: 2025-07-18  
**版本**: v1.0.0 