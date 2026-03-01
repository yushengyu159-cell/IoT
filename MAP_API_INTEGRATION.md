# 地图API集成指南

## 📍 概述

本系统已集成高德地图和百度地图API，提供完整的地理编码和逆地理编码功能。

## 🔧 配置说明

### 1. 环境变量配置

在 `.env` 文件中添加以下配置：

```bash
# 高德地图API Key (推荐)
MAP_AMAP_KEY=your_amap_api_key_here

# 百度地图API Key (备用)
MAP_BAIDU_KEY=your_baidu_api_key_here
```

### 2. 获取API Key

#### 高德地图API Key
1. 访问 [高德开放平台](https://lbs.amap.com/)
2. 注册账号并创建应用
3. 在控制台获取Web服务API Key
4. 建议申请企业版以获得更高配额

#### 百度地图API Key
1. 访问 [百度地图开放平台](https://lbsyun.baidu.com/)
2. 注册账号并创建应用
3. 在控制台获取服务端API Key
4. 建议申请企业版以获得更高配额

## 🚀 功能特性

### 1. 地理编码 (地址转坐标)
- 支持中文地址解析
- 自动地址标准化
- 返回精确经纬度坐标

### 2. 逆地理编码 (坐标转地址)
- 支持坐标转详细地址
- 返回省市区信息
- 支持格式化地址输出

### 3. 容错机制
- 高德地图为主，百度地图为备用
- API失败时自动切换
- 网络超时自动重试

## 📊 API接口

### 地理编码接口
```http
GET /api/location/geocode?address={地址}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "地理编码成功",
  "data": {
    "latitude": 39.9042,
    "longitude": 116.4074
  }
}
```

### 逆地理编码接口
```http
GET /api/location/reverse-geocode?latitude={纬度}&longitude={经度}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "逆地理编码成功",
  "data": {
    "address": "北京市朝阳区建国门外大街1号"
  }
}
```

## 🧪 测试工具

使用内置测试工具验证API功能：

```bash
# 运行地图API测试
go run test_map_api.go
```

测试内容包括：
- 多个城市地址的地理编码
- 多个坐标点的逆地理编码
- API容错机制验证

## 🔍 使用示例

### Go代码示例

```go
package main

import (
    "context"
    "fabric-sdk/internal/service"
)

func main() {
    // 创建地图API服务
    mapService := service.NewMapAPIService()
    ctx := context.Background()
    
    // 地理编码
    result, err := mapService.GeocodeAddress(ctx, "北京市朝阳区建国门外大街1号")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("坐标: %.6f, %.6f\n", result.Latitude, result.Longitude)
    
    // 逆地理编码
    address, err := mapService.ReverseGeocode(ctx, 39.9042, 116.4074)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("地址: %s\n", address)
}
```

### HTTP请求示例

```bash
# 地理编码
curl "http://localhost:8199/api/location/geocode?address=北京市朝阳区建国门外大街1号"

# 逆地理编码
curl "http://localhost:8199/api/location/reverse-geocode?latitude=39.9042&longitude=116.4074"
```

## ⚠️ 注意事项

### 1. API配额限制
- 高德地图：免费版每日1000次调用
- 百度地图：免费版每日6000次调用
- 建议申请企业版以获得更高配额

### 2. 网络要求
- 需要稳定的网络连接
- 建议配置代理服务器
- 超时时间设置为10秒

### 3. 数据精度
- 高德地图：精度较高，适合国内使用
- 百度地图：覆盖范围广，国际化支持好
- 建议根据使用场景选择主要API

## 🛠️ 故障排除

### 常见问题

1. **API Key无效**
   - 检查环境变量配置
   - 验证API Key是否正确
   - 确认API Key权限设置

2. **网络连接失败**
   - 检查网络连接
   - 验证防火墙设置
   - 尝试使用代理服务器

3. **配额超限**
   - 检查API调用次数
   - 考虑升级到企业版
   - 实现缓存机制

### 调试方法

```bash
# 查看详细日志
tail -f logs/sdk.log | grep -i "map\|geocode"

# 测试API连通性
curl "https://restapi.amap.com/v3/geocode/geo?key=YOUR_KEY&address=北京"
```

## 📈 性能优化

### 1. 缓存机制
- 实现Redis缓存
- 设置合理的缓存TTL
- 避免重复API调用

### 2. 批量处理
- 支持批量地理编码
- 减少API调用次数
- 提高处理效率

### 3. 异步处理
- 使用消息队列
- 异步处理大量请求
- 避免阻塞主流程

## 🔄 更新日志

- **v1.0.0**: 初始版本，支持高德地图和百度地图API
- **v1.1.0**: 添加容错机制和自动切换
- **v1.2.0**: 优化错误处理和日志记录
- **v1.3.0**: 添加测试工具和文档

## 📞 技术支持

如有问题，请联系开发团队或查看相关文档：
- 高德地图API文档: https://lbs.amap.com/api/webservice/summary
- 百度地图API文档: https://lbsyun.baidu.com/index.php?title=webapi
