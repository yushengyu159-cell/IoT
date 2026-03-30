# Dashboard页面空白问题诊断

## 问题现象

用户登录DID验证成功后，访问dashboard页面显示空白。

## 根本原因

**URL格式错误**：用户访问的是 （三个斜杠），这是无效的URL格式。

## 正确的访问方式

### 从服务器外部访问

```http://47.238.159.234:8199/static/dashboard.html```

### 从服务器本机访问

```http://localhost:8199/static/dashboard.html```

### 完整登录流程

1. 访问登录页面：`http://47.238.159.234:8199/static/index.html#login`
2. 使用邮箱和DID完成登录
3. 登录成功后自动跳转到：`http://47.238.159.234:8199/static/dashboard.html`

## 诊断步骤

### 1. 检查服务是否运行

```bash
docker ps | grep fabric-sdk-app
```

正常输出应该显示容器状态为 `healthy`。

### 2. 检查端口是否监听

```bash
netstat -tuln | grep 8199
# 或
ss -tuln | grep 8199
```

应该显示端口 8199 正在监听。

### 3. 测试API端点

```bash
curl http://localhost:8199/api/index
```

应该返回 JSON 响应。

### 4. 测试静态文件访问

```bash
curl -I http://localhost:8199/static/dashboard.html
```

应该返回 `HTTP/1.1 200 OK`。

### 5. 检查浏览器控制台

打开浏览器开发者工具（F12），查看Console标签页是否有JavaScript错误。

## 常见问题

### 问题1：登录后页面空白

**原因**：localStorage中没有用户信息
**解决**：重新登录，确保DID验证成功

### 问题2：页面一直显示Loading...

**原因**：API请求失败
**解决**：检查网络连接，确认API端点可访问

### 问题3：CSS样式未加载

**原因**：styles.css文件路径错误
**解决**：确认 `/static/styles.css` 可访问

### 问题4：JavaScript错误

**原因**：外部JS文件加载失败
**解决**：检查以下文件是否可访问：
- `/static/carbon-realtime.js`
- `/static/file-management.js`
- `/static/dashboard-autoload.js`

## 已验证的正常状态

✓ 服务运行正常（healthy）
✓ API端点可访问（/api/profile/building 返回正确数据）
✓ 静态文件存在（dashboard.html, styles.css, JS文件）
✓ 端口 8199 正常监听

## 推荐访问方式

使用以下完整URL进行访问：

```
http://47.238.159.234:8199/static/index.html#login
```

登录成功后会自动跳转到dashboard页面。

---
**创建时间**：2026-03-25  
**服务器IP**：47.238.159.234  
**服务端口**：8199
