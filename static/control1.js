// 权限管理系统 - Control1 JavaScript
// 仿照文件管理界面的交互逻辑

class PermissionManager {
    constructor() {
        this.currentUser = null;
        this.users = [];
        this.filteredUsers = [];
        this.permissions = [];
        this.init();
    }

    // 初始化
    init() {
        this.loadUserInfo();
        this.loadUsers();
        this.loadPermissions();
        this.bindEvents();
        this.updateStats();
    }

    // 加载用户信息
    loadUserInfo() {
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}');
        if (userInfo.Email) {
            // 检查是否为授权管理员
            if (userInfo.Email !== 'esgvisa@gmail.com') {
                this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.noPermission'):'您没有权限访问此页面'), 'error');
                setTimeout(() => {
                    this.redirectToFileManagement();
                }, 2000);
                return;
            }
            
            this.currentUser = userInfo;
            this.updateUserDisplay();
        } else {
            // 如果没有用户信息，跳转到登录页面
            this.redirectToLogin();
        }
    }

    // 更新用户显示
    updateUserDisplay() {
        if (this.currentUser) {
            document.getElementById('userName').textContent = this.currentUser.Name || this.currentUser.Email;
            document.getElementById('userRole').textContent = this.getRoleDisplayName(this.currentUser.Role);
            document.getElementById('userAvatar').textContent = (this.currentUser.Name || this.currentUser.Email).charAt(0).toUpperCase();
        }
    }

    // 获取角色显示名称
    getRoleDisplayName(role) {
        const roleMap = {
            'owner': window.LanguageConfig ? window.LanguageConfig.getText('role.owner') : '业主',
            'property_manager': window.LanguageConfig ? window.LanguageConfig.getText('role.property_manager') : '物业管理员',
            'institution': window.LanguageConfig ? window.LanguageConfig.getText('role.institution') : '机构用户',
            'admin': window.LanguageConfig ? window.LanguageConfig.getText('role.admin') : '系统管理员'
        };
        return roleMap[role] || (window.LanguageConfig ? window.LanguageConfig.getText('role.unknown') : '未知角色');
    }

    // 加载用户列表
    async loadUsers() {
        try {
            this.showLoading('userListContainer');
            
            // 调用获取所有用户的API（管理员用，包括待审核用户）
            const response = await fetch('/api/register/admin/users');
            const result = await response.json();
            
            if (result.code === 200) {
                // 获取每个用户的详细信息
                const userPromises = result.data.emails.map(email => this.getUserDetail(email));
                const userDetails = await Promise.all(userPromises);
                
                this.users = userDetails.filter(user => user !== null);
                this.filteredUsers = [...this.users];
                this.renderUserList();
                this.updateStats();
            } else {
                this.showError('userListContainer', (window.LanguageConfig?window.LanguageConfig.getText('loading.users'):'加载用户列表失败: ') + result.message);
            }
        } catch (error) {
            console.error('加载用户列表失败:', error);
            this.showError('userListContainer', (window.LanguageConfig?window.LanguageConfig.getText('loading.users'):'加载用户列表失败: ') + error.message);
        }
    }

    // 获取用户详细信息
    async getUserDetail(email) {
        try {
            const response = await fetch(`/api/register/user-detail?email=${encodeURIComponent(email)}&t=${Date.now()}`);
            const result = await response.json();
            
            if (result.code === 200) {
                return result.data;
            } else {
                console.warn(`获取用户 ${email} 详情失败:`, result.message);
                return null;
            }
        } catch (error) {
            console.error(`获取用户 ${email} 详情失败:`, error);
            return null;
        }
    }

    // 加载权限信息
    loadPermissions() {
        // 模拟权限数据
        this.permissions = [
            {
                id: 1,
                key: 'upload',
                name: '文件上传',
                description: '允许用户上传ESG文件',
                level: 'high',
                enabled: true
            },
            {
                id: 2,
                key: 'download',
                name: '文件下载',
                description: '允许用户下载ESG文件',
                level: 'medium',
                enabled: true
            },
            {
                id: 3,
                key: 'user_mgmt',
                name: '用户管理',
                description: '允许管理其他用户账户',
                level: 'high',
                enabled: false
            },
            {
                id: 4,
                key: 'system_settings',
                name: '系统设置',
                description: '允许修改系统配置',
                level: 'critical',
                enabled: false
            }
        ];
        this.renderPermissions();
    }

    // 渲染用户列表
    renderUserList() {
        const container = document.getElementById('userListContainer');
        
        if (this.filteredUsers.length === 0) {
            container.innerHTML = `
                <div class="empty-state">
                    <i class="fas fa-users"></i>
                    <h3 data-translate="empty.users.title">暂无用户数据</h3>
                    <p data-translate="empty.users.desc">系统中还没有注册用户</p>
                    <button class="btn btn-primary" onclick="permissionManager.loadUsers()">
                        <i class="fas fa-refresh"></i> <span data-translate="btn.refresh">刷新</span>
                    </button>
                </div>
            `;
            this.applyI18n();
            return;
        }

        const userListHTML = this.filteredUsers.map(user => `
            <div class="user-item">
                <div class="user-info-item">
                    <div class="user-avatar-sm">${(user.full_name || user.email).charAt(0).toUpperCase()}</div>
                    <div class="user-details-item">
                        <h4>${user.full_name || (window.LanguageConfig?window.LanguageConfig.getText('unknown'):'未设置姓名')}</h4>
                        <p>${user.email} • ${this.getRoleDisplayName(user.role)} • ${(window.LanguageConfig?window.LanguageConfig.getText('label.status'):'状态')}: ${(window.LanguageConfig?window.LanguageConfig.getText('status.' + (user.status || 'unknown')):(user.status || '未知'))}</p>
                        <p>${(window.LanguageConfig?window.LanguageConfig.getText('label.registerTime'):'注册时间')}: ${user.created_at || (window.LanguageConfig?window.LanguageConfig.getText('unknown'):'未知')}</p>
                    </div>
                </div>
                <div class="permission-buttons">
                    <span class="role-badge role-${user.role}">${this.getRoleDisplayName(user.role)}</span>
                    <button class="btn btn-sm btn-primary" onclick="permissionManager.viewUserDetail('${user.email}')">
                        <i class="fas fa-eye"></i> <span data-translate="btn.view">查看</span>
                    </button>
                    <button class="btn btn-sm btn-warning" onclick="permissionManager.editPermissions('${user.email}')">
                        <i class="fas fa-edit"></i> <span data-translate="btn.edit">编辑</span>
                    </button>
                    <button class="btn btn-sm btn-danger" onclick="permissionManager.deleteUser('${user.email}')">
                        <i class="fas fa-trash"></i> <span data-translate="btn.delete">删除</span>
                    </button>
                    ${user.status === 'pending_review' ? `
                    <button class="btn btn-sm btn-success" onclick="permissionManager.approveUser('${user.email}')">
                        <i class="fas fa-check"></i> <span data-translate="btn.approve">同意</span>
                    </button>
                    <button class="btn btn-sm btn-secondary" onclick="permissionManager.rejectUser('${user.email}')">
                        <i class="fas fa-times"></i> <span data-translate="btn.reject">不同意</span>
                    </button>` : ''}
                </div>
            </div>
        `).join('');

        container.innerHTML = userListHTML;
        this.applyI18n();
    }

    // 渲染权限列表
    renderPermissions() {
        const container = document.getElementById('permissionList');
        
        const permissionHTML = this.permissions.map(permission => `
            <div class="permission-item">
                <div class="permission-info">
                    <h4>${(window.LanguageConfig?window.LanguageConfig.getText('permissions.' + permission.key + '.name'):permission.name)}</h4>
                    <p>${(window.LanguageConfig?window.LanguageConfig.getText('permissions.' + permission.key + '.desc'):permission.description)}</p>
                </div>
                <div class="permission-actions">
                    <span class="role-badge role-${permission.level}">${(window.LanguageConfig?window.LanguageConfig.getText('level.' + permission.level):permission.level.toUpperCase())}</span>
                    <button class="btn btn-sm ${permission.enabled ? 'btn-success' : 'btn-secondary'}" 
                            onclick="permissionManager.togglePermission(${permission.id})">
                        <i class="fas fa-${permission.enabled ? 'check' : 'times'}"></i> 
                        ${permission.enabled ? (window.LanguageConfig?window.LanguageConfig.getText('modal.permission.status.active'):'启用') : (window.LanguageConfig?window.LanguageConfig.getText('modal.permission.status.inactive'):'禁用')}
                    </button>
                </div>
            </div>
        `).join('');

        container.innerHTML = permissionHTML;
        this.applyI18n();
    }

    // 更新统计信息
    updateStats() {
        const totalUsers = this.users.length;
        const activeUsers = this.users.filter(user => user.status === 'completed').length;
        const ownerUsers = this.users.filter(user => user.role === 'owner').length;
        const managerUsers = this.users.filter(user => user.role === 'property_manager').length;

        document.getElementById('totalUsers').textContent = totalUsers;
        document.getElementById('activeUsers').textContent = activeUsers;
        document.getElementById('ownerUsers').textContent = ownerUsers;
        document.getElementById('managerUsers').textContent = managerUsers;
    }

    // 搜索用户
    searchUsers() {
        const searchTerm = document.getElementById('searchInput').value.toLowerCase();
        const roleFilter = document.getElementById('roleFilter').value;

        this.filteredUsers = this.users.filter(user => {
            const matchesSearch = !searchTerm || 
                user.email.toLowerCase().includes(searchTerm) ||
                (user.full_name && user.full_name.toLowerCase().includes(searchTerm)) ||
                this.getRoleDisplayName(user.role).toLowerCase().includes(searchTerm);
            
            const matchesRole = !roleFilter || user.role === roleFilter;
            
            return matchesSearch && matchesRole;
        });

        this.renderUserList();
    }

    // 清除搜索
    clearSearch() {
        document.getElementById('searchInput').value = '';
        document.getElementById('roleFilter').value = '';
        this.filteredUsers = [...this.users];
        this.renderUserList();
    }

    // 查看用户详情
    viewUserDetail(email) {
        const user = this.users.find(u => u.email === email);
        if (!user) return;

        const userDetailHTML = `
            <div class="form-group">
                <label data-translate="label.email">用户邮箱</label>
                <input type="email" value="${user.email}" readonly>
            </div>
            <div class="form-group">
                <label data-translate="label.fullname">用户姓名</label>
                <input type="text" value="${user.full_name || '未设置'}" readonly>
            </div>
            <div class="form-group">
                <label data-translate="label.role">用户角色</label>
                <input type="text" value="${this.getRoleDisplayName(user.role)}" readonly>
            </div>
            <div class="form-group">
                <label data-translate="label.phone">手机号码</label>
                <input type="text" value="${user.phone || '未设置'}" readonly>
            </div>
            <div class="form-group">
                <label data-translate="label.did">DID标识</label>
                <input type="text" value="${user.did || '未设置'}" readonly>
            </div>
            <div class="form-group">
                <label data-translate="label.status">账户状态</label>
                <input type="text" value="${(window.LanguageConfig?window.LanguageConfig.getText('status.' + (user.status || 'unknown')):(user.status || '未知'))}" readonly>
            </div>
            <div class="form-group">
                <label data-translate="label.registerTime">注册时间</label>
                <input type="text" value="${user.created_at || '未知'}" readonly>
            </div>
            ${user.building_name ? `
            <div class="form-group">
                <label data-translate="label.building">建筑信息</label>
                <textarea readonly rows="3">${(window.LanguageConfig?window.LanguageConfig.getText('building.name.prefix'):'建筑名称')}: ${user.building_name}
${(window.LanguageConfig?window.LanguageConfig.getText('building.addr.prefix'):'建筑地址')}: ${user.building_addr || (window.LanguageConfig?window.LanguageConfig.getText('unknown'):'未设置')}
${(window.LanguageConfig?window.LanguageConfig.getText('building.type.prefix'):'建筑类型')}: ${user.building_type || (window.LanguageConfig?window.LanguageConfig.getText('unknown'):'未设置')}</textarea>
            </div>
            ` : ''}
            <div class="form-actions">
                <button type="button" class="btn btn-primary" onclick="permissionManager.editPermissions('${user.email}')">
                    <i class="fas fa-edit"></i> <span data-translate="btn.editPermission">编辑权限</span>
                </button>
                ${user.status === 'pending_review' ? `
                <button type="button" class="btn btn-success" onclick="permissionManager.approveUser('${user.email}')">
                    <i class="fas fa-check"></i> <span data-translate="btn.approve">同意</span>
                </button>
                <button type="button" class="btn btn-secondary" onclick="permissionManager.rejectUser('${user.email}')">
                    <i class="fas fa-times"></i> <span data-translate="btn.reject">不同意</span>
                </button>` : ''}
                <button type="button" class="btn btn-secondary" onclick="closeModal('userDetailModal')" data-translate="btn.close">关闭</button>
            </div>
        `;

        document.getElementById('userDetailContent').innerHTML = userDetailHTML;
        this.applyI18n();
        this.showModal('userDetailModal');
    }

    // 审核通过
    async approveUser(email) {
        try {
            const resp = await fetch('/api/admin/review-approve', {
                method: 'POST', headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email })
            });
            let result; const txt = await resp.text();
            try { result = JSON.parse(txt); } catch(e) { throw new Error(txt || (window.LanguageConfig?window.LanguageConfig.getText('server.nonjson'):'服务器返回非JSON')); }
            if (result.code === 200) {
                this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.approve.success'):'审核通过，已发送邮件通知'), 'success');
                await this.loadUsers();
                this.closeModal('userDetailModal');
            } else {
                throw new Error(result.message || (window.LanguageConfig?window.LanguageConfig.getText('msg.approve.fail'):'审批失败'));
            }
        } catch (e) {
            this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.approve.fail'):'审批失败') + ': ' + e.message, 'error');
        }
    }

    // 审核拒绝
    async rejectUser(email) {
        const reason = prompt(window.LanguageConfig?window.LanguageConfig.getText('msg.reject.prompt'):'请输入拒绝原因（可选）：') || '';
        try {
            const resp = await fetch('/api/admin/review-reject', {
                method: 'POST', headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, reason })
            });
            let result; const txt = await resp.text();
            try { result = JSON.parse(txt); } catch(e) { throw new Error(txt || (window.LanguageConfig?window.LanguageConfig.getText('server.nonjson'):'服务器返回非JSON')); }
            if (result.code === 200) {
                this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.reject.success'):'已拒绝，已发送邮件通知'), 'info');
                await this.loadUsers();
                this.closeModal('userDetailModal');
            } else {
                throw new Error(result.message || (window.LanguageConfig?window.LanguageConfig.getText('msg.reject.fail'):'拒绝失败'));
            }
        } catch (e) {
            this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.reject.fail'):'拒绝失败') + ': ' + e.message, 'error');
        }
    }

    // 编辑权限
    editPermissions(email) {
        const user = this.users.find(u => u.email === email);
        if (!user) return;

        document.getElementById('permissionUser').value = user.email;
        document.getElementById('permissionRole').value = user.role;
        document.getElementById('permissionStatus').value = user.status || 'active';
        document.getElementById('permissionNotes').value = '';

        this.showModal('permissionModal');
    }

    // 删除用户
    async deleteUser(email) {
        if (!confirm((window.LanguageConfig?window.LanguageConfig.getText('msg.delete.confirm'):'确定要删除用户 {email} 吗？此操作不可撤销！').replace('{email}', email))) {
            return;
        }

        try {
            // 这里应该调用删除用户的API
            // const response = await fetch(`/api/register/delete?email=${encodeURIComponent(email)}`, {
            //     method: 'DELETE'
            // });
            
            // 模拟删除成功
            this.users = this.users.filter(user => user.email !== email);
            this.filteredUsers = this.filteredUsers.filter(user => user.email !== email);
            this.renderUserList();
            this.updateStats();
            
            this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.delete.success'):'用户删除成功'), 'success');
        } catch (error) {
            console.error('删除用户失败:', error);
            this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.delete.fail'):'删除用户失败') + ': ' + error.message, 'error');
        }
    }

    // 切换权限状态
    togglePermission(permissionId) {
        const permission = this.permissions.find(p => p.id === permissionId);
        if (permission) {
            permission.enabled = !permission.enabled;
            this.renderPermissions();
            this.showMessage(((window.LanguageConfig?window.LanguageConfig.getText('msg.togglePermission.' + (permission.enabled ? 'enabled' : 'disabled')):'权限状态已更新')).replace('{name}', permission.name), 'success');
        }
    }

    // 保存权限更改
    async savePermissionChanges() {
        const form = document.getElementById('permissionForm');
        const formData = new FormData(form);
        
        const email = document.getElementById('permissionUser').value;
        const role = document.getElementById('permissionRole').value;
        const status = document.getElementById('permissionStatus').value;
        const notes = document.getElementById('permissionNotes').value;

        try {
            // 这里应该调用更新用户权限的API
            // const response = await fetch('/api/register/update-permissions', {
            //     method: 'POST',
            //     headers: {
            //         'Content-Type': 'application/json'
            //     },
            //     body: JSON.stringify({
            //         email: email,
            //         role: role,
            //         status: status,
            //         notes: notes
            //     })
            // });

            // 模拟更新成功
            const user = this.users.find(u => u.email === email);
            if (user) {
                user.role = role;
                user.status = status;
                this.renderUserList();
                this.updateStats();
            }

            this.closeModal('permissionModal');
            this.showMessage('权限更新成功', 'success');
        } catch (error) {
            console.error('更新权限失败:', error);
            this.showMessage('更新权限失败: ' + error.message, 'error');
        }
    }

    // 刷新数据
    async refreshData() {
        this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.refresh.start'):'正在刷新数据...'), 'info');
        await this.loadUsers();
        this.loadPermissions();
        this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.refresh.done'):'数据刷新完成'), 'success');
    }

    // 显示设置
    showSettings() {
        this.showMessage((window.LanguageConfig?window.LanguageConfig.getText('msg.settings.wip'):'设置功能开发中...'), 'info');
    }

    // 退出登录
    logout() {
        if (confirm(window.LanguageConfig?window.LanguageConfig.getText('msg.logout.confirm'):'确定要退出登录吗？')) {
            localStorage.removeItem('userInfo');
            localStorage.removeItem('chaincodeData');
            window.location.replace('/static/index.html');
        }
    }

    // 绑定事件
    bindEvents() {
        // 搜索输入框回车事件
        document.getElementById('searchInput').addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.searchUsers();
            }
        });

        // 角色筛选变化事件
        document.getElementById('roleFilter').addEventListener('change', () => {
            this.searchUsers();
        });

        // 权限表单提交事件
        document.getElementById('permissionForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.savePermissionChanges();
        });

        // 模态框点击外部关闭
        window.addEventListener('click', (e) => {
            if (e.target.classList.contains('modal')) {
                this.closeModal(e.target.id);
            }
        });
    }

    // 显示模态框
    showModal(modalId) {
        document.getElementById(modalId).style.display = 'block';
    }

    // 关闭模态框
    closeModal(modalId) {
        document.getElementById(modalId).style.display = 'none';
    }

    // 显示加载状态
    showLoading(containerId) {
        document.getElementById(containerId).innerHTML = `
            <div class="loading">
                <i class="fas fa-spinner"></i> ${(window.LanguageConfig?window.LanguageConfig.getText('loading.generic'):'加载中...')}
            </div>
        `;
        this.applyI18n();
    }

    // 显示错误信息
    showError(containerId, message) {
        document.getElementById(containerId).innerHTML = `
            <div class="empty-state">
                <i class="fas fa-exclamation-triangle"></i>
                <h3 data-translate="load.fail.title">加载失败</h3>
                <p>${message}</p>
                <button class="btn btn-primary" onclick="permissionManager.loadUsers()">
                    <i class="fas fa-refresh"></i> <span data-translate="btn.retry">重试</span>
                </button>
            </div>
        `;
        this.applyI18n();
    }

    // 应用多语言到当前渲染的动态内容
    applyI18n() {
        try {
            if (window.LanguageConfig && typeof window.LanguageConfig.applyLanguageToAllPages === 'function') {
                window.LanguageConfig.applyLanguageToAllPages();
            }
        } catch (e) {}
    }

    // 显示消息提示
    showMessage(message, type = 'info') {
        // 创建消息提示元素
        const messageDiv = document.createElement('div');
        messageDiv.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 15px 20px;
            border-radius: 8px;
            color: white;
            font-weight: 500;
            z-index: 10000;
            animation: slideIn 0.3s ease;
        `;

        // 根据类型设置颜色
        const colors = {
            success: '#27ae60',
            error: '#e74c3c',
            warning: '#f39c12',
            info: '#3498db'
        };
        messageDiv.style.backgroundColor = colors[type] || colors.info;
        messageDiv.textContent = message;

        // 添加到页面
        document.body.appendChild(messageDiv);

        // 3秒后自动移除
        setTimeout(() => {
            messageDiv.style.animation = 'slideOut 0.3s ease';
            setTimeout(() => {
                if (messageDiv.parentNode) {
                    messageDiv.parentNode.removeChild(messageDiv);
                }
            }, 300);
        }, 3000);
    }

    // 跳转到登录页面
    redirectToLogin() {
        window.location.replace('/static/index.html');
    }

    // 跳转到文件管理页面
    redirectToFileManagement() {
        window.location.replace('/static/file-management.html');
    }
}

// 全局函数
function refreshData() {
    permissionManager.refreshData();
}

function showSettings() {
    permissionManager.showSettings();
}

function logout() {
    permissionManager.logout();
}

function searchUsers() {
    permissionManager.searchUsers();
}

function clearSearch() {
    permissionManager.clearSearch();
}

function closeModal(modalId) {
    permissionManager.closeModal(modalId);
}

// 添加CSS动画
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// 初始化权限管理器
let permissionManager;
document.addEventListener('DOMContentLoaded', function() {
    permissionManager = new PermissionManager();
});
