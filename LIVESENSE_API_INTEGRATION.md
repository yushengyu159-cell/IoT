# LiveSense API 集成文档

## 📋 概述

本系统已集成 LiveSense GraphQL API，为 belucksapi 用户提供传感器数据获取功能。

## 🔐 权限说明

**重要**: 所有 LiveSense API 接口仅限 **belucksapi** 用户使用。

- 用户邮箱: `18027473816@163.com`
- 用户名: `belucksapi`

系统会自动验证用户身份，非授权用户将无法访问。

## 🚀 API 接口

### 1. 获取传感器数值

**接口**: `GET /api/livesense/values`

**参数**:
- `email` (必需): 用户邮箱，用于身份验证
- `id` (必需): 传感器ID
- `from` (可选): 起始时间戳（毫秒）
- `to` (可选): 结束时间戳（毫秒）
- `aggregation` (可选): 聚合方式
- `group` (可选): 分组方式
- `limit` (可选): 限制返回数量

**请求示例**:
```bash
curl -X GET "http://localhost:8199/api/livesense/values?email=18027473816@163.com&id=27732&from=1769327847000&to=1769414247000&group=allValues"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "values": [
      {
        "sensor": "sensor_id",
        "timezone": "UTC",
        "dataPoints": [
          {
            "value": 25.5,
            "timestampUtc": "2025-01-28T10:00:00Z"
          }
        ]
      }
    ]
  }
}
```

### 2. 获取上下文下的传感器列表

**接口**: `GET /api/livesense/sensors`

**参数**:
- `email` (必需): 用户邮箱，用于身份验证
- `contextId` (必需): 上下文ID（整数）
- `contextTypeId` (可选): 上下文类型ID
- `groupByTypeId` (可选): 分组类型ID

**请求示例**:
```bash
curl -X GET "http://localhost:8199/api/livesense/sensors?email=18027473816@163.com&contextId=7286&contextTypeId=502&groupByTypeId=599"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "获取成功",
  "data": {
    "sensorsUnderContext": [
      {
        "id": 7286,
        "name": "Building A",
        "sensors": [
          {
            "id": 27732,
            "name": "Temperature Sensor"
          }
        ]
      }
    ]
  }
}
```

### 3. 认证接口（测试用）

**接口**: `POST /api/livesense/authenticate`

**参数**:
- `email` (必需): 用户邮箱，用于身份验证

**请求示例**:
```bash
curl -X POST "http://localhost:8199/api/livesense/authenticate?email=18027473816@163.com"
```

**响应示例**:
```json
{
  "code": 200,
  "message": "认证成功",
  "data": {
    "authenticated": true
  }
}
```

## 🔧 技术实现

### 服务层 (`internal/service/livesense_api.go`)

- `LiveSenseAPIService`: LiveSense API 服务封装
- 自动处理认证和 token 管理
- Token 自动刷新（1小时有效期）
- GraphQL 查询封装

### 控制器层 (`internal/controller/livesense.go`)

- `LiveSenseController`: API 控制器
- 用户权限验证中间件
- 请求参数验证
- 错误处理

### 路由配置 (`internal/cmd/cmd.go`)

所有接口注册在 `/api/livesense` 路由组下，并自动应用用户验证中间件。

## 📝 使用说明

### 1. 用户身份验证

所有请求必须包含 `email` 参数，系统会验证：
- 邮箱是否为 `18027473816@163.com`，或
- 数据库中是否存在该邮箱且用户名为 `belucksapi`

### 2. 认证流程

系统会自动处理 LiveSense API 的认证：
1. 首次请求时自动获取 token
2. Token 缓存 1 小时
3. Token 过期后自动刷新

### 3. 错误处理

常见错误码：
- `400`: 请求参数错误
- `401`: 未提供用户邮箱
- `403`: 无权限访问（非 belucksapi 用户）
- `500`: 服务器错误或 API 调用失败

## 🧪 测试示例

### 使用 curl 测试

```bash
# 1. 测试认证
curl -X POST "http://localhost:8199/api/livesense/authenticate?email=18027473816@163.com"

# 2. 获取传感器数值
curl -X GET "http://localhost:8199/api/livesense/values?email=18027473816@163.com&id=27732&from=1769327847000&to=1769414247000"

# 3. 获取传感器列表
curl -X GET "http://localhost:8199/api/livesense/sensors?email=18027473816@163.com&contextId=7286"
```

### 使用 Postman

1. 导入 Postman 集合（参考 `调试文件/ec74a9ad-8152-4a14-9746-fe942efae337`）
2. 修改请求 URL 为系统地址
3. 添加 `email` 查询参数

## 🔒 安全注意事项

1. **用户验证**: 所有接口都经过用户身份验证
2. **Token 管理**: Token 自动管理，无需手动处理
3. **HTTPS**: 生产环境建议使用 HTTPS
4. **日志记录**: 所有 API 调用都会记录日志

## 📊 数据流程

```
用户请求 → 系统验证用户身份 → 调用 LiveSense API → 返回数据
```

1. 用户通过系统 API 发送请求
2. 系统验证用户是否为 belucksapi
3. 系统自动处理 LiveSense API 认证
4. 系统调用 LiveSense GraphQL API
5. 系统返回格式化后的数据

## 🐛 故障排除

### 认证失败
- 检查用户名和密码是否正确
- 检查网络连接
- 查看系统日志

### 权限错误
- 确认用户邮箱为 `18027473816@163.com`
- 确认数据库中用户名为 `belucksapi`

### API 调用失败
- 检查 LiveSense API 服务状态
- 查看系统日志获取详细错误信息
- 确认请求参数格式正确

## 📅 更新记录

- **2026-01-29**: 初始版本，集成 LiveSense GraphQL API

