# QQ邮箱SMTP配置指南

## 问题描述
当前QQ邮箱出现认证失败错误：
```
535 Login fail. Account is abnormal, service is not open, password is not correct, login frequency limited, or system is busy
```

## 解决方案

### 1. 获取新的SMTP授权码

#### 步骤1：登录QQ邮箱
- 访问：https://mail.qq.com
- 使用QQ账号登录

#### 步骤2：开启SMTP服务
1. 点击 **设置** → **账户**
2. 找到 **POP3/IMAP/SMTP/Exchange/CardDAV/CalDAV服务**
3. 开启 **SMTP服务**

#### 步骤3：获取授权码
1. 开启SMTP服务后，系统会要求验证身份
2. 验证成功后，会显示 **授权码**
3. 复制这个授权码（不是QQ密码）

### 2. 更新配置文件

#### 修改 `configs/email.yaml`
```yaml
# QQ邮箱配置
qq:
  host: "smtp.qq.com"
  port: "465"  # 使用SSL端口
  username: "tj18832045990@qq.com"
  password: "YOUR_NEW_SMTP_AUTH_CODE"  # 替换为新获取的授权码
  from: "tj18832045990@qq.com"
  enabled: true  # 重新启用
  use_ssl: true
  use_tls: false
```

### 3. 常见问题排查

#### 问题1：授权码错误
- 确保使用的是SMTP授权码，不是QQ密码
- 授权码通常是一串16位字符

#### 问题2：端口配置
- 推荐使用465端口（SSL）
- 587端口（TLS）也可以，但需要正确配置

#### 问题3：安全设置
- 确保QQ邮箱开启了SMTP服务
- 检查是否有登录频率限制

#### 问题4：网络问题
- 检查服务器网络连接
- 确认防火墙设置

### 4. 测试配置

#### 重启服务
```bash
cd /root/home/go/fabric-sdk
go run main.go
```

#### 测试邮件发送
1. 访问注册页面
2. 输入邮箱地址
3. 点击发送验证码
4. 检查邮箱是否收到

### 5. 临时解决方案

如果QQ邮箱暂时无法使用，可以：
1. 使用163邮箱进行注册（已配置正常）
2. 或者使用其他支持SMTP的邮箱服务

### 6. 配置验证

成功配置后，日志应该显示：
```
✅ 获取邮箱配置成功: smtp.qq.com 端口: 465 SSL: true TLS: false
📤 开始发送邮件...
✅ 邮件发送成功: tj18832045990@qq.com
```

---

**注意事项**：
- 授权码是敏感信息，不要泄露
- 定期更新授权码，提高安全性
- 生产环境建议使用企业邮箱服务

**技术支持**：
- QQ邮箱帮助：https://help.mail.qq.com/detail/108/1023

