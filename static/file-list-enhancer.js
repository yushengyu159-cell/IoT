// 文件列表显示增强脚本
// 提供更好的用户界面和交互

console.log('[FileList] 增强脚本加载...');

// 增强的文件列表显示函数
window.displayFileListEnhanced = function(files) {
    const fileListContainer = document.getElementById('fileList');
    if (!fileListContainer) {
        console.error('[FileList] 找不到 fileList 容器');
        return;
    }
    
    console.log('[FileList] 显示 ' + files.length + ' 个文件');
    
    if (!files || files.length === 0) {
        showEmptyFileList();
        return;
    }
    
    let html = '';
    files.forEach(function(file, index) {
        const fileName = file.Filename || 'unknown';
        const desc = file.Desc || '无描述';
        const uploadTime = file.UploadAt ? file.UploadAt.split('T')[0] : '';
        const cidShort = file.CID ? file.CID.substring(0, 16) + '...' : 'N/A';
        const fileSize = file.FileSize ? formatFileSize(file.FileSize) : '';
        
        html += '<div class=file-item style=padding: 16px
