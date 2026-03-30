// 显示用户详情
function showUserDetail() {
    const userInfo = localStorage.getItem('userInfo');
    if (!userInfo) {
        alert('未找到用户信息');
        return;
    }
    
    try {
        const user = JSON.parse(userInfo);
        const content = document.getElementById('userDetailContent');
        
        content.innerHTML = `
            <div style="padding: 20px;">
                <h4>用户基本信息</h4>
                <p><strong>姓名:</strong> ${user.FullName || user.fullName || user.Name || user.name || '用户'}</p>
                <p><strong>邮箱:</strong> ${user.Email || user.email || 'user@example.com'}</p>
                <p><strong>DID:</strong> ${user.DID || user.did || '未设置'}</p>
                <p><strong>角色:</strong> ${user.Role || user.role || '普通用户'}</p>
                <div style="margin-top: 20px; text-align: center;">
                    <button onclick="logout()" style="padding: 10px 20px; background: #f44336; color: white; border: none; border-radius: 4px; cursor: pointer;">
                        退出登录
                    </button>
                </div>
            </div>
        `;
        
        document.getElementById('userDetailModal').classList.remove('hidden');
    } catch (error) {
        console.error('解析用户信息失败:', error);
        alert('获取用户信息失败');
    }
}

// 关闭用户详情
function closeUserDetail() {
    document.getElementById('userDetailModal').classList.add('hidden');
}
