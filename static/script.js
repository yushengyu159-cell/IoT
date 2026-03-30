// 全局变量
const API_BASE_URL = '/api';

// ====== 简易多语言支持 ======
const languageData = {
    en: {
        toggleToRegister: 'Switch to Register',
        toggleToLogin: 'Switch to Sign in',
        loggingIn: 'Signing in…',
        registering: 'Registering…',
        fillLoginInfo: 'Please complete login information',
        loginSuccess: 'DID verified! Welcome back',
        requestFailed: 'Request failed: ',
        loginBtn: 'Sign in',
        registerBtn: 'Register',
        nameTooShort: 'Name must be at least 2 characters',
        invalidPhone: 'Please enter a valid mobile number',
        invalidEmail: 'Please enter a valid email address',
        selectRole: 'Please select a role',
        ageRange: 'Age must be between 18-120',
        passwordTooShort: 'Password must be at least 8 characters',
        didLabel: 'DID', nameLabel: 'Name', roleLabel: 'Role', emailLabel: 'Email', phoneLabel: 'Phone', ageLabel: 'Age', createdAtLabel: 'Created At',
        logout: 'Log out',
        welcomeBack: 'Welcome back', yourDigitalIdentity: 'Your digital identity'
    },
    zh: {
        toggleToRegister: '切换到注册',
        toggleToLogin: '切换到登录',
        loggingIn: '正在登录…',
        registering: '正在注册…',
        fillLoginInfo: '请填写完整的登录信息',
        loginSuccess: 'DID验证成功！欢迎回来',
        requestFailed: '请求失败: ',
        loginBtn: '登录',
        registerBtn: '注册',
        nameTooShort: '姓名至少需要2个字符',
        invalidPhone: '请输入有效的手机号码',
        invalidEmail: '请输入有效的邮箱地址',
        selectRole: '请选择用户角色',
        ageRange: '年龄必须在18-120岁之间',
        passwordTooShort: '密码至少需要8个字符',
        didLabel: 'DID标识符', nameLabel: '姓名', roleLabel: '角色', emailLabel: '邮箱', phoneLabel: '手机', ageLabel: '年龄', createdAtLabel: '创建时间',
        logout: '退出登录',
        welcomeBack: '欢迎回来', yourDigitalIdentity: '您的数字身份信息'
    },
    'zh-TW': {
        toggleToRegister: '切換到註冊',
        toggleToLogin: '切換到登入',
        loggingIn: '正在登入…',
        registering: '正在註冊…',
        fillLoginInfo: '請填寫完整的登入資訊',
        loginSuccess: 'DID 驗證成功！歡迎回來',
        requestFailed: '請求失敗: ',
        loginBtn: '登入',
        registerBtn: '註冊',
        nameTooShort: '姓名至少需要2個字元',
        invalidPhone: '請輸入有效的手機號碼',
        invalidEmail: '請輸入有效的郵箱地址',
        selectRole: '請選擇使用者角色',
        ageRange: '年齡必須在18-120之間',
        passwordTooShort: '密碼至少需要8個字元',
        didLabel: 'DID 識別符', nameLabel: '姓名', roleLabel: '角色', emailLabel: '郵箱', phoneLabel: '手機', ageLabel: '年齡', createdAtLabel: '建立時間',
        logout: '登出',
        welcomeBack: '歡迎回來', yourDigitalIdentity: '您的數位身分資訊'
    }
};

function getLang() {
    return localStorage.getItem('esg_language') || 'zh';
}

function tr(key, fallback) {
    const lang = getLang();
    return (languageData[lang] && languageData[lang][key]) || fallback || key;
}

function applyI18nStaticTexts() {
    const toggleText = document.getElementById('toggleText');
    if (toggleText) {
        // 文本在切换函数中动态赋值
    }
}


// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeForms();
    
    // 检查URL锚点，如果有#login则显示登录界面，否则显示注册界面
    if (window.location.hash === '#login') {
        showLogin();
    } else {
        showRegister(); // 默认显示注册界面（符合图片设计）
    }
    applyI18nStaticTexts();
});

// 初始化表单事件
function initializeForms() {
    // 登录表单提交
    document.getElementById('loginFormElement').addEventListener('submit', handleLogin);
    
    // 注册表单提交
    document.getElementById('registerFormElement').addEventListener('submit', handleRegister);
}

// 切换表单显示
function toggleForms() {
    const loginForm = document.getElementById('loginForm');
    const registerForm = document.getElementById('registerForm');
    const toggleText = document.getElementById('toggleText');
    
    if (loginForm.classList.contains('hidden')) {
        // 显示登录表单
        loginForm.classList.remove('hidden');
        registerForm.classList.add('hidden');
        toggleText.textContent = tr('toggleToRegister','切换到注册');
    } else {
        // 显示注册表单
        loginForm.classList.add('hidden');
        registerForm.classList.remove('hidden');
        toggleText.textContent = tr('toggleToLogin','切换到登录');
    }
    clearStatusMessage();
}

// 显示登录界面
function showLogin() {
    document.getElementById('loginForm').classList.remove('hidden');
    document.getElementById('registerForm').classList.add('hidden');
    document.getElementById('toggleText').textContent = tr('toggleToRegister','切换到注册');
    clearStatusMessage();
}

// 显示注册界面
function showRegister() {
    document.getElementById('loginForm').classList.add('hidden');
    document.getElementById('registerForm').classList.remove('hidden');
    document.getElementById('toggleText').textContent = tr('toggleToLogin','切换到登录');
    clearStatusMessage();
}

// 处理登录
async function handleLogin(event) {
    event.preventDefault();
    const submitBtn = event.target.querySelector('button[type="submit"]');
    if (submitBtn) { submitBtn.disabled = true; submitBtn.textContent = tr('loggingIn','正在登录…'); }
    
    const formData = new FormData(event.target);
    const loginData = {
        did: formData.get('did'),
        password: formData.get('password')
    };
    
    // 验证输入
    if (!loginData.did || !loginData.password) {
        showStatusMessage(tr('fillLoginInfo','请填写完整的登录信息'), 'error');
        return;
    }
    
    try {
        showLoading(true);
        
        console.log('发送登录请求:', loginData);
        console.log('请求URL:', `${API_BASE_URL}/did/verify`);
        
        const response = await fetch(`${API_BASE_URL}/did/verify`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(loginData)
        });
        
        console.log('响应状态:', response.status);
        console.log('响应头:', response.headers);
        
        if (!response.ok) {
            throw new Error(`HTTP错误: ${response.status} ${response.statusText}`);
        }
        
        const result = await response.json();
        console.log('响应数据:', result);
        
        if (result.code === 200) {
            showStatusMessage(tr('loginSuccess','DID验证成功！欢迎回来'), 'success');
            
            console.log('登录成功，用户数据:', result.data);
            console.log('用户数据类型:', typeof result.data);
            console.log('用户数据键:', Object.keys(result.data));
            
            // 存储用户信息到localStorage
            localStorage.setItem('userInfo', JSON.stringify(result.data));
            localStorage.setItem('isLoggedIn', 'true');
            
            // 根据用户邮箱决定跳转页面
            const userEmail = result.data.email;
            let redirectUrl = '/static/esg-dashboard.html'; // 默认跳转到仪表板
            
            // 只有 esgvisa@gmail.com 可以进入权限管理界面
            if (userEmail === 'esgvisa@gmail.com') {
                redirectUrl = '/static/control1.html';
                console.log('管理员用户，跳转到权限管理界面');
            } else {
                redirectUrl = '/static/esg-dashboard.html';
                console.log('普通用户，跳转到仪表板界面');
            }
            
            // 立即跳转到相应页面（使用绝对路径，避免相对路径问题）
            try {
                console.log('准备跳转到', redirectUrl);
                window.location.replace(redirectUrl);
                // 兜底：若浏览器阻止或未生效，100ms后再尝试一次
                setTimeout(() => {
                    if (!redirectUrl.includes(window.location.pathname)) {
                        window.location.href = redirectUrl;
                    }
                }, 100);
            } catch (e) {
                console.warn('location.replace 失败，使用 href 方式跳转', e);
                window.location.href = redirectUrl;
            }
            
        } else {
            showStatusMessage(result.message || tr('requestFailed','请求失败: '), 'error');
        }
        
    } catch (error) {
        console.error('登录错误详情:', error);
        console.error('错误类型:', error.name);
        console.error('错误消息:', error.message);
        showStatusMessage(tr('requestFailed','请求失败: ') + error.message, 'error');
    } finally {
        showLoading(false);
        if (submitBtn) { submitBtn.disabled = false; submitBtn.textContent = tr('loginBtn','登录'); }
    }
}

// 处理注册
async function handleRegister(event) {
    event.preventDefault();
    const submitBtn = event.target.querySelector('button[type="submit"]');
    if (submitBtn) { submitBtn.disabled = true; submitBtn.textContent = tr('registering','正在注册…'); }
    
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
        if (submitBtn) { submitBtn.disabled = false; submitBtn.textContent = tr('registerBtn','注册'); }
        return;
    }
    
    try {
        showLoading(true);
        
        console.log('发送注册请求:', registerData);
        console.log('请求URL:', `${API_BASE_URL}/did/register`);
        
        const response = await fetch(`${API_BASE_URL}/did/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(registerData)
        });
        
        console.log('响应状态:', response.status);
        console.log('响应头:', response.headers);
        
        if (!response.ok) {
            throw new Error(`HTTP错误: ${response.status} ${response.statusText}`);
        }
        
        const result = await response.json();
        console.log('响应数据:', result);
        
        if (result.code === 200) {
            showStatusMessage(`DID注册成功！您的DID是: ${result.data.did}`, 'success');
            
            // 清空表单
            event.target.reset();
            
            // 延迟显示登录界面
            setTimeout(() => {
                showLogin();
                showStatusMessage(tr('welcomeBack','欢迎回来') + ' · ' + tr('yourDigitalIdentity','您的数字身份信息'), 'info');
            }, 3000);
            
        } else {
            showStatusMessage(result.message || tr('requestFailed','请求失败: '), 'error');
        }
        
    } catch (error) {
        console.error('注册错误详情:', error);
        console.error('错误类型:', error.name);
        console.error('错误消息:', error.message);
        showStatusMessage(tr('requestFailed','请求失败: ') + error.message, 'error');
    } finally {
        showLoading(false);
        if (submitBtn) { submitBtn.disabled = false; submitBtn.textContent = tr('registerBtn','注册'); }
    }
}

// 验证注册数据
function validateRegisterData(data) {
    if (!data.name || data.name.trim().length < 2) {
        showStatusMessage(tr('nameTooShort','姓名至少需要2个字符'), 'error');
        return false;
    }
    
    if (!data.phone || !/^1[3-9]\d{9}$/.test(data.phone)) {
        showStatusMessage(tr('invalidPhone','请输入有效的手机号码'), 'error');
        return false;
    }
    
    if (!data.email || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email)) {
        showStatusMessage(tr('invalidEmail','请输入有效的邮箱地址'), 'error');
        return false;
    }
    
    if (!data.role) {
        showStatusMessage(tr('selectRole','请选择用户角色'), 'error');
        return false;
    }
    
    if (!data.age || data.age < 18 || data.age > 120) {
        showStatusMessage(tr('ageRange','年龄必须在18-120岁之间'), 'error');
        return false;
    }
    
    if (!data.password || data.password.length < 8) {
        showStatusMessage(tr('passwordTooShort','密码至少需要8个字符'), 'error');
        return false;
    }
    
    return true;
}

// 显示用户信息并跳转到仪表板
function showUserInfo(userData) {
    console.log('showUserInfo被调用，参数:', userData);
    
    // 保存用户信息到 localStorage
    localStorage.setItem('esg_user', JSON.stringify(userData));
    
    // 显示成功消息后跳转
    showStatusMessage(tr('loginSuccess', 'DID验证成功！欢迎回来'), 'success');
    
    setTimeout(() => {
        window.location.href = 'dashboard.html';
    }, 1000);
}

// 退出登录
function logout() {
    localStorage.removeItem('userInfo');
    localStorage.removeItem('isLoggedIn');
    window.location.replace('/static/index.html');
}

// 显示状态消息
function showStatusMessage(message, type = 'info') {
    const statusElement = document.getElementById('statusMessage');
    statusElement.textContent = message;
    statusElement.className = `status-message ${type}`;
    statusElement.classList.remove('hidden');
    
    // 自动隐藏消息
    setTimeout(() => {
        statusElement.classList.add('hidden');
    }, 5000);
}

// 清除状态消息
function clearStatusMessage() {
    const statusElement = document.getElementById('statusMessage');
    statusElement.classList.add('hidden');
}

// 显示/隐藏加载动画
function showLoading(show) {
    const loadingElement = document.getElementById('loadingOverlay');
    if (show) {
        loadingElement.classList.remove('hidden');
    } else {
        loadingElement.classList.add('hidden');
    }
}

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

// 页面加载时检查登录状态
if (checkLoginStatus()) {
    // 用户已登录，显示用户信息
} else {
    // 用户未登录，显示登录界面
    showLogin();
}

