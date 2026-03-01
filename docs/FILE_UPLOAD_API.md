# 文件上传接口 API 文档

## 概述

本文档描述了供外部系统调用的文件上传接口。该接口支持文件加密、分片上传、IPFS存储和区块链存证功能。

**系统域名**: `app.esgvisa.com`  
**基础URL**: `https://app.esgvisa.com` (生产环境) / `http://app.esgvisa.com:8199` (开发环境)  
**接口路径**: `/api/esg/upload-encrypted`  
**请求方法**: `POST`  
**内容类型**: `multipart/form-data`

> **注意**: 生产环境建议使用HTTPS协议，确保数据传输安全。开发环境可使用HTTP协议。

---

## 接口说明

### 1. 加密分片上传接口（推荐）

**接口地址**: `POST /api/esg/upload-encrypted`

**功能特点**:
- ✅ 文件自动加密存储
- ✅ 自动分片上传到IPFS（平均分为3片）
- ✅ 区块链存证（记录交易哈希）
- ✅ 数据库元数据存储
- ✅ 支持跨域访问（CORS已配置）
- ✅ 返回完整的文件信息和访问链接

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `file` | File | 是 | 要上传的文件（支持任意格式） |
| `desc` | String | 否 | 文件描述信息 |
| `uploader` | String | 否 | 上传者身份标识（建议格式：`did:example:xxx` 或用户邮箱） |

**请求示例**:

#### cURL 示例
```bash
# 生产环境（HTTPS）
curl -X POST "https://app.esgvisa.com/api/esg/upload-encrypted" \
  -F "file=@/path/to/your/file.pdf" \
  -F "desc=这是文件描述" \
  -F "uploader=external-system@example.com"

# 开发环境（HTTP）
curl -X POST "http://app.esgvisa.com:8199/api/esg/upload-encrypted" \
  -F "file=@/path/to/your/file.pdf" \
  -F "desc=这是文件描述" \
  -F "uploader=external-system@example.com"
```

#### JavaScript (Fetch API) 示例
```javascript
// 配置API基础URL（根据环境选择）
const API_BASE_URL = 'https://app.esgvisa.com'; // 生产环境
// const API_BASE_URL = 'http://app.esgvisa.com:8199'; // 开发环境

const formData = new FormData();
formData.append('file', fileInput.files[0]); // fileInput 是文件输入元素
formData.append('desc', '文件描述信息');
formData.append('uploader', 'external-system@example.com');

fetch(`${API_BASE_URL}/api/esg/upload-encrypted`, {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(data => {
  console.log('上传成功:', data);
  // 保存返回的key和cid，用于后续下载
  if (data.code === 200 && data.data) {
    localStorage.setItem('file_key', data.data.key);
    localStorage.setItem('file_cid', data.data.meta.cids[0]);
  }
})
.catch(error => {
  console.error('上传失败:', error);
});
```

#### Python 示例
```python
import requests

# 配置API基础URL（根据环境选择）
API_BASE_URL = 'https://app.esgvisa.com'  # 生产环境
# API_BASE_URL = 'http://app.esgvisa.com:8199'  # 开发环境

url = f"{API_BASE_URL}/api/esg/upload-encrypted"
files = {
    'file': ('document.pdf', open('/path/to/document.pdf', 'rb'), 'application/pdf')
}
data = {
    'desc': '文件描述信息',
    'uploader': 'external-system@example.com'
}

try:
    response = requests.post(url, files=files, data=data, timeout=300)  # 5分钟超时
    response.raise_for_status()  # 检查HTTP错误
    result = response.json()
    
    if result.get('code') == 200:
        print('上传成功!')
        print(f"文件CID: {result['data']['meta']['cids'][0]}")
        print(f"加密密钥: {result['data']['key']}")
        # 保存密钥用于后续下载
    else:
        print(f"上传失败: {result.get('message')}")
except requests.exceptions.RequestException as e:
    print(f"请求错误: {e}")
```

#### Java 示例
```java
import java.io.File;
import java.io.IOException;
import okhttp3.*;

public class FileUpload {
    // 配置API基础URL（根据环境选择）
    private static final String API_BASE_URL = "https://app.esgvisa.com"; // 生产环境
    // private static final String API_BASE_URL = "http://app.esgvisa.com:8199"; // 开发环境
    
    public static void main(String[] args) throws IOException {
        OkHttpClient client = new OkHttpClient.Builder()
            .connectTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
            .writeTimeout(300, java.util.concurrent.TimeUnit.SECONDS)
            .readTimeout(300, java.util.concurrent.TimeUnit.SECONDS)
            .build();
        
        RequestBody requestBody = new MultipartBody.Builder()
            .setType(MultipartBody.FORM)
            .addFormDataPart("file", "document.pdf",
                RequestBody.create(new File("/path/to/document.pdf"), 
                    MediaType.parse("application/pdf")))
            .addFormDataPart("desc", "文件描述信息")
            .addFormDataPart("uploader", "external-system@example.com")
            .build();
        
        Request request = new Request.Builder()
            .url(API_BASE_URL + "/api/esg/upload-encrypted")
            .post(requestBody)
            .build();
        
        try (Response response = client.newCall(request).execute()) {
            if (response.isSuccessful()) {
                String responseBody = response.body().string();
                System.out.println("上传成功: " + responseBody);
                // 解析JSON并保存key和cid
            } else {
                System.err.println("上传失败: " + response.code() + " " + response.message());
            }
        }
    }
}
```

**响应格式**:

#### 成功响应 (200)
```json
{
  "code": 200,
  "message": "加密分片上传并链上存证成功",
  "data": {
    "meta": {
      "fileName": "document.pdf",
      "desc": "文件描述信息",
      "uploader": "external-system@example.com",
      "fileSize": 1024000,
      "chunkSize": 341334,
      "cids": [
        "QmXXX...",
        "QmYYY...",
        "QmZZZ..."
      ],
      "iv": "base64_encoded_iv",
      "uploadAt": "2025-01-XX 12:34:56"
    },
    "key": "base64_encoded_encryption_key",
    "onChain": true,
    "cipherSample": "base64_encoded_cipher_sample",
    "chunkCount": 3,
    "chunkUrls": [
      "http://app.esgvisa.com:8081/ipfs/QmXXX...",
      "http://app.esgvisa.com:8081/ipfs/QmYYY...",
      "http://app.esgvisa.com:8081/ipfs/QmZZZ..."
    ],
    "fabricResult": {
      "txID": "mysql_storage_20250101123456",
      "status": "success",
      "message": "文件已存储到MySQL数据库",
      "assetID": "QmXXX..."
    },
    "timeStats": {
      "totalTime": "1.234s",
      "totalTimeMs": 1234,
      "startTime": "2025-01-XX 12:34:55.000",
      "endTime": "2025-01-XX 12:34:56.234"
    }
  }
}
```

#### 错误响应 (400)
```json
{
  "code": 400,
  "message": "未上传文件",
  "data": null
}
```

#### 错误响应 (500)
```json
{
  "code": 500,
  "message": "加密分片上传链上存证失败: 具体错误信息",
  "data": null
}
```

**响应字段说明**:

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | Integer | 响应状态码（200表示成功） |
| `message` | String | 响应消息 |
| `data.meta` | Object | 文件元数据 |
| `data.meta.fileName` | String | 原始文件名 |
| `data.meta.fileSize` | Integer | 文件大小（字节） |
| `data.meta.chunkCount` | Integer | 分片数量（通常为3） |
| `data.meta.cids` | Array | IPFS内容标识符数组 |
| `data.key` | String | Base64编码的加密密钥（用于下载解密） |
| `data.onChain` | Boolean | 是否已上链存证 |
| `data.chunkUrls` | Array | IPFS访问链接数组 |
| `data.fabricResult` | Object | 区块链存证结果 |
| `data.timeStats` | Object | 上传耗时统计 |

---

### 2. 普通上传接口（备选）

**接口地址**: `POST /api/esg/upload`

**功能特点**:
- ✅ 直接上传到IPFS
- ✅ 区块链存证
- ✅ 数据库记录
- ⚠️ 不进行加密和分片处理

**请求参数**:

| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| `file` | File | 是 | 要上传的文件 |
| `desc` | String | 否 | 文件描述 |
| `uploader` | String | 否 | 上传者身份 |

**响应格式**:
```json
{
  "code": 200,
  "message": "文件已成功存储到区块链",
  "data": {
    "cid": "QmXXX...",
    "onChain": true
  }
}
```

---

## 文件查询接口

### 查询文件详情

**接口地址**: `GET /api/esg/query?cid={cid}`

**请求示例**:
```bash
# 生产环境
curl "https://app.esgvisa.com/api/esg/query?cid=QmXXX..."

# 开发环境
curl "http://app.esgvisa.com:8199/api/esg/query?cid=QmXXX..."
```

**响应格式**:
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "ID": 1,
    "CID": "QmXXX...",
    "Filename": "document.pdf",
    "Desc": "文件描述",
    "Uploader": "external-system@example.com",
    "UploadAt": "2025-01-XX 12:34:56",
    "Txid": "mysql_storage_20250101123456",
    "ChunkCount": 3,
    "FileSize": 1024000,
    "uploaderDid": "did:example:xxx",
    "uploaderName": "上传者姓名",
    "uploaderEmail": "external-system@example.com",
    "uploaderRole": "owner"
  }
}
```

---

## 文件下载接口

### 下载加密文件

**接口地址**: `POST /api/esg/download-encrypted`

**请求参数** (JSON):
```json
{
  "cid": "QmXXX...",
  "key": "base64_encoded_encryption_key"
}
```

**请求示例**:
```bash
# 生产环境
curl -X POST "https://app.esgvisa.com/api/esg/download-encrypted" \
  -H "Content-Type: application/json" \
  -d '{
    "cid": "QmXXX...",
    "key": "base64_encoded_key"
  }' \
  --output downloaded_file.pdf

# 开发环境
curl -X POST "http://app.esgvisa.com:8199/api/esg/download-encrypted" \
  -H "Content-Type: application/json" \
  -d '{
    "cid": "QmXXX...",
    "key": "base64_encoded_key"
  }' \
  --output downloaded_file.pdf
```

**响应**: 直接返回文件二进制流

---

## 跨域访问配置

系统已配置CORS中间件，支持跨域访问。外部系统可以直接调用接口，无需额外配置。

**CORS配置**:
- 允许所有来源（`Access-Control-Allow-Origin: *`）
- 支持所有HTTP方法（GET, POST, PUT, DELETE, OPTIONS等）
- 支持所有请求头
- 支持携带凭证（Credentials）

**注意事项**:
- 生产环境建议配置具体的允许来源，提高安全性
- 如需限制特定域名访问，请联系系统管理员配置

---

## 安全建议

### 1. 生产环境部署

在生产环境中（`https://app.esgvisa.com`），建议：

1. **使用HTTPS协议**:
   - ✅ 系统已配置HTTPS支持
   - ✅ 所有数据传输均加密
   - ⚠️ 请勿在生产环境使用HTTP协议

2. **添加API认证机制**:
   - 使用API Key认证（推荐）
   - 或使用JWT Token认证
   - 限制访问来源IP白名单

3. **文件大小限制**:
   - 当前无限制，建议根据需求设置最大文件大小（如100MB）
   - 大文件上传建议实现进度条和断点续传

4. **访问频率限制**:
   - 添加Rate Limiting防止滥用（如每分钟最多10次上传）
   - 实现IP级别的访问控制

5. **日志和监控**:
   - 记录所有上传操作日志
   - 监控异常访问行为
   - 设置告警机制

### 2. 密钥管理

- **加密密钥**: 上传成功后返回的`key`字段是Base64编码的加密密钥，请妥善保管
- **密钥用途**: 下载文件时必须提供正确的密钥才能解密
- **密钥存储**: 建议将密钥存储在安全的密钥管理系统中

---

## 错误处理

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 参数错误 | 检查请求参数是否正确 |
| 500 | 服务器错误 | 检查服务器日志，联系管理员 |

### 错误处理示例

```javascript
const API_BASE_URL = 'https://app.esgvisa.com'; // 生产环境
// const API_BASE_URL = 'http://app.esgvisa.com:8199'; // 开发环境

fetch(`${API_BASE_URL}/api/esg/upload-encrypted`, {
  method: 'POST',
  body: formData
})
.then(async response => {
  // 检查HTTP状态码
  if (!response.ok) {
    throw new Error(`HTTP错误: ${response.status} ${response.statusText}`);
  }
  
  const data = await response.json();
  if (data.code === 200) {
    console.log('上传成功:', data);
    // 保存返回的key和cid，用于后续下载
    localStorage.setItem('file_key', data.data.key);
    localStorage.setItem('file_cid', data.data.meta.cids[0]);
    localStorage.setItem('file_name', data.data.meta.fileName);
  } else {
    console.error('上传失败:', data.message);
    // 根据错误码处理不同错误
    switch(data.code) {
      case 400:
        alert('参数错误: ' + data.message);
        break;
      case 500:
        alert('服务器错误: ' + data.message);
        break;
      default:
        alert('上传失败: ' + data.message);
    }
  }
})
.catch(error => {
  console.error('网络错误:', error);
  alert('网络连接失败，请检查网络后重试');
});
```

---

## 性能说明

### 上传流程

1. **文件读取**: 读取完整文件到内存
2. **文件加密**: 使用AES加密算法加密文件
3. **文件分片**: 将加密后的文件平均分为3片
4. **IPFS上传**: 依次上传每个分片到IPFS网络
5. **数据库存储**: 保存文件元数据到MySQL数据库
6. **区块链存证**: 记录交易哈希（可选）

### 性能指标

- **小文件 (< 10MB)**: 通常 < 5秒
- **中等文件 (10-100MB)**: 通常 5-30秒
- **大文件 (> 100MB)**: 根据网络情况，可能需要更长时间

### 优化建议

1. **分片上传**: 系统已自动分片，提高上传可靠性
2. **并发上传**: 可以同时上传多个文件（注意服务器负载）
3. **断点续传**: 当前不支持，建议实现客户端重试机制

---

## 测试示例

### 完整的上传和下载流程

```bash
#!/bin/bash

# 配置API基础URL（根据环境选择）
API_BASE_URL="https://app.esgvisa.com"  # 生产环境
# API_BASE_URL="http://app.esgvisa.com:8199"  # 开发环境

echo "=== 文件上传测试 ==="

# 1. 上传文件
echo "正在上传文件..."
UPLOAD_RESPONSE=$(curl -X POST "${API_BASE_URL}/api/esg/upload-encrypted" \
  -F "file=@test.pdf" \
  -F "desc=测试文件" \
  -F "uploader=test@example.com" \
  -s)

# 检查上传是否成功
if [ $? -ne 0 ]; then
  echo "❌ 上传失败: 网络错误"
  exit 1
fi

# 2. 提取CID和密钥（需要安装jq工具: apt-get install jq 或 brew install jq）
CID=$(echo $UPLOAD_RESPONSE | jq -r '.data.meta.cids[0]')
KEY=$(echo $UPLOAD_RESPONSE | jq -r '.data.key')
CODE=$(echo $UPLOAD_RESPONSE | jq -r '.code')

if [ "$CODE" != "200" ]; then
  echo "❌ 上传失败: $(echo $UPLOAD_RESPONSE | jq -r '.message')"
  exit 1
fi

echo "✅ 上传成功!"
echo "文件CID: $CID"
echo "加密密钥: $KEY"
echo ""

# 3. 查询文件信息
echo "正在查询文件信息..."
QUERY_RESPONSE=$(curl "${API_BASE_URL}/api/esg/query?cid=$CID" -s)
echo "文件信息:"
echo $QUERY_RESPONSE | jq '.'
echo ""

# 4. 下载文件（需要密钥）
echo "正在下载文件..."
curl -X POST "${API_BASE_URL}/api/esg/download-encrypted" \
  -H "Content-Type: application/json" \
  -d "{\"cid\":\"$CID\",\"key\":\"$KEY\"}" \
  --output downloaded_test.pdf

if [ $? -eq 0 ]; then
  echo "✅ 下载成功: downloaded_test.pdf"
else
  echo "❌ 下载失败"
fi
```

---

## 环境说明

### 生产环境
- **域名**: `app.esgvisa.com`
- **协议**: HTTPS
- **端口**: 443 (默认) 或 8199
- **API地址**: `https://app.esgvisa.com/api/esg/upload-encrypted`

### 开发环境
- **域名**: `app.esgvisa.com`
- **协议**: HTTP
- **端口**: 8199
- **API地址**: `http://app.esgvisa.com:8199/api/esg/upload-encrypted`

### 测试页面
访问测试页面进行接口测试：
- **生产环境**: `https://app.esgvisa.com/static/external-upload-test.html`
- **开发环境**: `http://app.esgvisa.com:8199/static/external-upload-test.html`

---

## 常见问题 (FAQ)

### Q1: 上传文件时出现网络错误？
**A**: 请检查：
1. 网络连接是否正常
2. 域名是否正确（`app.esgvisa.com`）
3. 是否使用了正确的协议（生产环境使用HTTPS）
4. 防火墙是否允许访问

### Q2: 如何获取上传后的文件？
**A**: 上传成功后会返回：
- `cid`: IPFS内容标识符，用于访问文件
- `key`: 加密密钥，下载文件时必须提供
- 使用 `/api/esg/download-encrypted` 接口下载文件

### Q3: 支持哪些文件格式？
**A**: 支持所有文件格式，包括但不限于：
- 文档：PDF, DOC, DOCX, XLS, XLSX, PPT, PPTX
- 图片：JPG, PNG, GIF, BMP, SVG
- 视频：MP4, AVI, MOV, WMV
- 其他：ZIP, RAR, TXT等

### Q4: 文件大小有限制吗？
**A**: 当前无硬性限制，但建议：
- 小文件（< 10MB）：上传速度快
- 中等文件（10-100MB）：可能需要较长时间
- 大文件（> 100MB）：建议分块上传或使用专用上传工具

### Q5: 上传失败如何处理？
**A**: 
1. 检查返回的错误码和错误信息
2. 确认文件格式和大小是否正常
3. 重试上传（系统支持重试）
4. 联系技术支持

---

## 联系与支持

**系统域名**: app.esgvisa.com  
**技术支持**: 如有问题或需要技术支持，请联系系统管理员  
**API文档**: 本文档持续更新，请关注最新版本

---

**文档版本**: v1.1  
**最后更新**: 2025-01-XX  
**维护者**: ESG Visa System Administrator  
**系统域名**: app.esgvisa.com

