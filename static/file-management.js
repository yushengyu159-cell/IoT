// 文件管理系统 JavaScript
// 全局变量
const API_BASE_URL = 'http://47.238.159.234:8199/api';
let currentUser = null;

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 DOM加载完成，开始初始化文件管理系统');
    console.log('🔍 页面加载时的 localStorage.userInfo:', localStorage.getItem('userInfo'));
    
    try {
        initializeFileManagement();
        // 先检查用户登录状态，初始化 currentUser
        checkUserLogin();
        // 然后强制刷新用户信息
        forceRefreshUserInfo().then((result) => {
            console.log('🔄 页面加载完成后的强制刷新结果:', result);
            console.log('🔍 强制刷新后的 localStorage.userInfo:', localStorage.getItem('userInfo'));
            console.log('🔍 强制刷新后的 currentUser:', currentUser);
        });
        console.log('✅ 文件管理系统初始化完成');
    } catch (error) {
        console.error('❌ 文件管理系统初始化失败:', error);
    }
});
// 简易翻译函数（依赖 file-management.html 中的 languageData）
function tr(key, fallback) {
    try {
        const lang = localStorage.getItem('esg_language') || 'en';
        const data = (window.languageData && window.languageData[lang] && window.languageData[lang].t) || {};
        return data[key] || fallback || key;
    } catch (_) {
        return fallback || key;
    }
}


// 刷新并覆盖本地的用户信息，确保 DID/txID 与后端一致
async function refreshUserInfoFromServer() {
    try {
        console.log('🔄 开始刷新用户信息...');
        const local = JSON.parse(localStorage.getItem('userInfo') || '{}');
        console.log('📱 本地用户信息:', local);
        
        const email = local.Email || local.email;
        if (!email) {
            console.warn('⚠️ 本地用户信息中没有邮箱');
            return null;
        }
        
        console.log('📧 查询邮箱:', email);
        const resp = await fetch(`/api/register/status?email=${encodeURIComponent(email)}`);
        if (!resp.ok) {
            console.error('❌ 后端接口调用失败:', resp.status, resp.statusText);
            return null;
        }
        
        const json = await resp.json();
        console.log('📥 后端返回数据:', json);
        
        if (json && json.code === 200 && json.data) {
            const merged = { ...local, ...json.data };
            console.log('🔄 合并前:', local);
            console.log('🔄 合并后:', merged);
            
            // 兼容字段：统一 DID/txID 键名（同时兼容 did/DID 两种返回）
            if (json.data.DID && !merged.did) merged.did = json.data.DID;
            if (json.data.did && !merged.DID) merged.DID = json.data.did;
            if (json.data.did && !merged.did) merged.did = json.data.did;
            if (json.data.txID && !merged.txid) merged.txid = json.data.txID;
            if (json.data.txid && !merged.txID) merged.txID = json.data.txid;
            
            // 映射后端字段为前端字段
            const data = json.data;
            if (data.building_name) merged.BuildingName = data.building_name;
            if (data.building_addr) merged.BuildingAddr = data.building_addr;
            if (data.building_type) merged.BuildingType = data.building_type;
            if (data.property_name) merged.PropertyName = data.property_name;
            if (data.occupation) merged.Occupation = data.occupation;
            if (data.institution) merged.Institution = data.institution;
            if (data.Role && !merged.Role) merged.Role = data.Role;
            if (data.FullName && !merged.FullName) merged.FullName = data.FullName;
            if (data.Phone && !merged.Phone) merged.Phone = data.Phone;
            if (data.RegisterTime && !merged.RegisterTime) merged.RegisterTime = data.RegisterTime;

            // 强制更新 DID 字段，确保使用正确的链码 DID
            if (json.data.DID || json.data.did) {
                const ensuredDid = json.data.DID || json.data.did;
                merged.DID = ensuredDid;
                merged.did = ensuredDid;
                console.log('✅ 强制更新 DID 字段:', ensuredDid);
            }
            
            // 保存到 localStorage
            localStorage.setItem('userInfo', JSON.stringify(merged));
            console.log('💾 已保存到 localStorage:', merged);
            
            // 同时更新全局 currentUser 对象
            if (currentUser) {
                currentUser.DID = merged.DID;
                currentUser.did = merged.did;
                console.log('🌍 已更新 currentUser:', currentUser);
            } else {
                console.warn('⚠️ currentUser 对象未初始化');
            }
            
            console.log('✅ 用户信息刷新完成，DID:', merged.DID);
            return merged; // 返回更新后的用户信息
        } else {
            console.warn('⚠️ 后端返回数据格式不正确:', json);
        }
    } catch (err) {
        console.error('❌ 刷新用户信息失败:', err);
    }
    return null;
}

// 强制刷新用户信息，清理错误的 DID
async function forceRefreshUserInfo() {
    console.log('🔄 强制刷新用户信息...');
    
    try {
        const local = JSON.parse(localStorage.getItem('userInfo') || '{}');
        const email = local.Email || local.email;
        
        if (!email) {
            console.error('❌ 无法获取用户邮箱');
            return null;
        }
        
        console.log('📧 强制刷新邮箱:', email);
        
        // 直接调用后端接口获取最新信息
        const resp = await fetch(`/api/register/status?email=${encodeURIComponent(email)}`);
        if (!resp.ok) {
            throw new Error(`HTTP错误: ${resp.status}`);
        }
        
        const json = await resp.json();
        console.log('📥 强制刷新返回数据:', json);
        
        if (json && json.code === 200 && json.data) {
            // 完全覆盖本地存储，确保使用正确的 DID
            const correctUserInfo = {
                ...local,
                ...json.data,
                DID: json.data.did,        // 强制使用后端的 did 字段
                did: json.data.did,        // 同时设置 did 字段
                txID: json.data.txID || json.data.txid || '',
                txid: json.data.txID || json.data.txid || ''
            };
            // 同步映射字段
            if (json.data.building_name) correctUserInfo.BuildingName = json.data.building_name;
            if (json.data.building_addr) correctUserInfo.BuildingAddr = json.data.building_addr;
            if (json.data.building_type) correctUserInfo.BuildingType = json.data.building_type;
            if (json.data.property_name) correctUserInfo.PropertyName = json.data.property_name;
            if (json.data.occupation) correctUserInfo.Occupation = json.data.occupation;
            if (json.data.institution) correctUserInfo.Institution = json.data.institution;
            if (json.data.Role && !correctUserInfo.Role) correctUserInfo.Role = json.data.Role;
            if (json.data.FullName && !correctUserInfo.FullName) correctUserInfo.FullName = json.data.FullName;
            if (json.data.Phone && !correctUserInfo.Phone) correctUserInfo.Phone = json.data.Phone;
            if (json.data.RegisterTime && !correctUserInfo.RegisterTime) correctUserInfo.RegisterTime = json.data.RegisterTime;
            
            console.log('✅ 强制刷新后的用户信息:', correctUserInfo);
            
            // 保存到 localStorage
            localStorage.setItem('userInfo', JSON.stringify(correctUserInfo));
            
            // 更新 currentUser
            if (currentUser) {
                currentUser.DID = correctUserInfo.DID;
                currentUser.did = correctUserInfo.did;
                currentUser.txID = correctUserInfo.txID;
                currentUser.txid = correctUserInfo.txid;
            }
            
            console.log('✅ 强制刷新完成，DID:', correctUserInfo.DID);
            return correctUserInfo;
        }
    } catch (error) {
        console.error('❌ 强制刷新失败:', error);
    }
    
    return null;
}

// 初始化文件管理系统
function initializeFileManagement() {
    console.log('🔧 开始初始化文件管理系统');
    
    // 绑定文件上传表单事件
    const uploadForm = document.getElementById('fileUploadForm');
    if (uploadForm) {
        console.log('✅ 找到上传表单，绑定提交事件');
        uploadForm.addEventListener('submit', handleFileUpload);
        
        // 添加额外的调试信息
        console.log('📝 表单元素详情:', {
            id: uploadForm.id,
            action: uploadForm.action,
            method: uploadForm.method,
            elements: uploadForm.elements.length
        });
    } else {
        console.error('❌ 未找到上传表单 fileUploadForm');
        // 尝试延迟绑定
        setTimeout(() => {
            const retryForm = document.getElementById('fileUploadForm');
            if (retryForm) {
                console.log('🔄 延迟绑定成功');
                retryForm.addEventListener('submit', handleFileUpload);
            }
        }, 1000);
    }
    
    // 绑定文件选择事件
    const fileInput = document.getElementById('fileInput');
    if (fileInput) {
        console.log('✅ 找到文件输入框，绑定选择事件');
        fileInput.addEventListener('change', handleFileSelect);
    } else {
        console.error('❌ 未找到文件输入框 fileInput');
    }
    
    // 加载文件列表
    loadFileList();
    
    console.log('✅ 文件管理系统初始化完成');
}

// 检查用户登录状态
function checkUserLogin() {
    console.log('🔐 检查用户登录状态...');
    const userInfo = localStorage.getItem('userInfo');
    console.log('💾 localStorage.userInfo:', userInfo);
    
    if (userInfo) {
        try {
            currentUser = JSON.parse(userInfo);
            console.log('✅ 成功解析用户信息:', currentUser);
            console.log('🔑 用户 DID:', currentUser.DID || currentUser.did);
            displayUserInfo(currentUser);
        } catch (error) {
            console.error('❌ 解析用户信息失败:', error);
            redirectToLogin();
        }
    } else {
        console.warn('⚠️ localStorage 中没有用户信息');
        redirectToLogin();
    }
}

// 显示用户信息
function displayUserInfo(userData) {
    document.getElementById('userName').textContent = userData.Name || userData.name || tr('unknown','未知用户');
    document.getElementById('userEmail').textContent = userData.Email || userData.email || tr('unknown','未知邮箱');
}

// 文件上传相关功能
function showUploadModal() {
    console.log('🚀 显示上传弹窗...');
    
    // 调试：检查当前用户信息状态
    debugUserInfoState();
    
    const modal = document.getElementById('uploadModal');
    modal.classList.remove('hidden');
    
    // 重置表单
    document.getElementById('fileUploadForm').reset();
    document.getElementById('fileInputDisplay').innerHTML = `
        <i class="fas fa-cloud-upload-alt"></i>
        <span>${tr('clickOrDrag','点击选择文件或拖拽文件到此处')}</span>
    `;
    document.querySelector('.file-input-wrapper').classList.remove('has-file');
    
    console.log('📱 当前 currentUser:', currentUser);
    console.log('💾 当前 localStorage.userInfo:', localStorage.getItem('userInfo'));
    
    // 先刷新用户信息，确保 DID 是最新的
    forceRefreshUserInfo().then((updatedUserInfo) => {
        console.log('🔄 强制刷新完成，返回数据:', updatedUserInfo);
        
        // 使用最新的 DID 设置上传者字段
        let uploaderValue = '';
        
        // 优先使用强制刷新返回的 DID
        if (updatedUserInfo && updatedUserInfo.DID) {
            uploaderValue = updatedUserInfo.DID;
            console.log('✅ 使用强制刷新的 DID:', uploaderValue);
        } else if (currentUser && (currentUser.DID || currentUser.did)) {
            uploaderValue = currentUser.DID || currentUser.did;
            console.log('✅ 使用 currentUser 的 DID:', uploaderValue);
        } else {
            // 最后尝试从 localStorage 获取
            try {
                const local = JSON.parse(localStorage.getItem('userInfo') || '{}');
                uploaderValue = local.DID || local.did || '';
                console.log('✅ 使用 localStorage 的 DID:', uploaderValue);
            } catch (_) {
                console.warn('⚠️ 无法从 localStorage 获取 DID');
            }
        }
        
        const uploaderInput = document.getElementById('uploader');
        uploaderInput.value = uploaderValue;
        uploaderInput.setAttribute('readonly', 'readonly');
        uploaderInput.title = tr('uploaderHelp','上传者身份固定为当前登录者的 DID');
        console.log('✅ 上传弹窗设置 DID:', uploaderValue);
        
        // 验证 DID 格式
        if (uploaderValue && !/^did:example:/.test(uploaderValue)) {
            console.error('❌ 错误：上传者 DID 格式不正确:', uploaderValue);
            console.error('❌ 期望格式: did:example:...');
            
            // 如果 DID 格式仍然不正确，显示错误提示
            showStatusMessage(tr('cannotGetUploader','无法获取正确的上传者身份，请重新登录'), 'error');
            closeUploadModal();
            return;
        } else if (uploaderValue) {
            console.log('✅ DID 格式正确:', uploaderValue);
        }
    });
}

// 调试函数：检查当前用户信息状态
function debugUserInfoState() {
    console.log('🔍 === 调试用户信息状态 ===');
    console.log('📱 currentUser:', currentUser);
    console.log('💾 localStorage.userInfo:', localStorage.getItem('userInfo'));
    
    try {
        const local = JSON.parse(localStorage.getItem('userInfo') || '{}');
        console.log('🔑 本地存储的 DID 字段:');
        console.log('  - DID:', local.DID);
        console.log('  - did:', local.did);
        console.log('  - Email:', local.Email || local.email);
        console.log('  - FullName:', local.FullName || local.full_name);
        
        // 检查 DID 格式
        if (local.DID || local.did) {
            const did = local.DID || local.did;
            if (/^did:esg:user:/.test(did)) {
                console.error('❌ 检测到错误的 DID 格式:', did);
                console.error('❌ 期望格式: did:example:...');
            } else if (/^did:example:/.test(did)) {
                console.log('✅ DID 格式正确:', did);
            } else {
                console.warn('⚠️ DID 格式未知:', did);
            }
        }
    } catch (error) {
        console.error('❌ 解析 localStorage.userInfo 失败:', error);
    }
    
    // 检查是否有其他可能的 DID 存储位置
    const allKeys = Object.keys(localStorage);
    console.log('🔑 localStorage 中的所有键:', allKeys);
    
    // 查找包含 DID 信息的键
    allKeys.forEach(key => {
        if (key.toLowerCase().includes('did') || key.toLowerCase().includes('user')) {
            console.log(`🔍 检查键 ${key}:`, localStorage.getItem(key));
        }
    });
    
    console.log('�� === 调试结束 ===');
}

function closeUploadModal() {
    const modal = document.getElementById('uploadModal');
    modal.classList.add('hidden');
}

// 处理文件选择
function handleFileSelect(event) {
    const file = event.target.files[0];
    const fileDisplay = document.getElementById('fileInputDisplay');
    const fileWrapper = document.querySelector('.file-input-wrapper');
    
    if (file) {
        // 显示选中的文件信息
        fileDisplay.innerHTML = `
            <i class="fas fa-file"></i>
            <span>${file.name}</span>
            <small>${formatFileSize(file.size)}</small>
        `;
        fileWrapper.classList.add('has-file');
    } else {
        // 重置显示
        fileDisplay.innerHTML = `
            <i class="fas fa-cloud-upload-alt"></i>
            <span>${tr('clickOrDrag','点击选择文件或拖拽文件到此处')}</span>
        `;
        fileWrapper.classList.remove('has-file');
    }
}

// 处理文件上传
async function handleFileUpload(event) {
    console.log('🚀 文件上传事件触发');
    event.preventDefault();
    
    const fileInput = document.getElementById('fileInput');
    const fileDesc = document.getElementById('fileDesc').value;
    const uploaderInputEl = document.getElementById('uploader');
    let uploader = uploaderInputEl.value || '';
    
    console.log('📤 上传前检查 - 文件:', fileInput.files[0]?.name);
    console.log('📤 上传前检查 - 描述:', fileDesc);
    console.log('📤 上传前检查 - 上传者:', uploader);
    console.log('📤 上传前检查 - currentUser:', currentUser);
    console.log('📤 上传前检查 - localStorage.userInfo:', localStorage.getItem('userInfo'));
    
    if (!fileInput.files[0]) {
        showStatusMessage(tr('pleaseSelectFile','请选择要上传的文件'), 'error');
        return;
    }
    
    if (!fileDesc.trim()) {
        showStatusMessage(tr('pleaseEnterDesc','请输入文件描述'), 'error');
        return;
    }
    
    if (!uploader.trim()) {
        showStatusMessage(tr('pleaseEnterUploader','请输入上传者信息'), 'error');
        return;
    }
    
    // 强制使用 DID 作为上传者身份
    if (!/^did:/.test(uploader)) {
        console.log('⚠️ 上传者不是 DID 格式，尝试强制设置...');
        try {
            const local = JSON.parse(localStorage.getItem('userInfo') || '{}');
            const enforced = local.DID || local.did || '';
            if (enforced) {
                uploader = enforced;
                uploaderInputEl.value = enforced;
                console.log('✅ 已强制设置为 DID:', enforced);
            }
        } catch (_) {
            console.warn('⚠️ 强制设置 DID 失败');
        }
    }
    
    // 简化：直接使用当前uploader，不进行DID格式验证
    console.log('📤 使用上传者:', uploader);

    // 如果uploader为空，尝试从localStorage获取
    if (!uploader) {
        try {
            const local = JSON.parse(localStorage.getItem('userInfo') || '{}');
            uploader = local.Email || local.email || local.DID || local.did || '';
            console.log('✅ 从localStorage获取上传者:', uploader);
        } catch (e) {
            console.warn('⚠️ 无法获取上传者信息');
        }
    }

    console.log('📤 最终上传者:', uploader);
    
    try {
        showLoading(true);
        
        const formData = new FormData();
        formData.append('file', fileInput.files[0]);
        formData.append('desc', fileDesc);
        formData.append('uploader', uploader);
        
        console.log('开始上传文件:', fileInput.files[0].name);
        console.log('📤 准备发送到后端的数据:', {
            fileName: fileInput.files[0].name,
            fileSize: fileInput.files[0].size,
            description: fileDesc,
            uploader: uploader
        });
        console.log('📤 请求URL:', `${API_BASE_URL}/esg/upload-encrypted`);
        
        // 使用后端ESG路由，上传后会写入数据库
        console.log('📤 切换为加密上传接口: /esg/upload-encrypted');
        const response = await fetch(`${API_BASE_URL}/esg/upload-encrypted`, {
            method: 'POST',
            body: formData
        });
        
         console.log('📥 后端响应状态:', response.status);
         console.log('📥 后端响应头:', response.headers);
        
         if (!response.ok) {
             throw new Error(`HTTP错误: ${response.status} ${response.statusText}`);
         }
        
         const result = await response.json();
         console.log('📥 后端响应数据:', result);
        
        if (result.code === 200) {
            showStatusMessage(tr('uploadSuccess', '加密分片上传成功！'), 'success');
            console.log('✅ [上传成功] 即将关闭模态框并刷新列表');
            
            // 先关闭模态框
            closeUploadModal();
            
            // 立即执行第一次刷新（尝试）
            console.log('🔄 [上传成功] 立即执行第一次刷新尝试...');
            setTimeout(() => {
                console.log('⚡ [上传成功] 1秒后执行刷新...');
                loadFileList();
            }, 1000);
            
            // 1.5秒后执行第二次刷新（确保后端写入完成）
            setTimeout(() => {
                console.log('🔄 [上传成功] 2秒后执行确认刷新...');
                loadFileList();
            }, 2000);
        } else {
            showStatusMessage(result.message || tr('uploadFailed','上传失败'), 'error');
        }
        
    } catch (error) {
        console.error('文件上传错误:', error);
        showStatusMessage(tr('uploadFailed','上传失败') + ': ' + error.message, 'error');
    } finally {
        showLoading(false);
    }
}

// 文件列表管理功能
function loadFileList() {
    console.log('🔄 [loadFileList] 开始加载文件列表...');
    fetchFileList()
        .then(files => {
            console.log('📥 [loadFileList] 获取到文件数量:', files.length);
            displayFileList(files);
        })
        .catch(error => {
            console.error('加载文件列表失败:', error);
            showStatusMessage(tr('loadingListFailed','加载文件列表失败'), 'error');
        });
}

async function fetchFileList() {
    try {
        // 获取当前登录用户信息
        const userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}');
        const userEmail = userInfo.Email || userInfo.email;
        
        if (!userEmail) {
            console.error('未找到用户邮箱信息');
            showStatusMessage(tr('pleaseLogin','请先登录'), 'error');
            return [];
        }
        
        // 调用文件列表API，传递用户邮箱参数实现用户隔离
        const response = await fetch(`${API_BASE_URL}/esg/list?userEmail=${encodeURIComponent(userEmail)}`);
        if (!response.ok) {
            throw new Error(`HTTP错误: ${response.status}`);
        }
        
        const result = await response.json();
        if (result.code === 200) {
            return result.data || [];
        } else {
            throw new Error(result.message || '获取文件列表失败');
        }
    } catch (error) {
        console.error('获取文件列表错误:', error);
        return [];
    }
}

function displayFileList(files) {
    console.log('🎨 [displayFileList] 开始渲染文件列表, 文件数:', files.length);
    const fileListContainer = document.getElementById('fileList');
    console.log('📍 [displayFileList] 找到fileList容器:', fileListContainer);
    
    if (!fileListContainer) {
        console.error('❌ [displayFileList] 无法找到fileList元素!');
        return;
    }
    
    // 强制触发重排以确保DOM更新
    console.log('🔄 [displayFileList] 准备更新DOM...');
    
    if (!files || files.length === 0) {
        fileListContainer.innerHTML = `
            <div class="empty-state">
                <div class="empty-hero">
                    <svg viewBox="0 0 200 140" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                        <defs>
                            <linearGradient id="g2" x1="0" x2="1">
                                <stop offset="0" stop-color="#cfead4"/>
                                <stop offset="1" stop-color="#a6d9ad"/>
                            </linearGradient>
                        </defs>
                        <rect x="20" y="60" width="160" height="60" rx="12" fill="url(#g2)" opacity="0.35"/>
                        <circle cx="80" cy="52" r="18" fill="#2e7d32" opacity="0.15"/>
                        <circle cx="120" cy="52" r="12" fill="#2e7d32" opacity="0.15"/>
                        <rect x="60" y="82" width="20" height="16" rx="4" fill="#2e7d32" opacity="0.2"/>
                        <rect x="94" y="82" width="20" height="16" rx="4" fill="#2e7d32" opacity="0.2"/>
                        <rect x="128" y="82" width="20" height="16" rx="4" fill="#2e7d32" opacity="0.2"/>
                    </svg>
                    <div class="empty-title">${tr('emptyTitle','Drag & Drop or Click to Upload Building Data')}</div>
                    <div class="empty-subtext">${tr('emptySubtext','Instant ESG Rating')}</div>
                    <ul class="empty-bullets">
                        <li>${tr('supportedFormats','Supported formats: PDF/DOC/DOCX/CSV (Max 50MB)')}</li>
                        <li>${tr('includeData','Include green building certification, energy & carbon')}</li>
                    </ul>
                </div>
            </div>
        `;
        return;
    }
    
    // 按上传时间排序，最新的文件在最上面
    const sortedFiles = files.sort((a, b) => {
        const timeA = a.uploadAt || a.UploadAt || '';
        const timeB = b.uploadAt || b.UploadAt || '';
        
        // 如果时间格式是 "2025-08-16 01:54:53" 或 "2025 08 16 01:54:59"
        if (timeA && timeB) {
            // 统一时间格式，将 "2025 08 16 01:54:59" 转换为 "2025-08-16 01:54:59"
            const normalizedTimeA = timeA.replace(/(\d{4})\s+(\d{2})\s+(\d{2})/, '$1-$2-$3');
            const normalizedTimeB = timeB.replace(/(\d{4})\s+(\d{2})\s+(\d{2})/, '$1-$2-$3');
            
            // 转换为Date对象进行比较
            const dateA = new Date(normalizedTimeA);
            const dateB = new Date(normalizedTimeB);
            
            // 最新的时间在前（降序）
            return dateB - dateA;
        }
        
        // 如果无法解析时间，保持原有顺序
        return 0;
    });
    
    const fileListHTML = sortedFiles.map(file => {
        const filename = file.filename || file.Filename || tr('unknown','未知文件');
        const fileIcon = getFileIcon(filename);
        const uploader = file.uploader || file.Uploader || tr('unknown','未知');
        const uploadTime = file.uploadAt || file.UploadAt || tr('unknown','未知');
        const description = file.desc || file.Desc || tr('noDescription','无描述');
        
        return `
        <div class="file-item" data-cid="${file.cid || file.CID || ''}" data-filename="${filename}">
            <div class="file-info">
                <div class="file-icon">
                    ${fileIcon}
                </div>
                <div class="file-details">
                    <h4 class="file-name" title="${filename}">${filename}</h4>
                    <p class="file-desc" title="${description}">${description || tr('noDescription','无描述')}</p>
                    <p class="file-meta">
                        <span class="uploader">${tr('uploaderNameLabel','上传者')}: ${uploader}</span>
                        <span class="upload-time">${tr('uploadTimeLabel','上传时间')}: ${uploadTime}</span>
                    </p>
                </div>
            </div>
            <div class="file-actions">
                <button class="btn btn-info btn-sm" onclick="viewFileDetail('${file.cid || file.CID || ''}')">
                    <i class="fas fa-info-circle"></i> ${tr('viewDetails','文件详情')}
                </button>
                <button class="btn btn-primary btn-sm" onclick="downloadFile('${file.cid || file.CID || ''}', '${filename}')">
                    <i class="fas fa-download"></i> ${tr('downloadBtn','下载')}
                </button>
            </div>
        </div>
        `;
    }).join('');
    
    fileListContainer.innerHTML = fileListHTML;
}

// 格式化CID列表显示
function formatCidList(cidsJson) {
    try {
        if (!cidsJson || cidsJson === '[]') {
            return '无分片信息';
        }
        
        const cids = JSON.parse(cidsJson);
        if (!Array.isArray(cids) || cids.length === 0) {
            return '无分片信息';
        }
        
        return cids.map((cid, index) => 
            `分片${index + 1}: ${cid}`
        ).join('; ');
        
    } catch (error) {
        console.error('解析CID列表失败:', error);
        return '解析失败';
    }
}

// 获取CID数量
function getCidCount(cidsJson) {
    try {
        if (!cidsJson || cidsJson === '[]') {
            return '0';
        }
        
        const cids = JSON.parse(cidsJson);
        if (Array.isArray(cids)) {
            return cids.length.toString();
        }
        
        return '0';
    } catch (error) {
        console.error('获取CID数量失败:', error);
        return '0';
    }
}



// 根据文件名获取对应的文件图标
function getFileIcon(filename) {
    if (!filename) return '<i class="fas fa-file"></i>';
    
    const extension = filename.toLowerCase().split('.').pop();
    
    switch (extension) {
        case 'pdf':
            return '<i class="fas fa-file-pdf" style="color: #e74c3c;"></i>';
        case 'doc':
        case 'docx':
            return '<i class="fas fa-file-word" style="color: #3498db;"></i>';
        case 'xls':
        case 'xlsx':
            return '<i class="fas fa-file-excel" style="color: #27ae60;"></i>';
        case 'ppt':
        case 'pptx':
            return '<i class="fas fa-file-powerpoint" style="color: #e67e22;"></i>';
        case 'txt':
            return '<i class="fas fa-file-alt" style="color: #95a5a6;"></i>';
        case 'jpg':
        case 'jpeg':
        case 'png':
        case 'gif':
        case 'bmp':
            return '<i class="fas fa-file-image" style="color: #9b59b6;"></i>';
        case 'mp3':
        case 'wav':
        case 'flac':
            return '<i class="fas fa-file-audio" style="color: #f39c12;"></i>';
        case 'mp4':
        case 'avi':
        case 'mov':
        case 'wmv':
            return '<i class="fas fa-file-video" style="color: #e74c3c;"></i>';
        case 'zip':
        case 'rar':
        case '7z':
            return '<i class="fas fa-file-archive" style="color: #34495e;"></i>';
        default:
            return '<i class="fas fa-file" style="color: #7f8c8d;"></i>';
    }
}

function refreshFileList() {
    showStatusMessage('正在刷新文件列表...', 'info');
    loadFileList();
}

// 文件详情查看和下载功能
function viewFileDetail(cid) {
    fetchFileDetail(cid)
        .then(fileDetail => {
            displayFileDetail(fileDetail);
        })
        .catch(error => {
            console.error('获取文件详情失败:', error);
            showStatusMessage(tr('uploadFailed','获取文件详情失败'), 'error');
        });
}

async function fetchFileDetail(cid) {
    try {
        const response = await fetch(`${API_BASE_URL}/esg/query?cid=${cid}`);
        if (!response.ok) {
            throw new Error(`HTTP错误: ${response.status}`);
        }
        
        const result = await response.json();
        if (result.code === 200) {
            return result.data;
        } else {
            throw new Error(result.message || '获取文件详情失败');
        }
    } catch (error) {
        console.error('获取文件详情错误:', error);
        throw error;
    }
}

// 显示文件详情
function displayFileDetail(fileDetail) {
    const modal = document.getElementById('fileDetailModal');
    const content = modal.querySelector('.modal-body');
    
    // 构建详情内容
    let detailHTML = `
        <div class="file-detail-section">
            <h6 class="section-title">${tr('basicInfo','基础信息')}</h6>
            <div class="detail-grid">
                <div class="detail-item">
                    <span class="label">${tr('fileName','文件名')}:</span>
                    <span class="value">${fileDetail.Filename || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('fileDescLabel','文件描述')}:</span>
                    <span class="value">${fileDetail.Desc || tr('noDescription','无描述')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('uploaderNameLabel','上传者')}:</span>
                    <span class="value">${fileDetail.uploaderName || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('uploaderIdentityLabel','上传者身份')}:</span>
                    <span class="value">${fileDetail.uploaderDid || fileDetail.Uploader || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('uploadTimeLabel','上传时间')}:</span>
                    <span class="value">${fileDetail.UploadAt || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('txHashLabel','交易哈希')}:</span>
                    <span class="value">${fileDetail.Txid || tr('unknown','未知')}</span>
                </div>
            </div>
        </div>
    `;
    
    // 如果有分片信息，显示分片详情
    if (fileDetail.ChunkCount && fileDetail.ChunkCount > 0) {
        detailHTML += `
            <div class="file-detail-section">
                <h6 class="section-title">${tr('shardInfo','分片信息')}</h6>
                <div class="detail-grid">
                    <div class="detail-item">
                        <span class="label">${tr('shardCountLabel','分片数量')}:</span>
                        <span class="value">${fileDetail.ChunkCount}</span>
                    </div>
                    <div class="detail-item">
                        <span class="label">${tr('shardSizeLabel','分片大小')}:</span>
                        <span class="value">${formatFileSize(fileDetail.ChunkSize || 0)}</span>
                    </div>
                    <div class="detail-item">
                        <span class="label">${tr('totalFileSizeLabel','总文件大小')}:</span>
                        <span class="value">${formatFileSize(fileDetail.FileSize || 0)}</span>
                    </div>
                </div>
            </div>
        `;
    }
    
    // 如果有加密信息，显示加密详情
    if (fileDetail.EncryptionKey && fileDetail.IV) {
        detailHTML += `
            <div class="file-detail-section">
                <h6 class="section-title">${tr('encryptionInfo','加密信息')}</h6>
                <div class="detail-grid">
                    <div class="detail-item">
                        <span class="label">${tr('encryptionKeyLabel','加密密钥')}:</span>
                        <span class="value">${fileDetail.EncryptionKey}</span>
                    </div>
                    <div class="detail-item">
                        <span class="label">${tr('ivLabel','初始化向量')}:</span>
                        <span class="value">${fileDetail.IV}</span>
                    </div>
                    <div class="detail-item">
                        <span class="label">${tr('cipherSampleLabel','密文样本')}:</span>
                        <span class="value">${fileDetail.CipherSample || tr('unknown','无')}</span>
                    </div>
                </div>
            </div>
        `;
    }
    
    // 如果有上传时间信息，显示性能详情
    if (fileDetail.TotalTime) {
        detailHTML += `
            <div class="file-detail-section">
                <h6 class="section-title">${tr('uploadPerformance','上传性能')}</h6>
                <div class="detail-grid">
                    <div class="detail-item">
                        <span class="label">${tr('startTimeLabel','开始时间')}:</span>
                        <span class="value">${fileDetail.UploadStartTime || tr('unknown','未知')}</span>
                    </div>
                    <div class="detail-item">
                        <span class="label">${tr('endTimeLabel','结束时间')}:</span>
                        <span class="value">${fileDetail.UploadEndTime || tr('unknown','未知')}</span>
                    </div>
                    <div class="detail-item">
                        <span class="label">${tr('totalDurationLabel','总耗时')}:</span>
                        <span class="value">${fileDetail.TotalTime}</span>
                    </div>
                </div>
            </div>
        `;
    }
    
    content.innerHTML = detailHTML;
    
    // 显示模态框 - 使用 classList 方法，与 CSS 保持一致
    modal.classList.remove('hidden');
    
    // 设置下载按钮的文件信息
    const downloadBtn = modal.querySelector('.btn-primary');
    if (downloadBtn) {
        downloadBtn.onclick = () => downloadFile(fileDetail.CID, fileDetail.Filename);
    }
}

function closeFileDetail() {
    document.getElementById('fileDetailModal').classList.add('hidden');
}

// 显示用户详情模态框
async function showUserDetail() {
    // 先强制刷新一次，确保拿到后端的最新 DID/txID
    await refreshUserInfoFromServer();
    const modal = document.getElementById('userDetailModal');
    const content = document.getElementById('userDetailContent');

    // 优先从后端获取最新用户信息，失败再回退 localStorage
    let userInfo = JSON.parse(localStorage.getItem('userInfo') || '{}');
    try {
        const email = userInfo.Email || userInfo.email;
        if (email) {
            const resp = await fetch(`/api/register/status?email=${encodeURIComponent(email)}`);
            if (resp.ok) {
                const result = await resp.json();
                if (result && result.code === 200 && result.data) {
                    userInfo = Object.assign({}, userInfo, result.data); // 合并后端信息
                }
            }
        }
    } catch (e) {
        console.warn('获取后端用户状态失败，使用本地缓存', e);
    }

    // 统一键名与缺省值：如果缺少 txid，则用 DID 兜底（链码这里 txID 与 DID 一致的场景）
    if (!userInfo.txid && userInfo.txID) userInfo.txid = userInfo.txID;
    if (!userInfo.txID && userInfo.txid) userInfo.txID = userInfo.txid;
    if (!userInfo.DID && userInfo.did) userInfo.DID = userInfo.did;
    if (!userInfo.did && userInfo.DID) userInfo.did = userInfo.DID;

    // 选择更可信的 DID：优先链码/数据库 did（以 did: 开头且不是 did:esg:user）
    const candidateDids = [
        userInfo.DID,
        userInfo.did,
        userInfo.chaincodeData && userInfo.chaincodeData.did,
    ].filter(Boolean);
    let displayDid = candidateDids.find(d => /^did:/.test(d) && !/^did:esg:user/.test(d))
        || candidateDids.find(d => /^did:/.test(d))
        || candidateDids[0]
        || tr('unknown','未知');

    // 区块链信息优先取 chaincodeData，其次取平铺字段
    const cc = userInfo.chaincodeData || {};
    const chaincodeStatus = cc.status || userInfo.chaincodeStatus || userInfo.status || tr('unknown','未知');
    const chaincodeTimestamp = cc.timestamp || userInfo.timestamp || (userInfo.timeStats && (userInfo.timeStats.endTime || userInfo.timeStats.startTime)) || tr('unknown','未知');
    const chaincodeMessage = cc.message || userInfo.message || tr('unknown','未知');

    // 若角色相关信息为空，增加一次 DID 兜底查询并合并
    const isOwner = (userInfo.Role || userInfo.role) === 'owner';
    const missingOwnerFields = isOwner && (!userInfo.BuildingName && !userInfo.buildingName);
    if (missingOwnerFields && typeof displayDid === 'string' && displayDid.startsWith('did:')) {
        try {
            const respByDid = await fetch(`/api/register/status?did=${encodeURIComponent(displayDid)}`);
            if (respByDid.ok) {
                const resultByDid = await respByDid.json();
                if (resultByDid && resultByDid.code === 200 && resultByDid.data) {
                    userInfo = Object.assign({}, userInfo, resultByDid.data);
                    // 同步本地缓存，便于后续页面使用
                    localStorage.setItem('userInfo', JSON.stringify(userInfo));
                }
            }
        } catch (e) {
            console.warn('按DID兜底查询角色相关信息失败:', e);
        }
    }

    // 选择 tx：优先使用链码返回；过滤掉形如 tx_123 的模拟值；缺失时用 DID 去掉前缀兜底
    const chaincodeTx = userInfo.chaincodeData && (userInfo.chaincodeData.txID || userInfo.chaincodeData.txid);
    const txCandidates = [chaincodeTx, userInfo.txID, userInfo.txid].filter(Boolean).filter(tx => !/^tx_\d+$/.test(tx));
    let displayTx = txCandidates[0] || (typeof displayDid === 'string' ? displayDid.replace(/^did:[^:]+:/, '') : tr('unknown','未知'));

    const userDetailHTML = `
        <div class="user-detail-content">
            <div class="detail-section">
                <h4>${tr('basicInfoTitle','基本信息')}</h4>
                <div class="detail-item">
                    <span class="label">${tr('nameLabel','姓名')}:</span>
                    <span class="value">${userInfo.FullName || userInfo.fullName || userInfo.Name || userInfo.name || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('phoneLabel','手机号码')}:</span>
                    <span class="value">${userInfo.Phone || userInfo.phone || tr('unknown','未填写')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('emailLabel','邮箱地址')}:</span>
                    <span class="value">${userInfo.Email || userInfo.email || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('roleLabel','角色')}:</span>
                    <span class="value">${userInfo.Role || userInfo.role || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('ageLabel','年龄')}:</span>
                    <span class="value">${userInfo.Age || userInfo.age || tr('unknown','未知')}</span>
                </div>
            </div>
            
            <div class="detail-section">
                <h4>${tr('roleRelatedTitle','角色相关信息')}</h4>
                ${userInfo.Role === 'owner' ? `
                <div class="detail-item">
                    <span class="label">${tr('buildingNameLabel','建筑名称')}:</span>
                    <span class="value">${userInfo.BuildingName || userInfo.buildingName || userInfo.building_name || tr('unknown','未填写')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('buildingAddrLabel','建筑地址')}:</span>
                    <span class="value">${userInfo.BuildingAddr || userInfo.buildingAddr || userInfo.building_addr || tr('unknown','未填写')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('buildingTypeLabel','建筑类型')}:</span>
                    <span class="value">${userInfo.BuildingType || userInfo.buildingType || userInfo.building_type || tr('unknown','未知')}</span>
                </div>
                ` : ''}
                ${userInfo.Role === 'property_manager' ? `
                <div class="detail-item">
                    <span class="label">${tr('propertyNameLabel','物业名称')}:</span>
                    <span class="value">${userInfo.PropertyName || userInfo.propertyName || userInfo.property_name || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('occupationLabel','职业')}:</span>
                    <span class="value">${userInfo.Occupation || userInfo.occupation || tr('unknown','未知')}</span>
                </div>
                ` : ''}
                ${userInfo.Role === 'institution' ? `
                <div class="detail-item">
                    <span class="label">${tr('institutionLabel','机构名称')}:</span>
                    <span class="value">${userInfo.Institution || userInfo.institution || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('occupationLabel','职业')}:</span>
                    <span class="value">${userInfo.Occupation || userInfo.occupation || tr('unknown','未知')}</span>
                </div>
                ` : ''}
            </div>
            
            <div class="detail-section">
                <h4>${tr('didInfoTitle','DID身份信息')}</h4>
                <div class="detail-item">
                    <span class="label">${tr('didIdentifierLabel','DID标识符')}:</span>
                    <span class="value">${displayDid}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('registerTimeLabel','注册时间')}:</span>
                    <span class="value">${userInfo.RegisterTime || userInfo.registerTime || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('identityStatusLabel','身份状态')}:</span>
                    <span class="value status-success">${tr('verifiedText','已验证')}</span>
                </div>
            </div>
            
            <div class="detail-section">
                <h4>${tr('chainInfoTitle','区块链信息')}</h4>
                <div class="detail-item">
                    <span class="label">${tr('chainStatusLabel','链上状态')}:</span>
                    <span class="value status-success">${tr('onChainText','已上链')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('assetIdLabel','资产ID')}:</span>
                    <span class="value">${userInfo.assetID || userInfo.Email || userInfo.email || tr('unknown','未知')}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('txHashLabel','交易哈希')}:</span>
                    <span class="value">${displayTx}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('chaincodeStatusLabel','链码状态')}:</span>
                    <span class="value">${chaincodeStatus}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('chainTimeLabel','上链时间')}:</span>
                    <span class="value">${chaincodeTimestamp}</span>
                </div>
                <div class="detail-item">
                    <span class="label">${tr('chainMsgLabel','链码消息')}:</span>
                    <span class="value">${chaincodeMessage}</span>
                </div>
            </div>
            
            <div class="detail-section">
                <div class="logout-section">
                    <button class="btn btn-danger logout-btn" onclick="handleLogout()">
                        <i class="fas fa-sign-out-alt"></i> ${tr('logoutText','退出登录')}
                    </button>
                </div>
            </div>
        </div>
    `;
    
    content.innerHTML = userDetailHTML;
    modal.classList.remove('hidden');
}

// 关闭用户详情模态框
function closeUserDetail() {
    const modal = document.getElementById('userDetailModal');
    modal.classList.add('hidden');
}

// 处理退出登录
function handleLogout() {
    // 显示确认对话框
    if (confirm(tr('confirmLogout','确定要退出登录吗？'))) {
        // 清除本地存储的用户信息
        localStorage.removeItem('userInfo');
        localStorage.removeItem('isLoggedIn');
        
        // 关闭用户详情模态框
        closeUserDetail();
        
        // 立即跳转到首页（index），避免浏览器历史返回到旧 register 页面
        try {
            window.location.replace('/static/index.html');
        } catch (_e) {
            window.location.href = '/static/index.html';
        }
        // 兜底：再强制一次
        setTimeout(() => {
            if (!/index\.html$/i.test(window.location.pathname)) {
                window.location.href = '/static/index.html';
            }
        }, 150);
    }
}

async function downloadFile(cid, filename = '') {
    try {
        console.log('开始下载文件，CID:', cid, '文件名:', filename);
        showStatusMessage(tr('preparingDownloadInfo','正在准备下载...'), 'info');
        
        // 首先查询文件详情，判断是否为加密文件
        console.log('正在查询文件详情...');
        const fileDetail = await getFileDetail(cid);
        console.log('文件详情:', fileDetail);
        
        if (fileDetail && fileDetail.EncryptionKey && fileDetail.AllCIDs) {
            console.log('检测到加密文件，使用解密下载');
            // 加密文件，使用解密下载接口
            await downloadEncryptedFile(fileDetail, filename);
        } else {
            console.log('检测到普通文件，使用普通下载');
            // 普通文件，使用IPFS下载接口
            await downloadNormalFile(cid, filename);
        }
        
    } catch (error) {
        console.error('下载文件失败:', error);
        showStatusMessage(tr('downloadFailedPrefix','下载失败: ') + error.message, 'error');
    }
}

// 获取文件详情
async function getFileDetail(cid) {
    try {
        console.log('正在调用文件详情API:', `${API_BASE_URL}/esg/query?cid=${cid}`);
        const response = await fetch(`${API_BASE_URL}/esg/query?cid=${cid}`);
        console.log('API响应状态:', response.status);
        
        const result = await response.json();
        console.log('API响应结果:', result);
        console.log('API响应数据类型:', typeof result);
        console.log('API响应数据结构:', JSON.stringify(result, null, 2));
        
        if (result.code === 200) {
            console.log('返回的文件详情数据:', result.data);
            console.log('文件详情数据类型:', typeof result.data);
            return result.data;
        } else {
            throw new Error(result.message || '获取文件详情失败');
        }
    } catch (error) {
        console.error('获取文件详情失败:', error);
        throw error;
    }
}

// 下载普通文件
async function downloadNormalFile(cid, originalFilename = '') {
    try {
        console.log('开始下载普通文件，CID:', cid, '原始文件名:', originalFilename);
        showStatusMessage(tr('downloadingNormalInfo','正在下载普通文件...'), 'info');
        
        // 使用后端IPFS下载接口
        const downloadUrl = `${API_BASE_URL}/ipfs/download?cid=${cid}`;
        console.log('调用下载API:', downloadUrl);
        
        const response = await fetch(downloadUrl);
        console.log('下载响应状态:', response.status);
        console.log('下载响应头:', response.headers);
        
        if (!response.ok) {
            throw new Error(`下载失败: ${response.status} ${response.statusText}`);
        }
        
        // 优先使用传入的原始文件名，如果没有则从响应头获取
        let filename = originalFilename || `file_${cid}`;
        const contentDisposition = response.headers.get('Content-Disposition');
        if (!originalFilename && contentDisposition) {
            const filenameMatch = contentDisposition.match(/filename=(.+)/);
            if (filenameMatch) {
                filename = filenameMatch[1].replace(/"/g, '');
            }
        }
        console.log('最终下载文件名:', filename);
        
        // 优化blob创建过程，使用流式处理
        console.log('正在创建文件blob...');
        showStatusMessage(tr('processingFileData','正在处理文件数据，请稍候...'), 'info');
        
        // 使用流式读取优化大文件处理
        const blob = await createOptimizedBlob(response);
        console.log('文件blob大小:', blob.size, '字节');
        
        // 验证文件内容
        if (blob.size === 0) {
            throw new Error('下载的文件为空');
        }
        
        // 快速下载文件，不进行文件头验证以提高性能
        console.log('正在触发文件下载...');
        
        // 使用setTimeout让UI更新，然后开始下载
        setTimeout(async () => {
            try {
                console.log('准备调用downloadBlob，参数:', { filename, blobSize: blob.size });
                await downloadBlob(blob, filename);
                console.log('文件下载完成');
                showStatusMessage(tr('downloadCompleted','下载完成！'), 'success');
            } catch (error) {
                console.error('下载失败:', error);
                showStatusMessage(tr('downloadFailedPrefix','下载失败: ') + error.message, 'error');
            }
        }, 100);
        
    } catch (error) {
        console.error('普通文件下载失败:', error);
        throw error;
    }
}

// 下载加密文件
async function downloadEncryptedFile(fileDetail, originalFilename = '') {
    try {
        console.log('开始下载加密文件，文件详情:', fileDetail, '原始文件名:', originalFilename);
        showStatusMessage(tr('processingFileData','正在下载并解密文件...'), 'info');
        
        // 解析分片CIDs
        let cids = [];
        try {
            console.log('正在解析分片CIDs:', fileDetail.AllCIDs);
            cids = JSON.parse(fileDetail.AllCIDs);
            console.log('解析后的分片CIDs:', cids);
        } catch (e) {
            console.error('解析分片信息失败:', e);
            throw new Error('解析分片信息失败');
        }
        
        // 准备下载参数
        console.log('文件详情字段:', Object.keys(fileDetail));
        console.log('EncryptionKey值:', fileDetail.EncryptionKey);
        console.log('IV值:', fileDetail.IV);
        console.log('AllCIDs值:', fileDetail.AllCIDs);
        
        // 尝试多种可能的字段名
        const encryptionKey = fileDetail.EncryptionKey || fileDetail.encryptionKey || fileDetail.key || fileDetail.Key;
        const iv = fileDetail.IV || fileDetail.iv || fileDetail.IVBase64 || fileDetail.ivBase64;
        
        console.log('最终使用的EncryptionKey:', encryptionKey);
        console.log('最终使用的IV:', iv);
        
        if (!encryptionKey || !iv) {
            throw new Error('缺少加密密钥或IV，无法下载加密文件');
        }
        
        const downloadData = {
            cids: cids,
            key: encryptionKey,
            iv: iv
        };
        console.log('准备下载参数:', downloadData);
        
        // 调用解密下载接口
        const downloadUrl = `${API_BASE_URL}/ipfs/download-decrypted`;
        console.log('调用解密下载API:', downloadUrl);
        
        const response = await fetch(downloadUrl, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(downloadData)
        });
        
        console.log('解密下载响应状态:', response.status);
        console.log('解密下载响应头:', response.headers);
        
        if (!response.ok) {
            throw new Error(`解密下载失败: ${response.status} ${response.statusText}`);
        }
        
        // 优先使用传入的原始文件名，如果没有则使用文件详情中的文件名
        const finalFilename = originalFilename || fileDetail.Filename || `decrypted_${fileDetail.CID}`;
        console.log('最终解密后文件名:', finalFilename);
        
        // 优化blob创建过程，使用流式处理
        console.log('正在创建解密文件blob...');
        showStatusMessage(tr('processingFileData','正在处理文件数据，请稍候...'), 'info');
        
        // 使用流式读取优化大文件处理
        const blob = await createOptimizedBlob(response);
        console.log('解密文件blob大小:', blob.size, '字节');
        
        // 验证解密后的文件内容
        if (blob.size === 0) {
            throw new Error('解密后的文件为空');
        }
        
        // 快速下载文件，不进行文件头验证以提高性能
        console.log('正在触发文件下载...');
        
        // 使用setTimeout让UI更新，然后开始下载
        setTimeout(async () => {
            try {
                console.log('准备调用downloadBlob（加密文件），参数:', { finalFilename, blobSize: blob.size });
                await downloadBlob(blob, finalFilename);
                console.log('加密文件下载并解密完成');
                showStatusMessage(tr('decryptedDownloadSuccess','加密文件下载并解密成功'), 'success');
            } catch (error) {
                console.error('下载失败:', error);
                showStatusMessage(tr('downloadFailedPrefix','下载失败: ') + error.message, 'error');
            }
        }, 100);
        
    } catch (error) {
        console.error('加密文件下载失败:', error);
        throw error;
    }
}

/**
 * 优化的Blob创建函数 - 提升大文件处理性能
 */
async function createOptimizedBlob(response) {
    const contentLength = response.headers.get('content-length');
    const fileSize = contentLength ? parseInt(contentLength) : 0;
    
    // 对于大文件（>5MB），使用分块读取优化性能
    if (fileSize > 5 * 1024 * 1024) {
        console.log('检测到大文件，使用分块读取优化...');
        return await createChunkedBlob(response, fileSize);
    } else {
        // 小文件直接使用blob()
        console.log('小文件，使用标准blob()方法');
        return await response.blob();
    }
}

/**
 * 分块读取大文件，提升性能
 */
async function createChunkedBlob(response, totalSize) {
    // 根据文件大小智能调整chunk大小
    let chunkSize;
    if (totalSize > 100 * 1024 * 1024) { // 大于100MB
        chunkSize = 5 * 1024 * 1024; // 5MB chunks
    } else if (totalSize > 50 * 1024 * 1024) { // 大于50MB
        chunkSize = 2 * 1024 * 1024; // 2MB chunks
    } else {
        chunkSize = 1024 * 1024; // 1MB chunks
    }
    
    console.log(`使用分块大小: ${formatFileSize(chunkSize)}`);
    
    const chunks = [];
    let downloadedSize = 0;
    const startTime = Date.now();
    
    const reader = response.body.getReader();
    
    try {
        while (true) {
            const { done, value } = await reader.read();
            
            if (done) break;
            
            chunks.push(value);
            downloadedSize += value.length;
            
            // 显示处理进度和速度
            const progress = Math.round((downloadedSize / totalSize) * 100);
            const elapsed = (Date.now() - startTime) / 1000;
            const speed = downloadedSize / elapsed;
            
            if (progress % 10 === 0) { // 每10%更新一次
                console.log(`文件数据处理进度: ${progress}%, 速度: ${formatFileSize(speed)}/s`);
                showStatusMessage(`正在处理文件数据... ${progress}% (${formatFileSize(speed)}/s)`, 'info');
            }
        }
        
        // 合并所有分块
        console.log('所有分块读取完成，正在合并...');
        showStatusMessage('正在合并文件数据...', 'info');
        
        const mergeStartTime = Date.now();
        const blob = new Blob(chunks);
        const mergeTime = Date.now() - mergeStartTime;
        
        console.log(`文件合并完成，耗时: ${mergeTime}ms`);
        showStatusMessage('文件数据处理完成', 'success');
        
        return blob;
        
    } catch (error) {
        console.error('分块读取失败:', error);
        throw new Error('文件数据处理失败: ' + error.message);
    } finally {
        reader.releaseLock();
    }
}

/**
 * 高效的Blob下载函数 - 性能优化版本
 * 优化点：
 * 1. 移除文件头验证，避免大文件的ArrayBuffer转换
 * 2. 使用更高效的DOM操作方式
 * 3. 智能的文件大小检测和提示
 * 4. 优化的URL对象管理
 */
async function downloadBlob(blob, filename) {
    try {
        // 验证参数
        if (!filename || !blob || blob.size === 0) {
            console.error('下载参数无效:', { filename, blobSize: blob?.size });
            throw new Error('下载参数无效');
        }
        
        // 显示下载开始消息
        showDownloadProgress('开始下载文件...', 'info');
        
        // 对于大文件，使用流式下载
        if (blob.size > 10 * 1024 * 1024) { // 大于10MB的文件
            console.log('检测到大文件，使用流式下载优化...');
            showDownloadProgress('正在优化大文件下载...', 'info');
        }
        
        // 使用更高效的下载方式
        if (navigator.msSaveBlob) {
            // IE 10+
            navigator.msSaveBlob(blob, filename);
        } else {
            // 现代浏览器 - 使用更高效的下载方式
            const url = window.URL.createObjectURL(blob);
            
            // 创建隐藏的下载链接
        const link = document.createElement('a');
            link.style.cssText = 'position:absolute;left:-9999px;top:-9999px;opacity:0;pointer-events:none;';
            link.href = url;
            link.download = filename;
            link.setAttribute('download', filename);
            
            // 添加到DOM并立即触发下载
        document.body.appendChild(link);
        link.click();
            
            // 立即移除元素并释放URL
        document.body.removeChild(link);
        
            // 延迟释放URL对象，确保下载开始
            setTimeout(() => {
                window.URL.revokeObjectURL(url);
            }, 200);
        }
        
        console.log(`文件 ${filename} 下载已触发，大小: ${formatFileSize(blob.size)}`);
        
    } catch (error) {
        console.error('下载Blob失败:', error);
        throw error;
    }
}

// 工具函数
function showStatusMessage(message, type = 'info', duration = 5000) {
    const statusElement = document.getElementById('statusMessage');
    
    // 如果消息包含"下载"关键词，延长显示时间
    if (message.includes('下载') && type === 'info') {
        duration = 8000;
    }
    
    statusElement.textContent = message;
    statusElement.className = `status-message ${type}`;
    statusElement.classList.remove('hidden');
    
    // 自动隐藏消息
    setTimeout(() => {
        statusElement.classList.add('hidden');
    }, duration);
}

// 显示下载进度消息（简化版，仅显示状态消息）
function showDownloadProgress(message, type = 'info') {
    showStatusMessage(message, type, 10000); // 下载相关消息显示更长时间
}

function showLoading(show) {
    const loadingElement = document.getElementById('loadingOverlay');
    if (show) {
        loadingElement.classList.remove('hidden');
    } else {
        loadingElement.classList.add('hidden');
    }
}

// 文件大小格式化函数
function formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// 兼容旧代码：如有其他文件引用到旧的下载实现，这里转发到上面的主下载流程
async function legacyDownload(cid, filename) {
    try {
        console.log('[兼容] legacyDownload 调用，转发到 downloadFile', { cid, filename });
        await downloadFile(cid, filename);
    } catch (error) {
        console.error('[兼容] legacyDownload 失败:', error);
        showStatusMessage('下载文件失败: ' + error.message, 'error');
    }
}

// 页面跳转功能
function redirectToLogin() {
    try {
        window.location.replace('/static/index.html');
        setTimeout(() => {
            if (!/index\.html$/i.test(window.location.pathname)) {
                window.location.href = '/static/index.html';
            }
        }, 100);
    } catch (_e) {
        window.location.href = '/static/index.html';
    }
}

function logout() {
    localStorage.removeItem('userInfo');
    localStorage.removeItem('isLoggedIn');
    redirectToLogin();
}

// 页面加载完成
document.addEventListener('DOMContentLoaded', function() {
    console.log('🚀 DOM加载完成');
});

// 点击模态框外部关闭
document.addEventListener('click', function(event) {
    const modal = document.getElementById('fileDetailModal');
    if (event.target === modal) {
        closeFileDetail();
    }
});

// 键盘事件处理
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        closeFileDetail();
    }
});

// 页面可见性变化检测（已移除下载进度条相关代码）
document.addEventListener('visibilitychange', function() {
    // 可以在这里添加其他页面可见性相关的逻辑
});
