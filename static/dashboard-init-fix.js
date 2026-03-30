// Dashboard 初始化修复脚本
console.log('[Dashboard] 初始化修复脚本 v2');

window.loadFileListFixed = async function() {
    console.log('[Dashboard] 加载文件列表...');
    
    try {
        const userInfoStr = localStorage.getItem("userInfo");
        if (!userInfoStr) {
            console.error("未找到用户信息");
            window.location.href = "/static/index.html#login";
            return;
        }
        
        const userInfo = JSON.parse(userInfoStr);
        const userEmail = userInfo.Email || userInfo.email;
        
        if (!userEmail) {
            alert("用户信息不完整，请重新登录");
            window.location.href = "/static/index.html#login";
            return;
        }
        
        console.log("用户邮箱:", userEmail);
        
        const url = "/api/esg/list?userEmail=" + encodeURIComponent(userEmail);
        const response = await fetch(url);
        
        if (!response.ok) {
            throw new Error("HTTP " + response.status);
        }
        
        const result = await response.json();
        
        if (result.code === 200) {
            const files = result.data || [];
            console.log("成功获取 " + files.length + " 个文件");
            displayFileListFixed(files);
        } else {
            alert("加载文件列表失败: " + result.message);
        }
    } catch (error) {
        console.error("加载文件列表异常:", error);
        alert("加载文件列表失败: " + error.message);
    }
};

window.displayFileListFixed = function(files) {
    const fileListContainer = document.getElementById("fileList");
    if (!fileListContainer) {
        console.error("找不到 fileList 容器");
        return;
    }
    
    let html = "";
    files.forEach(function(file) {
        const uploadTime = file.UploadAt || new Date().toISOString().split("T")[0];
        const fileName = file.Filename || "unknown";
        const desc = file.Desc || "无描述";
        const cidShort = file.CID ? file.CID.substring(0, 16) + "..." : "N/A";
        
        html += "<div class=\"file-item\" style=\"padding: 16px; background: #f8f9fa; margin-bottom: 12px; border-radius: 8px; border-left: 4px solid #2ECC71;\">";
        html += "  <div style=\"display: flex; justify-content: space-between; align-items: start;\">";
        html += "    <div style=\"flex: 1;\">";
        html += "      <div style=\"font-weight: 600; font-size: 16px; margin-bottom: 8px;\">";
        html += "        <i class=\"fas fa-file\" style=\"color: #2ECC71; margin-right: 8px;\"></i>";
        html += fileName;
        html += "      </div>";
        html += "      <div style=\"font-size: 14px; color: #666; margin-bottom: 4px;\">" + desc + "</div>";
        html += "      <div style=\"font-size: 12px; color: #999;\">";
        html += "        <i class=\"fas fa-clock\"></i> " + uploadTime;
        html += "        <span style=\"margin-left: 16px;\"><i class=\"fas fa-fingerprint\"></i> " + cidShort + "</span>";
        html += "      </div>";
        html += "    </div>";
        html += "    <div style=\"display: flex; gap: 8px;\">";
        html += "      <button class=\"btn btn-sm btn-primary\" onclick=\"downloadFile('" + file.CID + "', '" + fileName + "')\">";
        html += "        <i class=\"fas fa-download\"></i> 下载";
        html += "      </button>";
        html += "      <button class=\"btn btn-sm btn-secondary\" onclick=\"viewFileDetails('" + file.CID + "')\">";
        html += "        <i class=\"fas fa-info-circle\"></i> 详情";
        html += "      </button>";
        html += "    </div>";
        html += "  </div>";
        html += "</div>";
    });
    
    fileListContainer.innerHTML = html;
};

document.addEventListener("DOMContentLoaded", function() {
    setTimeout(function() {
        console.log("[Dashboard] 延迟加载文件列表...");
        loadFileListFixed();
    }, 1000);
});

console.log("[Dashboard] 修复脚本已加载");
