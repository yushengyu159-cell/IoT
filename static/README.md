# ESG VISA 前端界面

## 概述

这是一个基于现代Web技术构建的ESG VISA数字身份认证系统前端界面，包含登录和注册功能，并与后端DID API进行交互。

## 文件结构

```
static/
├── index.html          # 主页面HTML
├── styles.css          # 样式文件
├── script.js           # JavaScript逻辑
└── README.md           # 说明文档
```

## 功能特性

### 1. 用户注册
- 完整的用户信息收集
- 实时表单验证
- 密码强度检查
- 角色选择（普通用户/管理员/企业用户）

### 2. 用户登录
- DID标识符登录
- 密码验证
- 登录状态保持
- 自动跳转用户信息页面

### 3. 用户信息展示
- 完整的DID信息显示
- 响应式设计
- 退出登录功能

### 4. 界面特性
- 现代化UI设计
- 渐变背景和动画效果
- 响应式布局
- 加载状态提示
- 错误和成功消息显示

## 使用方法

### 1. 启动服务
确保后端fabric-sdk服务正在运行：
```bash
cd /root/home/go/fabric-sdk
go run main.go
```

### 2. 访问前端
在浏览器中打开：
```
http://localhost:8199/static/
```

### 3. 测试流程

#### 注册新用户
1. 点击"立即注册"
2. 填写完整的用户信息
3. 提交注册表单
4. 系统会返回生成的DID标识符
5. 自动跳转到登录界面

#### 用户登录
1. 输入注册时获得的DID标识符
2. 输入密码
3. 点击登录按钮
4. 登录成功后显示用户信息

## API接口

### 注册接口
```
POST /api/did/register
Content-Type: application/json

{
    "name": "用户名",
    "phone": "手机号",
    "email": "邮箱",
    "role": "角色",
    "age": 年龄,
    "password": "密码"
}
```

### 验证接口
```
POST /api/did/verify
Content-Type: application/json

{
    "did": "DID标识符",
    "password": "密码"
}
```

## 技术栈

- **HTML5**: 语义化标签和表单
- **CSS3**: 现代样式、动画和响应式设计
- **JavaScript ES6+**: 异步操作、表单处理、API调用
- **Font Awesome**: 图标库
- **LocalStorage**: 本地状态管理

## 浏览器兼容性

- Chrome 60+
- Firefox 55+
- Safari 12+
- Edge 79+

## 开发说明

### 样式定制
- 主色调：`#667eea` 到 `#764ba2` 的渐变
- 圆角：20px
- 阴影：现代化的box-shadow效果
- 动画：CSS3 transitions和keyframes

### JavaScript功能
- 表单验证和提交
- API调用和错误处理
- 状态管理和本地存储
- 动态内容渲染

## 部署说明

1. 将static文件夹放在fabric-sdk项目的根目录
2. 确保后端服务正常运行
3. 访问对应的URL即可使用

## 注意事项

1. 确保后端API地址正确（默认：http://localhost:8199/api）
2. 密码至少8位字符
3. 手机号格式：1开头的11位数字
4. 邮箱格式验证
5. 年龄范围：18-120岁

## 故障排除

### 常见问题
1. **API连接失败**: 检查后端服务是否启动
2. **样式不显示**: 检查CSS文件路径
3. **功能不工作**: 检查浏览器控制台错误信息

### 调试方法
1. 打开浏览器开发者工具
2. 查看Console标签页的错误信息
3. 检查Network标签页的API请求状态
4. 验证HTML元素的ID和类名是否正确




