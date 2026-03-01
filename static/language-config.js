// 统一多语言配置文件
// 所有页面共享的语言数据和转换逻辑

const LanguageConfig = {
    // 语言数据
    data: {
        en: {
            code: 'EN', 
            flag: '🇺🇸', 
            name: 'English',
            translations: {
                // 通用文本
                'welcome.title': 'Welcome to ESG VISA',
                'welcome.subtitle': 'Verified Intelligence for Sustainable Architecture',
                'create.account': 'Create your account',
                'login.account': 'Login your account',
                'sustainable.architecture': 'Sustainable architecture verification intelligence',
                'continue': 'Continue',
                'back': 'Back',
                'next': 'Next →',
                'return.login': 'Return to login',
                'terms.prefix': 'By continuing, you agree to our',
                'terms.service': 'Terms of Service',
                'privacy.policy': 'Privacy Policy',
                
                // 表单字段
                'email.placeholder': 'Enter your email',
                'password.placeholder': 'Create a password (8-16 characters)',
                'fullname.placeholder': 'Full Name',
                'phone.placeholder': 'Phone Number (Optional)',
                'remember.me': 'Remember for 30 days',
                'forgot.password': 'Forgot password',
                'loginBtn': 'Login',
                
                // 验证页面
                'verification.title': 'Verification Required',
                'verification.subtitle': 'Sustainable architecture verification intelligence',
                'verification.code.placeholder': 'Enter verification code',
                'verification.expires': 'Verification code expires in:',
                'verification.verify': 'Verify',
                'verification.resend': "Didn't get a code? Click to resend.",
                'verification.send': 'Send code',
                'verification.sending': 'Sending...',
                'verification.sent': 'Sent',
                'verification.sending.message': 'Sending verification code, please wait...',
                'verification.sent.message': 'Verification code sent successfully',
                'verification.failed': 'Failed to send verification code: ',
                'email.required': 'Email is required',
                'password.reset.success': 'Password reset successfully',
                'password.reset.failed': 'Password reset failed: ',
                
                // Index页面状态消息
                'index.email.empty': 'Please enter your email address',
                'index.email.registered': 'Email already registered, redirecting to login page...',
                'index.email.not.registered': 'Email not registered, redirecting to registration page...',
                'index.check.email.failed': 'Failed to check email',
                'index.network.error': 'Network error, please try again',
                'index.redirecting.wechat': 'Redirecting to WeChat login interface...',
                
                // 建筑位置页面
                'location.title': 'Where is your building located?',
                'location.subtitle': 'Sustainable architecture verification intelligence',
                'building.name.placeholder': 'Enter your building name',
                'address.placeholder': 'Enter your Address',
                'use.location': 'Use Location',
                
                // 位置选择器页面
                'location.page.title': 'ESG VISA - Location Selection',
                'location.search.placeholder': 'Search address or location',
                'location.map.loading': 'Map loading...',
                'location.selected.info': 'Selected Location Information',
                'location.detail.address': 'Detailed Address:',
                'location.detail.coordinates': 'Coordinates:',
                'location.confirm': 'Confirm Selection',
                'location.cancel': 'Cancel',
                'location.search.empty': 'Please enter search keywords',
                'location.search.processing': 'Searching location...',
                'location.search.success': 'Location search successful',
                'location.search.failed': 'Search failed',
                'location.search.retry': 'Search failed, please try again',
                'location.confirm.empty': 'Please select a location first',
                'location.current.failed': 'Unable to get current location, please search manually',
                'location.geolocation.unsupported': 'Browser does not support geolocation',
                'location.selected.position': 'Selected Position',
                'location.relocate': 'Relocate',
                'location.current.position': 'Current Position',
                'location.default.beijing': 'Beijing Tiananmen Square',
                'location.search.title': 'Search',
                'location.click.select': 'Click to select location',
                'location.interactive.map': 'Interactive Map',
                'location.click.anywhere': 'Click anywhere to select address',
                'location.selector.title': 'Location Selector',
                'location.selector.subtitle': 'Please use search function or click button below to select location',
                'location.get.current': 'Get Current Location',
                
                // 建筑类型页面
                'building.type.title': 'What type of building do you manage?',
                'building.type.subtitle': 'Sustainable architecture verification intelligence',
                'building.types.private': 'Private Residence',
                'building.types.village': 'Village / Ding-style House',
                'building.types.subsidized': 'Subsidized Housing',
                'building.types.retail': 'Retail Shop',
                'building.types.office': 'Office Building',
                'building.types.industrial': 'Industrial Unit',
                'building.types.serviced': 'Serviced Apartment',
                'building.types.shopping': 'Shopping Mall',
                'building.types.mixed': 'Mixed-Use Building',
                
                // 角色选择页面
                'role.title': 'What is your role?',
                'role.subtitle': 'Sustainable architecture verification intelligence',
                'role.owner': 'Owner',
                'role.manager': 'Property Manager',
                'role.institution': 'Institution(Other)',
                'property.name.placeholder': 'Enter your property name (Optional)',
                'institution.name.placeholder': 'Enter your Institution name (Optional)',
                'occupation.placeholder': 'Enter your Occupation (Optional)',
                
                // 第三方登录
                'google.login': 'Sign in with Google',
                'wechat.login': 'Sign in with Wechat',
                'apple.login': 'Sign in with Apple',
                
                // 微信扫码登录页面
                'wechat.login.title': 'WeChat QR Code Login',
                'wechat.login.subtitle': 'Scan the QR code with WeChat to login',
                'wechat.qrcode.waiting': 'Waiting for scan...',
                'wechat.qrcode.scan': 'Please scan the QR code with WeChat',
                'wechat.qrcode.scanned': 'QR code scanned, please confirm on your phone',
                'wechat.qrcode.confirmed': 'Login successful, redirecting...',
                'wechat.qrcode.expired': 'QR code expired, please refresh',
                'wechat.qrcode.refresh': 'Refresh QR Code',
                'wechat.login.back': 'Back',
                'wechat.login.getting.info': 'Getting user information...',
                'wechat.login.info.failed': 'Failed to get user information, please try again',
                'wechat.login.network.error': 'Network error, please try again',
                'wechat.login.incomplete': 'Login information incomplete, please login again',
                'wechat.qrcode.generating': 'Generating new QR code...',
                'wechat.qrcode.email.missing': 'Email information lost, returning to login page...',
                'wechat.qrcode.email.not.exists': 'Email does not exist, redirecting to registration page...',
                
                // 微信确认登录页面（手机端）
                'wechat.confirm.title': 'Confirm Login',
                'wechat.confirm.subtitle': 'Please confirm login on your phone',
                'wechat.confirm.info.title': 'You are about to log in to ESG VISA system',
                'wechat.confirm.info.desc': 'After logging in, you can manage your files and identity information',
                'wechat.confirm.button': 'Confirm Login',
                'wechat.confirm.cancel': 'Cancel',
                'wechat.confirm.success': 'Login successful! Please check on your computer',
                'wechat.confirm.success.title': '✓ Login successful',
                'wechat.confirm.success.desc': 'You can continue using the system on your computer',
                'wechat.confirm.success.close': 'You can close this page',
                'wechat.confirm.error.token': 'Missing login token',
                'wechat.confirm.error.failed': 'Login failed, please try again',
                'wechat.confirm.error.network': 'Network error, please try again',
                'wechat.confirm.processing': 'Processing login request...'
                ,
                // Control1 权限管理页面
                'control1.title': 'Permission Management - Control1',
                'control1.header': 'Permission Management System',
                'control1.subtitle': 'Control1 - User permission management and monitoring platform',
                'control1.refresh': 'Refresh Data',
                'control1.settings': 'Settings',
                'control1.logout': 'Logout',
                'control1.search.placeholder': 'Search user email, name or role...',
                'control1.filter.all': 'All Roles',
                'control1.search.btn': 'Search',
                'control1.search.clear': 'Clear',
                'control1.stats': 'System Stats',
                'control1.stats.total': 'Total Users',
                'control1.stats.active': 'Active Users',
                'control1.stats.owner': 'Owner Users',
                'control1.stats.manager': 'Managers',
                'control1.permissions': 'Permissions',
                'control1.userlist': 'User List',
                'loading.permissions': 'Loading permissions...',
                'loading.users': 'Loading users...',
                // 角色显示
                'role.owner': 'Owner',
                'role.property_manager': 'Property Manager',
                'role.institution': 'Institution User',
                // 模态框
                'modal.user.detail.title': 'User Details',
                'modal.permission.edit.title': 'Edit Permission',
                'modal.permission.userEmail': 'User Email',
                'modal.permission.userRole': 'User Role',
                'modal.permission.accountStatus': 'Account Status',
                'modal.permission.status.active': 'Active',
                'modal.permission.status.inactive': 'Inactive',
                'modal.permission.status.pending': 'Pending Review',
                'modal.permission.notes': 'Notes',
                'modal.permission.notes.placeholder': 'Enter change notes...',
                'modal.cancel': 'Cancel',
                'modal.save': 'Save Changes'
                ,
                // 通用与按钮/标签/占位
                'unknown': 'Unknown',
                'btn.refresh': 'Refresh',
                'btn.view': 'View',
                'btn.edit': 'Edit',
                'btn.delete': 'Delete',
                'btn.approve': 'Approve',
                'btn.reject': 'Reject',
                'btn.close': 'Close',
                'btn.retry': 'Retry',
                'btn.editPermission': 'Edit Permission',
                'label.email': 'Email',
                'label.fullname': 'Full Name',
                'label.role': 'Role',
                'label.phone': 'Phone',
                'label.did': 'DID',
                'label.status': 'Status',
                'label.registerTime': 'Register Time',
                'label.building': 'Building',
                'building.name.prefix': 'Name',
                'building.addr.prefix': 'Address',
                'building.type.prefix': 'Type',
                // 状态显示
                'status.pending_review': 'Pending Review',
                'status.completed': 'Completed',
                'status.rejected': 'Rejected',
                'status.active': 'Active',
                'status.inactive': 'Inactive',
                'status.unknown': 'Unknown',
                'empty.users.title': 'No Users',
                'empty.users.desc': 'There are no registered users yet',
                'loading.generic': 'Loading...'
                ,
                // 角色补充
                'role.admin': 'Administrator',
                'role.unknown': 'Unknown'
                ,
                // 消息与提示
                'msg.noPermission': 'You do not have permission to access this page',
                'msg.approve.success': 'Approved and email notification sent',
                'msg.approve.fail': 'Approval failed',
                'msg.reject.prompt': 'Please enter reason for rejection (optional):',
                'msg.reject.success': 'Rejected and email notification sent',
                'msg.reject.fail': 'Rejection failed',
                'msg.delete.confirm': 'Are you sure to delete user {email}? This cannot be undone!',
                'msg.delete.success': 'User deleted successfully',
                'msg.delete.fail': 'Failed to delete user',
                'msg.togglePermission.enabled': 'Permission "{name}" enabled',
                'msg.togglePermission.disabled': 'Permission "{name}" disabled',
                'msg.refresh.start': 'Refreshing data...',
                'msg.refresh.done': 'Data refreshed',
                'msg.settings.wip': 'Settings is under development...',
                'msg.logout.confirm': 'Confirm to logout?'
                ,
                // 权限名称与描述
                'permissions.upload.name': 'File Upload',
                'permissions.upload.desc': 'Allow user to upload ESG files',
                'permissions.download.name': 'File Download',
                'permissions.download.desc': 'Allow user to download ESG files',
                'permissions.user_mgmt.name': 'User Management',
                'permissions.user_mgmt.desc': 'Allow managing other user accounts',
                'permissions.system_settings.name': 'System Settings',
                'permissions.system_settings.desc': 'Allow modifying system configuration',
                // 等级
                'level.high': 'HIGH',
                'level.medium': 'MEDIUM',
                'level.critical': 'CRITICAL'
            }
        },
        'zh-cn': {
            code: '简体', 
            flag: '🇨🇳', 
            name: '简体中文',
            translations: {
                // 通用文本
                'welcome.title': '欢迎来到 ESG VISA',
                'welcome.subtitle': '可持续架构的验证智能',
                'create.account': '创建你的账户',
                'login.account': '登录你的账户',
                'sustainable.architecture': '可持续架构的验证智能',
                'continue': '继续',
                'back': '返回',
                'next': '下一步 →',
                'return.login': '返回登录',
                'terms.prefix': '继续即表示您同意我们的',
                'terms.service': '服务条款',
                'privacy.policy': '隐私政策',
                
                // 表单字段
                'email.placeholder': '请输入您的邮箱',
                'password.placeholder': '创建密码（8-16个字符）',
                'fullname.placeholder': '全名',
                'phone.placeholder': '手机号（可选）',
                'remember.me': '记住30天',
                'forgot.password': '忘记密码',
                'loginBtn': '登录',
                
                // 验证页面
                'verification.title': '需要验证',
                'verification.subtitle': '可持续架构的验证智能',
                'verification.code.placeholder': '输入验证码',
                'verification.expires': '验证码有效期:',
                'verification.verify': '验证',
                'verification.resend': '没收到验证码？点击重新发送。',
                'verification.send': '发送验证码',
                'verification.sending': '发送中...',
                'verification.sent': '已发送',
                'verification.sending.message': '正在发送验证码，请稍候...',
                'verification.sent.message': '验证码发送成功',
                'verification.failed': '验证码发送失败：',
                'email.required': '邮箱不能为空',
                'password.reset.success': '密码重置成功',
                'password.reset.failed': '密码重置失败：',
                
                // Index页面状态消息
                'index.email.empty': '请输入邮箱地址',
                'index.email.registered': '邮箱已注册，正在跳转到登录页面...',
                'index.email.not.registered': '邮箱未注册，正在跳转到注册页面...',
                'index.check.email.failed': '检查邮箱失败',
                'index.network.error': '网络错误，请重试',
                'index.redirecting.wechat': '正在跳转到微信登录界面...',
                
                // 建筑位置页面
                'location.title': '你的建筑位于哪里？',
                'location.subtitle': '可持续架构的验证智能',
                'building.name.placeholder': '输入你的建筑名称',
                'address.placeholder': '输入你的地址',
                'use.location': '使用位置',
                
                // 位置选择器页面
                'location.page.title': 'ESG VISA - 位置选择',
                'location.search.placeholder': '搜索地址或地点',
                'location.map.loading': '地图加载中...',
                'location.selected.info': '选择的位置信息',
                'location.detail.address': '详细地址：',
                'location.detail.coordinates': '坐标：',
                'location.confirm': '确认选择',
                'location.cancel': '取消',
                'location.search.empty': '请输入搜索关键词',
                'location.search.processing': '正在搜索位置...',
                'location.search.success': '位置搜索成功',
                'location.search.failed': '搜索失败',
                'location.search.retry': '搜索失败，请重试',
                'location.confirm.empty': '请先选择一个位置',
                'location.current.failed': '无法获取当前位置，请手动搜索',
                'location.geolocation.unsupported': '浏览器不支持地理定位',
                'location.selected.position': '选择的位置',
                'location.relocate': '重新定位',
                'location.current.position': '当前位置',
                'location.default.beijing': '北京天安门广场',
                'location.search.title': '搜索',
                'location.click.select': '点击选择位置',
                'location.interactive.map': '交互式地图',
                'location.click.anywhere': '点击任意位置选择地址',
                'location.selector.title': '位置选择器',
                'location.selector.subtitle': '请使用搜索功能或点击下方按钮选择位置',
                'location.get.current': '获取当前位置',
                
                // 建筑类型页面
                'building.type.title': '你管理的是哪类建筑？',
                'building.type.subtitle': '可持续架构的验证智能',
                'building.types.private': '私人住宅',
                'building.types.village': '村庄/丁式房屋',
                'building.types.subsidized': '补贴住房',
                'building.types.retail': '零售商店',
                'building.types.office': '办公楼',
                'building.types.industrial': '工业单元',
                'building.types.serviced': '服务式公寓',
                'building.types.shopping': '购物中心',
                'building.types.mixed': '混合用途建筑',
                
                // 角色选择页面
                'role.title': '你的角色是？',
                'role.subtitle': '可持续架构的验证智能',
                'role.owner': '业主',
                'role.manager': '物业经理',
                'role.institution': '机构（其他）',
                'property.name.placeholder': '输入你的物业名称（可选）',
                'institution.name.placeholder': '输入你的机构名称（可选）',
                'occupation.placeholder': '输入你的职业（可选）',
                
                // 第三方登录
                'google.login': '使用 Google 登录',
                'wechat.login': '使用 微信 登录',
                'apple.login': '使用 Apple 登录',
                
                // 微信扫码登录页面
                'wechat.login.title': '微信扫码登录',
                'wechat.login.subtitle': '使用微信扫描下方二维码进行登录',
                'wechat.qrcode.waiting': '等待扫描...',
                'wechat.qrcode.scan': '请使用微信扫描二维码',
                'wechat.qrcode.scanned': '二维码已扫描，请在手机上确认登录',
                'wechat.qrcode.confirmed': '登录成功，正在获取用户信息...',
                'wechat.qrcode.expired': '二维码已过期，请点击刷新按钮',
                'wechat.qrcode.refresh': '刷新二维码',
                'wechat.login.back': '返回',
                'wechat.login.getting.info': '正在获取用户信息...',
                'wechat.login.info.failed': '获取用户信息失败，请重试',
                'wechat.login.network.error': '网络错误，请重试',
                'wechat.login.incomplete': '登录信息不完整，请重新登录',
                'wechat.qrcode.generating': '正在生成新的二维码...',
                'wechat.qrcode.email.missing': '邮箱信息丢失，正在返回登录页面...',
                'wechat.qrcode.email.not.exists': '邮箱不存在，正在跳转到注册页面...',
                
                // 微信确认登录页面（手机端）
                'wechat.confirm.title': '确认登录',
                'wechat.confirm.subtitle': '请在手机上确认登录',
                'wechat.confirm.info.title': '您即将登录 ESG VISA 系统',
                'wechat.confirm.info.desc': '登录后您可以管理您的文件和身份信息',
                'wechat.confirm.button': '确认登录',
                'wechat.confirm.cancel': '取消',
                'wechat.confirm.success': '登录成功！请在电脑上查看',
                'wechat.confirm.success.title': '✓ 登录成功',
                'wechat.confirm.success.desc': '您可以在电脑上继续使用系统',
                'wechat.confirm.success.close': '可以关闭此页面',
                'wechat.confirm.error.token': '缺少登录令牌',
                'wechat.confirm.error.failed': '登录失败，请重试',
                'wechat.confirm.error.network': '网络错误，请重试',
                'wechat.confirm.processing': '正在处理登录请求...'
                ,
                // Control1 权限管理页面
                'control1.title': '权限管理系统 - Control1',
                'control1.header': '权限管理系统',
                'control1.subtitle': 'Control1 - 用户权限管理与监控平台',
                'control1.refresh': '刷新数据',
                'control1.settings': '设置',
                'control1.logout': '退出',
                'control1.search.placeholder': '搜索用户邮箱、姓名或角色...',
                'control1.filter.all': '所有角色',
                'control1.search.btn': '搜索',
                'control1.search.clear': '清除',
                'control1.stats': '系统统计',
                'control1.stats.total': '总用户数',
                'control1.stats.active': '活跃用户',
                'control1.stats.owner': '业主用户',
                'control1.stats.manager': '管理员',
                'control1.permissions': '权限管理',
                'control1.userlist': '用户列表',
                'loading.permissions': '加载权限信息...',
                'loading.users': '加载用户列表...',
                // 角色显示
                'role.owner': '业主',
                'role.property_manager': '物业管理员',
                'role.institution': '机构用户',
                // 模态框
                'modal.user.detail.title': '用户详情',
                'modal.permission.edit.title': '编辑权限',
                'modal.permission.userEmail': '用户邮箱',
                'modal.permission.userRole': '用户角色',
                'modal.permission.accountStatus': '账户状态',
                'modal.permission.status.active': '活跃',
                'modal.permission.status.inactive': '禁用',
                'modal.permission.status.pending': '待审核',
                'modal.permission.notes': '备注',
                'modal.permission.notes.placeholder': '输入权限变更备注...',
                'modal.cancel': '取消',
                'modal.save': '保存更改'
                ,
                // 通用与按钮/标签/占位
                'unknown': '未知',
                'btn.refresh': '刷新',
                'btn.view': '查看',
                'btn.edit': '编辑',
                'btn.delete': '删除',
                'btn.approve': '同意',
                'btn.reject': '不同意',
                'btn.close': '关闭',
                'btn.retry': '重试',
                'btn.editPermission': '编辑权限',
                'label.email': '用户邮箱',
                'label.fullname': '用户姓名',
                'label.role': '用户角色',
                'label.phone': '手机号码',
                'label.did': 'DID标识',
                'label.status': '账户状态',
                'label.registerTime': '注册时间',
                'label.building': '建筑信息',
                'building.name.prefix': '建筑名称',
                'building.addr.prefix': '建筑地址',
                'building.type.prefix': '建筑类型',
                // 状态显示
                'status.pending_review': '待审核',
                'status.completed': '已完成',
                'status.rejected': '已拒绝',
                'status.active': '活跃',
                'status.inactive': '禁用',
                'status.unknown': '未知',
                'empty.users.title': '暂无用户数据',
                'empty.users.desc': '系统中还没有注册用户',
                'loading.generic': '加载中...'
                ,
                // 角色补充
                'role.admin': '系统管理员',
                'role.unknown': '未知角色'
                ,
                // 消息与提示
                'msg.noPermission': '您没有权限访问此页面',
                'msg.approve.success': '审核通过，已发送邮件通知',
                'msg.approve.fail': '审批失败',
                'msg.reject.prompt': '请输入拒绝原因（可选）：',
                'msg.reject.success': '已拒绝，已发送邮件通知',
                'msg.reject.fail': '拒绝失败',
                'msg.delete.confirm': '确定要删除用户 {email} 吗？此操作不可撤销！',
                'msg.delete.success': '用户删除成功',
                'msg.delete.fail': '删除用户失败',
                'msg.togglePermission.enabled': '权限 "{name}" 已启用',
                'msg.togglePermission.disabled': '权限 "{name}" 已禁用',
                'msg.refresh.start': '正在刷新数据...',
                'msg.refresh.done': '数据刷新完成',
                'msg.settings.wip': '设置功能开发中...',
                'msg.logout.confirm': '确定要退出登录吗？'
                ,
                // 权限名称与描述
                'permissions.upload.name': '文件上传',
                'permissions.upload.desc': '允许用户上传ESG文件',
                'permissions.download.name': '文件下载',
                'permissions.download.desc': '允许用户下载ESG文件',
                'permissions.user_mgmt.name': '用户管理',
                'permissions.user_mgmt.desc': '允许管理其他用户账户',
                'permissions.system_settings.name': '系统设置',
                'permissions.system_settings.desc': '允许修改系统配置',
                // 等级
                'level.high': '高',
                'level.medium': '中',
                'level.critical': '危急'
            }
        },
        'zh-tw': {
            code: '繁體', 
            flag: '🇭🇰', 
            name: '繁體中文',
            translations: {
                // 通用文本
                'welcome.title': '歡迎來到 ESG VISA',
                'welcome.subtitle': '可持續架構的驗證智能',
                'create.account': '創建你的賬戶',
                'login.account': '登錄你的賬戶',
                'sustainable.architecture': '可持續架構的驗證智能',
                'continue': '繼續',
                'back': '返回',
                'next': '下一步 →',
                'return.login': '返回登錄',
                'terms.prefix': '繼續即表示您同意我們的',
                'terms.service': '服務條款',
                'privacy.policy': '隱私政策',
                
                // 表单字段
                'email.placeholder': '請輸入您的郵箱',
                'password.placeholder': '創建密碼（8-16個字符）',
                'fullname.placeholder': '全名',
                'phone.placeholder': '手機號（可選）',
                'remember.me': '記住30天',
                'forgot.password': '忘記密碼',
                'loginBtn': '登入',
                
                // 验证页面
                'verification.title': '需要驗證',
                'verification.subtitle': '可持續架構的驗證智能',
                'verification.code.placeholder': '輸入驗證碼',
                'verification.expires': '驗證碼有效期:',
                'verification.verify': '驗證',
                'verification.resend': '沒收到驗證碼？點擊重新發送。',
                'verification.send': '發送驗證碼',
                'verification.sending': '發送中...',
                'verification.sent': '已發送',
                'verification.sending.message': '正在發送驗證碼，請稍候...',
                'verification.sent.message': '驗證碼發送成功',
                'verification.failed': '驗證碼發送失敗：',
                'email.required': '郵箱不能為空',
                'password.reset.success': '密碼重置成功',
                'password.reset.failed': '密碼重置失敗：',
                
                // Index页面状态消息
                'index.email.empty': '請輸入郵箱地址',
                'index.email.registered': '郵箱已註冊，正在跳轉到登錄頁面...',
                'index.email.not.registered': '郵箱未註冊，正在跳轉到註冊頁面...',
                'index.check.email.failed': '檢查郵箱失敗',
                'index.network.error': '網絡錯誤，請重試',
                'index.redirecting.wechat': '正在跳轉到微信登錄界面...',
                
                // 建筑位置页面
                'location.title': '你的建築位於哪裡？',
                'location.subtitle': '可持續架構的驗證智能',
                'building.name.placeholder': '輸入你的建築名稱',
                'address.placeholder': '輸入你的地址',
                'use.location': '使用位置',
                
                // 位置选择器页面
                'location.page.title': 'ESG VISA - 位置選擇',
                'location.search.placeholder': '搜索地址或地點',
                'location.map.loading': '地圖加載中...',
                'location.selected.info': '選擇的位置信息',
                'location.detail.address': '詳細地址：',
                'location.detail.coordinates': '坐標：',
                'location.confirm': '確認選擇',
                'location.cancel': '取消',
                'location.search.empty': '請輸入搜索關鍵詞',
                'location.search.processing': '正在搜索位置...',
                'location.search.success': '位置搜索成功',
                'location.search.failed': '搜索失敗',
                'location.search.retry': '搜索失敗，請重試',
                'location.confirm.empty': '請先選擇一個位置',
                'location.current.failed': '無法獲取當前位置，請手動搜索',
                'location.geolocation.unsupported': '瀏覽器不支持地理定位',
                'location.selected.position': '選擇的位置',
                'location.relocate': '重新定位',
                'location.current.position': '當前位置',
                'location.default.beijing': '北京天安門廣場',
                'location.search.title': '搜索',
                'location.click.select': '點擊選擇位置',
                'location.interactive.map': '交互式地圖',
                'location.click.anywhere': '點擊任意位置選擇地址',
                'location.selector.title': '位置選擇器',
                'location.selector.subtitle': '請使用搜索功能或點擊下方按鈕選擇位置',
                'location.get.current': '獲取當前位置',
                
                // 建筑类型页面
                'building.type.title': '你管理的是哪類建築？',
                'building.type.subtitle': '可持續架構的驗證智能',
                'building.types.private': '私人住宅',
                'building.types.village': '村莊/丁式房屋',
                'building.types.subsidized': '補貼住房',
                'building.types.retail': '零售商店',
                'building.types.office': '辦公樓',
                'building.types.industrial': '工業單元',
                'building.types.serviced': '服務式公寓',
                'building.types.shopping': '購物中心',
                'building.types.mixed': '混合用途建築',
                
                // 角色选择页面
                'role.title': '你的角色是？',
                'role.subtitle': '可持續架構的驗證智能',
                'role.owner': '業主',
                'role.manager': '物業經理',
                'role.institution': '機構（其他）',
                'property.name.placeholder': '輸入你的物業名稱（可選）',
                'institution.name.placeholder': '輸入你的機構名稱（可選）',
                'occupation.placeholder': '輸入你的職業（可選）',
                
                // 第三方登录
                'google.login': '使用 Google 登錄',
                'wechat.login': '使用 微信 登錄',
                'apple.login': '使用 Apple 登錄',
                
                // 微信扫码登录页面
                'wechat.login.title': '微信掃碼登錄',
                'wechat.login.subtitle': '使用微信掃描下方二維碼進行登錄',
                'wechat.qrcode.waiting': '等待掃描...',
                'wechat.qrcode.scan': '請使用微信掃描二維碼',
                'wechat.qrcode.scanned': '二維碼已掃描，請在手機上確認登錄',
                'wechat.qrcode.confirmed': '登錄成功，正在獲取用戶信息...',
                'wechat.qrcode.expired': '二維碼已過期，請點擊刷新按鈕',
                'wechat.qrcode.refresh': '刷新二維碼',
                'wechat.login.back': '返回',
                'wechat.login.getting.info': '正在獲取用戶信息...',
                'wechat.login.info.failed': '獲取用戶信息失敗，請重試',
                'wechat.login.network.error': '網絡錯誤，請重試',
                'wechat.login.incomplete': '登錄信息不完整，請重新登錄',
                'wechat.qrcode.generating': '正在生成新的二維碼...',
                'wechat.qrcode.email.missing': '郵箱信息丟失，正在返回登錄頁面...',
                'wechat.qrcode.email.not.exists': '郵箱不存在，正在跳轉到註冊頁面...',
                
                // 微信確認登錄頁面（手機端）
                'wechat.confirm.title': '確認登錄',
                'wechat.confirm.subtitle': '請在手機上確認登錄',
                'wechat.confirm.info.title': '您即將登錄 ESG VISA 系統',
                'wechat.confirm.info.desc': '登錄後您可以管理您的文件和身份信息',
                'wechat.confirm.button': '確認登錄',
                'wechat.confirm.cancel': '取消',
                'wechat.confirm.success': '登錄成功！請在電腦上查看',
                'wechat.confirm.success.title': '✓ 登錄成功',
                'wechat.confirm.success.desc': '您可以在電腦上繼續使用系統',
                'wechat.confirm.success.close': '可以關閉此頁面',
                'wechat.confirm.error.token': '缺少登錄令牌',
                'wechat.confirm.error.failed': '登錄失敗，請重試',
                'wechat.confirm.error.network': '網絡錯誤，請重試',
                'wechat.confirm.processing': '正在處理登錄請求...'
                ,
                // Control1 權限管理頁面
                'control1.title': '權限管理系統 - Control1',
                'control1.header': '權限管理系統',
                'control1.subtitle': 'Control1 - 用戶權限管理與監控平台',
                'control1.refresh': '刷新數據',
                'control1.settings': '設置',
                'control1.logout': '退出',
                'control1.search.placeholder': '搜索用戶郵箱、姓名或角色...',
                'control1.filter.all': '所有角色',
                'control1.search.btn': '搜索',
                'control1.search.clear': '清除',
                'control1.stats': '系統統計',
                'control1.stats.total': '總用戶數',
                'control1.stats.active': '活躍用戶',
                'control1.stats.owner': '業主用戶',
                'control1.stats.manager': '管理員',
                'control1.permissions': '權限管理',
                'control1.userlist': '用戶列表',
                'loading.permissions': '加載權限信息...',
                'loading.users': '加載用戶列表...',
                // 角色顯示
                'role.owner': '業主',
                'role.property_manager': '物業管理員',
                'role.institution': '機構用戶',
                // 模態框
                'modal.user.detail.title': '用戶詳情',
                'modal.permission.edit.title': '編輯權限',
                'modal.permission.userEmail': '用戶郵箱',
                'modal.permission.userRole': '用戶角色',
                'modal.permission.accountStatus': '賬戶狀態',
                'modal.permission.status.active': '活躍',
                'modal.permission.status.inactive': '禁用',
                'modal.permission.status.pending': '待審核',
                'modal.permission.notes': '備註',
                'modal.permission.notes.placeholder': '輸入權限變更備註...',
                'modal.cancel': '取消',
                'modal.save': '保存更改'
                ,
                // 通用与按钮/标签/占位
                'unknown': '未知',
                'btn.refresh': '刷新',
                'btn.view': '查看',
                'btn.edit': '編輯',
                'btn.delete': '刪除',
                'btn.approve': '同意',
                'btn.reject': '不同意',
                'btn.close': '關閉',
                'btn.retry': '重試',
                'btn.editPermission': '編輯權限',
                'label.email': '用戶郵箱',
                'label.fullname': '用戶姓名',
                'label.role': '用戶角色',
                'label.phone': '手機號碼',
                'label.did': 'DID標識',
                'label.status': '賬戶狀態',
                'label.registerTime': '註冊時間',
                'label.building': '建築信息',
                'building.name.prefix': '建築名稱',
                'building.addr.prefix': '建築地址',
                'building.type.prefix': '建築類型',
                // 狀態顯示
                'status.pending_review': '待審核',
                'status.completed': '已完成',
                'status.rejected': '已拒絕',
                'status.active': '活躍',
                'status.inactive': '禁用',
                'status.unknown': '未知',
                'empty.users.title': '暫無用戶數據',
                'empty.users.desc': '系統中還沒有註冊用戶',
                'loading.generic': '加載中...'
                ,
                // 角色補充
                'role.admin': '系統管理員',
                'role.unknown': '未知角色'
                ,
                // 消息與提示
                'msg.noPermission': '您沒有權限訪問此頁面',
                'msg.approve.success': '審核通過，已發送郵件通知',
                'msg.approve.fail': '審批失敗',
                'msg.reject.prompt': '請輸入拒絕原因（可選）：',
                'msg.reject.success': '已拒絕，已發送郵件通知',
                'msg.reject.fail': '拒絕失敗',
                'msg.delete.confirm': '確定要刪除用戶 {email} 嗎？此操作不可撤銷！',
                'msg.delete.success': '用戶刪除成功',
                'msg.delete.fail': '刪除用戶失敗',
                'msg.togglePermission.enabled': '權限 "{name}" 已啟用',
                'msg.togglePermission.disabled': '權限 "{name}" 已禁用',
                'msg.refresh.start': '正在刷新數據...',
                'msg.refresh.done': '數據刷新完成',
                'msg.settings.wip': '設置功能開發中...',
                'msg.logout.confirm': '確定要退出登錄嗎？'
                ,
                // 權限名稱與描述
                'permissions.upload.name': '文件上傳',
                'permissions.upload.desc': '允許用戶上傳ESG文件',
                'permissions.download.name': '文件下載',
                'permissions.download.desc': '允許用戶下載ESG文件',
                'permissions.user_mgmt.name': '用戶管理',
                'permissions.user_mgmt.desc': '允許管理其他用戶賬戶',
                'permissions.system_settings.name': '系統設置',
                'permissions.system_settings.desc': '允許修改系統配置',
                // 等級
                'level.high': '高',
                'level.medium': '中',
                'level.critical': '嚴重'
            }
        }
    },

    // 当前语言
    currentLanguage: 'en',

    // 初始化
    init() {
        // 从localStorage获取保存的语言设置
        const saved = localStorage.getItem('esg_language') || 'en';
        this.currentLanguage = saved;
        this.applyLanguageToAllPages();
    },

    // 切换语言
    switchLanguage(lang) {
        if (!this.data[lang]) return;
        
        this.currentLanguage = lang;
        localStorage.setItem('esg_language', lang);
        
        // 更新所有页面的语言
        this.applyLanguageToAllPages();
        
        // 更新语言选择器UI
        this.updateLanguageSwitcher();
    },

    // 获取翻译文本
    getText(key) {
        const langData = this.data[this.currentLanguage];
        return langData?.translations[key] || key;
    },

    // 更新语言选择器UI
    updateLanguageSwitcher() {
        const langBtn = document.getElementById('langBtn');
        const currentLangSpan = document.getElementById('currentLang');
        const langOptions = document.querySelectorAll('#langMenu .lang-option');
        
        if (currentLangSpan) {
            currentLangSpan.textContent = this.data[this.currentLanguage].code;
        }
        
        langOptions.forEach(option => {
            const lang = option.getAttribute('data-lang');
            option.classList.toggle('active', lang === this.currentLanguage);
        });
    },

    // 应用语言到所有页面
    applyLanguageToAllPages() {
        const langData = this.data[this.currentLanguage];
        if (!langData) return;

        // 通用元素翻译
        this.translateCommonElements(langData);
        
        // 页面特定翻译
        this.translatePageSpecific(langData);
        
        // 自动翻译所有带有data-translate属性的元素
        this.autoTranslateElements(langData);
    },

    // 自动翻译所有带有data-translate属性的元素
    autoTranslateElements(langData) {
        // 翻译所有带有data-translate属性的元素
        document.querySelectorAll('[data-translate]').forEach(element => {
            const key = element.getAttribute('data-translate');
            const translation = langData.translations[key];
            
            if (translation) {
                if (element.tagName === 'INPUT' && element.type !== 'submit' && element.type !== 'button') {
                    // 输入框更新placeholder
                    element.placeholder = translation;
                } else if (element.tagName === 'BUTTON' || element.type === 'submit') {
                    // 按钮更新文本内容
                    element.textContent = translation;
                } else if (element.tagName === 'TITLE') {
                    // 页面标题更新
                    element.textContent = translation;
                } else {
                    // 其他元素更新文本内容
                    element.textContent = translation;
                }
            }
        });
    },

    // 翻译通用元素
    translateCommonElements(langData) {
        // 标题
        const titles = document.querySelectorAll('h1[data-translate]');
        titles.forEach(title => {
            const key = title.getAttribute('data-translate');
            title.textContent = langData.translations[key] || title.textContent;
        });

        // 副标题
        const subtitles = document.querySelectorAll('p[data-translate]');
        subtitles.forEach(subtitle => {
            const key = subtitle.getAttribute('data-translate');
            subtitle.textContent = langData.translations[key] || subtitle.textContent;
        });

        // 按钮
        const buttons = document.querySelectorAll('button[data-translate]');
        buttons.forEach(button => {
            const key = button.getAttribute('data-translate');
            button.textContent = langData.translations[key] || button.textContent;
        });

        // 输入框占位符
        const inputs = document.querySelectorAll('input[data-translate]');
        inputs.forEach(input => {
            const key = input.getAttribute('data-translate');
            input.placeholder = langData.translations[key] || input.placeholder;
        });

        // 链接
        const links = document.querySelectorAll('a[data-translate]');
        links.forEach(link => {
            const key = link.getAttribute('data-translate');
            link.textContent = langData.translations[key] || link.textContent;
        });
    },

    // 页面特定翻译
    translatePageSpecific(langData) {
        // 根据页面ID进行特定翻译
        const pageId = document.body.getAttribute('data-ui-id') || 
                      document.querySelector('[data-ui-id]')?.getAttribute('data-ui-id');
        
        switch(pageId) {
            case 'register-1':
            case 'register-2':
                this.translateRegistrationPage(langData);
                break;
            case 'verification':
                this.translateVerificationPage(langData);
                break;
            case 'location':
            case 'location-selector':
                this.translateLocationPage(langData);
                break;
            case 'building-type':
                this.translateBuildingTypePage(langData);
                break;
            case 'role-selection':
                this.translateRoleSelectionPage(langData);
                break;
            case 'register-12':
                this.translateForgotPasswordPage(langData);
                break;
        }
    },

    // 注册页面翻译
    translateRegistrationPage(langData) {
        // 实现注册页面特定翻译
        const elements = {
            'h1': 'create.account',
            'p.subtitle': 'sustainable.architecture',
            'input[type="email"]': 'email.placeholder',
            'input[type="password"]': 'password.placeholder',
            'input[name="fullname"]': 'fullname.placeholder',
            'input[name="phone"]': 'phone.placeholder',
            'button[type="submit"]': 'continue',
            'a.return-login': 'return.login'
        };

        Object.entries(elements).forEach(([selector, key]) => {
            const element = document.querySelector(selector);
            if (element) {
                if (selector.includes('input')) {
                    element.placeholder = langData.translations[key];
                } else {
                    element.textContent = langData.translations[key];
                }
            }
        });
    },

    // 验证页面翻译
    translateVerificationPage(langData) {
        // 实现验证页面特定翻译
    },

    // 位置页面翻译
    translateLocationPage(langData) {
        // 实现位置页面特定翻译
        const elements = {
            'h1': 'location.title',
            'p.subtitle': 'location.subtitle',
            'input[name="building_name"]': 'building.name.placeholder',
            'input[name="address"]': 'address.placeholder',
            'button.use-location': 'use.location',
            'button[type="submit"]': 'continue'
        };

        Object.entries(elements).forEach(([selector, key]) => {
            const element = document.querySelector(selector);
            if (element) {
                if (selector.includes('input')) {
                    element.placeholder = langData.translations[key];
                } else {
                    element.textContent = langData.translations[key];
                }
            }
        });
    },

    // 建筑类型页面翻译
    translateBuildingTypePage(langData) {
        // 实现建筑类型页面特定翻译
    },

    // 角色选择页面翻译
    translateRoleSelectionPage(langData) {
        // 实现角色选择页面特定翻译
    },

    // 忘记密码页面翻译
    translateForgotPasswordPage(langData) {
        // 实现忘记密码页面特定翻译
        const elements = {
            'h1': 'forgot.password',
            'p.subtitle': 'verification.subtitle',
            'input[type="email"]': 'email.placeholder',
            'input[type="text"]': 'verification.code.placeholder',
            'input[type="password"]': 'password.placeholder',
            'button[id="sendCodeBtn"]': 'verification.send',
            'button[id="verifyBtn"]': 'verification.verify',
            'a.back-link span': 'back'
        };

        Object.entries(elements).forEach(([selector, key]) => {
            const element = document.querySelector(selector);
            if (element) {
                if (selector.includes('input')) {
                    element.placeholder = langData.translations[key];
                } else {
                    element.textContent = langData.translations[key];
                }
            }
        });
    }
};

// 全局语言切换函数
window.switchLanguage = function(lang) {
    LanguageConfig.switchLanguage(lang);
};

// 页面加载时初始化
document.addEventListener('DOMContentLoaded', function() {
    LanguageConfig.init();
});

// 导出供其他脚本使用
window.LanguageConfig = LanguageConfig;
