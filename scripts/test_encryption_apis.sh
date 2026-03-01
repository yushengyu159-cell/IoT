#!/bin/bash

# 分片加密存储接口测试脚本
# 测试ESG和IPFS相关的加密存储功能

set -e

BASE_URL="http://localhost:8199"
API_BASE="$BASE_URL/api"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查服务状态
check_service() {
    log_info "检查服务状态..."
    if curl -s "$API_BASE" > /dev/null; then
        log_success "服务运行正常"
        return 0
    else
        log_error "服务未运行，请先启动服务"
        return 1
    fi
}

# 测试基础API
test_basic_api() {
    log_info "测试基础API..."
    response=$(curl -s "$API_BASE")
    if echo "$response" | grep -q "Hello Fabric SDK"; then
        log_success "基础API正常"
    else
        log_error "基础API异常"
    fi
}

# 测试ESG文件列表
test_esg_list() {
    log_info "测试ESG文件列表..."
    response=$(curl -s "$API_BASE/esg/list")
    echo "ESG文件列表响应: $response"
}

# 测试ESG文件查询
test_esg_query() {
    log_info "测试ESG文件查询..."
    response=$(curl -s "$API_BASE/esg/query?cid=test_cid")
    echo "ESG文件查询响应: $response"
}

# 测试ESG文件上传（加密）
test_esg_upload_encrypted() {
    log_info "测试ESG加密文件上传..."
    
    # 创建测试文件
    echo "这是一个测试文件内容" > /tmp/test_file.txt
    
    response=$(curl -s -X POST \
        -F "file=@/tmp/test_file.txt" \
        -F "description=测试加密上传" \
        -F "uploader=test_user" \
        "$API_BASE/esg/upload-encrypted")
    
    echo "ESG加密上传响应: $response"
    
    # 清理测试文件
    rm -f /tmp/test_file.txt
}

# 测试ESG文件上传（普通）
test_esg_upload() {
    log_info "测试ESG普通文件上传..."
    
    # 创建测试文件
    echo "这是一个普通测试文件" > /tmp/test_file_normal.txt
    
    response=$(curl -s -X POST \
        -F "file=@/tmp/test_file_normal.txt" \
        -F "description=测试普通上传" \
        -F "uploader=test_user" \
        "$API_BASE/esg/upload")
    
    echo "ESG普通上传响应: $response"
    
    # 清理测试文件
    rm -f /tmp/test_file_normal.txt
}

# 测试IPFS上传（加密）
test_ipfs_upload_encrypted() {
    log_info "测试IPFS加密上传..."
    
    # 创建测试文件
    echo "IPFS加密测试文件内容" > /tmp/ipfs_test_encrypted.txt
    
    response=$(curl -s -X POST \
        -F "file=@/tmp/ipfs_test_encrypted.txt" \
        -F "description=IPFS加密测试" \
        "$API_BASE/ipfs/upload-encrypted")
    
    echo "IPFS加密上传响应: $response"
    
    # 清理测试文件
    rm -f /tmp/ipfs_test_encrypted.txt
}

# 测试IPFS上传（普通）
test_ipfs_upload() {
    log_info "测试IPFS普通上传..."
    
    # 创建测试文件
    echo "IPFS普通测试文件内容" > /tmp/ipfs_test_normal.txt
    
    response=$(curl -s -X POST \
        -F "file=@/tmp/ipfs_test_normal.txt" \
        -F "description=IPFS普通测试" \
        "$API_BASE/ipfs/upload")
    
    echo "IPFS普通上传响应: $response"
    
    # 清理测试文件
    rm -f /tmp/ipfs_test_normal.txt
}

# 测试DID注册
test_did_register() {
    log_info "测试DID注册..."
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "name": "测试用户",
            "phone": "13800138000",
            "email": "test@example.com",
            "password": "test123456",
            "role": "user",
            "age": 25
        }' \
        "$API_BASE/did/register")
    
    echo "DID注册响应: $response"
}

# 测试DID验证
test_did_verify() {
    log_info "测试DID验证..."
    
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "did": "did:example:123456789",
            "signature": "test_signature"
        }' \
        "$API_BASE/did/verify")
    
    echo "DID验证响应: $response"
}

# 测试链码功能
test_chaincode() {
    log_info "测试链码功能..."
    
    # 测试写入记录
    log_info "测试写入记录..."
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "key": "test_key_001",
            "value": "test_value_001"
        }' \
        "$API_BASE/chaincode/write-record")
    echo "写入记录响应: $response"
    
    # 测试读取记录
    log_info "测试读取记录..."
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "key": "test_key_001"
        }' \
        "$API_BASE/chaincode/read-record")
    echo "读取记录响应: $response"
    
    # 测试创建资产
    log_info "测试创建资产..."
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "id": "asset_001",
            "owner": "test_owner",
            "value": 100
        }' \
        "$API_BASE/chaincode/create-asset")
    echo "创建资产响应: $response"
}

# 主测试流程
main() {
    echo "🔧 开始测试分片加密存储接口..."
    echo "=================================="
    
    # 检查服务状态
    if ! check_service; then
        exit 1
    fi
    
    # 测试基础功能
    test_basic_api
    
    echo ""
    log_info "开始测试DID功能..."
    test_did_register
    test_did_verify
    
    echo ""
    log_info "开始测试链码功能..."
    test_chaincode
    
    echo ""
    log_info "开始测试ESG功能..."
    test_esg_list
    test_esg_query
    test_esg_upload
    test_esg_upload_encrypted
    
    echo ""
    log_info "开始测试IPFS功能..."
    test_ipfs_upload
    test_ipfs_upload_encrypted
    
    echo ""
    echo "=================================="
    log_success "所有接口测试完成！"
    echo "=================================="
    
    echo ""
    log_info "手动测试命令参考:"
    echo "1. 查看所有ESG文件: curl -s '$API_BASE/esg/list'"
    echo "2. 查询特定文件: curl -s '$API_BASE/esg/query?cid=your_cid'"
    echo "3. 上传加密文件: curl -X POST -F 'file=@your_file.txt' '$API_BASE/esg/upload-encrypted'"
    echo "4. 下载加密文件: curl -X POST -d '{\"cid\":\"your_cid\"}' '$API_BASE/esg/download-encrypted'"
    echo "5. IPFS上传: curl -X POST -F 'file=@your_file.txt' '$API_BASE/ipfs/upload-encrypted'"
    echo "6. 链码测试: curl -s '$API_BASE/chaincode/test'"
}

# 运行主函数
main "$@" 