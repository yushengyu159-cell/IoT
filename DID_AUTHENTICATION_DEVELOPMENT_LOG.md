# DID身份认证功能开发过程记录

## 项目背景

本项目是一个基于区块链技术的ESG VISA数字身份管理系统，使用GoFrame框架开发后端，前端采用原生HTML/CSS/JavaScript技术栈。

### 技术架构
- **后端框架**: GoFrame v2
- **数据库**: MySQL 8.0
- **前端技术**: HTML5 + CSS3 + JavaScript ES6+
- **部署环境**: 阿里云服务器 (IP: 47.238.159.234)
- **服务端口**: 8199

## 功能需求分析

### 1. DID身份注册功能
- 用户填写个人信息（姓名、手机、邮箱、角色、年龄、密码）
- 系统生成唯一DID标识符
- 数据加密存储到MySQL数据库
- 返回注册成功信息和DID标识符

### 2. DID身份验证功能
- 用户输入DID标识符和密码
- 系统验证身份信息
- 验证成功后显示完整的用户信息
- 支持用户登出功能

### 3. 用户界面要求
- 现代化、简洁的设计风格
- 响应式布局，支持移动端
- 表单验证和错误提示
- 加载状态和成功反馈

## 开发时间线

**开始时间**: 2025年8月15日
**完成时间**: 2025年8月15日
**开发周期**: 1天

---

## 后端架构设计与实现

### 1. 数据模型设计

#### DID模型 (internal/model/did.go)
```go
type DID struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string `gorm:"size:64"`
    Phone     string `gorm:"size:32"`
    Email     string `gorm:"size:128"`
    Password  string `gorm:"size:128"`
    Role      string `gorm:"size:32"`
    Age       int
    CreatedAt string `gorm:"size:32"`
    DID       string `gorm:"column:did;size:128;uniqueIndex"`
}
```

**设计要点**:
- 使用GORM标签定义数据库约束
- 密码字段存储SHA256哈希值
- DID字段使用UUID生成唯一标识符
- 支持用户角色管理

### 2. 服务层实现

#### DID服务 (internal/service/did.go)
```go
// CreateDID 新增DID登记
func CreateDID(name, phone, email, role, password string, age int) (string, error)

// VerifyDID 通过DID和密码校验
func VerifyDID(didStr, password string) (*model.DID, error)
```

**核心功能**:
- 密码加密：使用SHA256算法
- DID生成：基于UUID的唯一标识符
- 数据验证：完整的输入参数检查
- 错误处理：详细的错误信息返回

### 3. 控制器层实现

#### DID控制器 (internal/controller/did.go)
```go
// POST /api/did/register - DID身份注册
func (c *DIDController) RegisterDID(r *ghttp.Request)

// POST /api/did/verify - DID身份验证
func (c *DIDController) VerifyDID(r *ghttp.Request)
```

**API设计**:
- RESTful风格接口设计
- 统一的响应格式
- 完整的参数验证
- 错误状态码规范

---

## 数据库配置与初始化

### 1. 环境配置

#### .env文件配置
```bash
# MySQL数据库配置
MYSQL_DSN=root:Test@123456@tcp(127.0.0.1:3306)/esg?charset=utf8mb4&parseTime=True&loc=Local

# 服务配置
SERVER_PORT=8199
```

**配置说明**:
- 使用Docker容器化MySQL 8.0
- 数据库名称：esg
- 字符集：utf8mb4
- 时区：CST (中国标准时间)

### 2. 数据库初始化

#### MySQL服务启动
```bash
# 检查MySQL容器状态
docker ps | grep mysql

# 创建数据库
docker exec -it fabric-sdk-mysql mysql -u root -pTest@123456 \
  -e "CREATE DATABASE IF NOT EXISTS esg CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

#### 数据库连接初始化 (internal/service/mysql.go)
```go
func InitMySQL() error {
    dsn := os.Getenv("MYSQL_DSN")
    if dsn == "" {
        dsn = "root:Test@123456@tcp(127.0.0.1:3306)/esg?charset=utf8mb4&parseTime=True&loc=Local"
    }
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return fmt.Errorf("MySQL连接失败: %v", err)
    }
    
    // 自动迁移表结构
    err = safeDropAndMigrate(db)
    if err != nil {
        return fmt.Errorf("MySQL自动迁移失败: %v", err)
    }
    
    DB = db
    return nil
}
```

### 3. 表结构自动迁移

**迁移策略**:
- 使用GORM AutoMigrate自动创建表结构
- 支持安全的外键约束处理
- 自动处理字段类型和大小限制

**迁移的表**:
- `dids` - DID身份信息表
- `esg_files` - ESG文件信息表

---

## 前端界面设计与实现

### 1. 设计理念

**UI设计原则**:
- 现代化、简洁的视觉风格
- 响应式布局，支持多设备访问
- 直观的用户交互体验
- 清晰的信息层次结构

**色彩方案**:
- 主色调：青绿色 (#00BCD4)
- 辅助色：绿色系 (#4CAF50, #8BC34A)
- 背景色：纯白色 (#ffffff)
- 文字色：深灰色 (#333333)

### 2. 页面结构设计

#### 主页面布局 (static/index.html)
```html
<!-- 顶部网格装饰 -->
<div class="grid-pattern"></div>

<!-- Logo区域 -->
<div class="logo-container">
    <div class="logo">
        <div class="logo-shape shape-1"></div>
        <div class="logo-shape shape-2"></div>
        <div class="logo-shape shape-3"></div>
    </div>
</div>

<!-- 标题区域 -->
<div class="header">
    <h1>DID身份管理</h1>
    <p>基于区块链技术的可信身份管理平台</p>
</div>

<!-- 表单区域 -->
<div class="form-container">
    <!-- DID注册表单 -->
    <!-- DID登录表单 -->
</div>
```

### 3. 表单设计

#### 注册表单字段
- **姓名**: 文本输入，最小2个字符
- **手机号码**: 电话输入，格式验证
- **邮箱地址**: 邮箱输入，格式验证
- **角色**: 下拉选择（普通用户/管理员/企业用户）
- **年龄**: 数字输入，范围18-120岁
- **密码**: 密码输入，最小8个字符

#### 登录表单字段
- **DID标识符**: 文本输入
- **密码**: 密码输入

### 4. 响应式设计

**断点设置**:
- 桌面端: > 768px
- 平板端: 480px - 768px
- 移动端: < 480px

**适配策略**:
- 弹性布局 (Flexbox)
- 网格布局 (CSS Grid)
- 媒体查询 (Media Queries)

---

## 前端功能实现与交互逻辑

### 1. 核心功能模块

#### 表单切换功能
```javascript
function toggleForms() {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const toggleText = document.getElementById('toggleText');
    
    if (loginForm.classList.contains('hidden')) {
        // 显示登录表单
        loginForm.classList.remove('hidden');
        registerForm.classList.add('hidden');
        toggleText.textContent = '切换到注册';
    } else {
        // 显示注册表单
        loginForm.classList.add('hidden');
        registerForm.classList.remove('hidden');
        toggleText.textContent = '切换到登录';
    }
    clearStatusMessage();
}
```

#### 表单验证逻辑
```javascript
function validateRegisterData(data) {
    if (!data.name || data.name.trim().length < 2) {
        showStatusMessage('姓名至少需要2个字符', 'error');
        return false;
    }
    
    if (!data.phone || !/^1[3-9]\d{9}$/.test(data.phone)) {
        showStatusMessage('请输入有效的手机号码', 'error');
        return false;
    }
    
    if (!data.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
        showStatusMessage('请输入有效的邮箱地址', 'error');
        return false;
    }
    
    if (!data.role) {
        showStatusMessage('请选择用户角色', 'error');
        return false;
    }
    
    if (!data.age || data.age < 18 || data.age > 120) {
        showStatusMessage('年龄必须在18-120岁之间', 'error');
        return false;
    }
    
    if (!data.password || data.password.length < 8) {
        showStatusMessage('密码至少需要8个字符', 'error');
        return false;
    }
    
    return true;
}
```

### 2. API通信模块

#### 注册请求处理
```javascript
async function handleRegister(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const registerData = {
        name: formData.get('name'),
        phone: formData.get('phone'),
        email: formData.get('email'),
        role: formData.get('role'),
        age: parseInt(formData.get('age')),
        password: formData.get('password')
    };
    
    // 验证输入
    if (!validateRegisterData(registerData)) {
        return;
    }
    
    try {
        showLoading(true);
        
        const response = await fetch(`${API_BASE_URL}/did/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(registerData)
        });
        
        const result = await response.json();
        
        if (result.code === 200) {
            showStatusMessage(`DID注册成功！您的DID是: ${result.data.did}`, 'success');
            event.target.reset();
            
            setTimeout(() => {
                showLogin();
                showStatusMessage('请使用新注册的DID进行登录', 'info');
            }, 3000);
        } else {
            showStatusMessage(result.message || '注册失败', 'error');
        }
        
    } catch (error) {
        console.error('注册错误详情:', error);
        showStatusMessage(`请求失败: ${error.message}`, 'error');
    } finally {
        showLoading(false);
    }
}
```

### 3. 状态管理

#### 用户信息存储
```javascript
// 存储用户信息到localStorage
localStorage.setItem('userInfo', JSON.stringify(result.data));
localStorage.setItem('isLoggedIn', 'true');

// 检查登录状态
function checkLoginStatus() {
    const isLoggedIn = localStorage.getItem('isLoggedIn');
    const userInfo = localStorage.getItem('userInfo');
    
    if (isLoggedIn === 'true' && userInfo) {
        try {
            const userData = JSON.parse(userInfo);
            showUserInfo(userData);
            return true;
        } catch (error) {
            console.error('解析用户信息失败:', error);
            localStorage.removeItem('userInfo');
            localStorage.removeItem('isLoggedIn');
        }
    }
    return false;
}
```

---

