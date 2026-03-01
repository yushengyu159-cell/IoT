#!/bin/bash

# Fabric SDK 集成测试脚本
# 测试WriteRecord和ReadRecord功能的SDK集成

set -e

echo "🚀 开始Fabric SDK集成测试..."
echo "=================================="

# 配置
API_BASE="http://localhost:8199/api/chaincode"
SLEEP_TIME=2

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

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    log_info "测试: $description"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$API_BASE$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$API_BASE$endpoint" -d "$data")
    fi
    
    # 分离响应体和状态码
    body=$(echo "$response" | head -n -1)
    status_code=$(echo "$response" | tail -n 1)
    
    if [ "$status_code" = "200" ]; then
        log_success "$description - 成功 (HTTP $status_code)"
        echo "响应: $body" | jq . 2>/dev/null || echo "响应: $body"
    else
        log_error "$description - 失败 (HTTP $status_code)"
        echo "响应: $body"
        return 1
    fi
    
    echo ""
    sleep $SLEEP_TIME
}

# 检查服务是否运行
check_service() {
    log_info "检查SDK服务状态..."
    
    if curl -s "http://localhost:8199/api" > /dev/null; then
        log_success "SDK服务正在运行"
    else
        log_error "SDK服务未运行，请先启动服务"
        exit 1
    fi
}

# 主测试流程
main() {
    echo "🔧 开始Fabric SDK集成测试..."
    echo "=================================="
    
    # 检查服务状态
    check_service
    
    # 1. 初始化链码服务
    log_info "步骤 1: 初始化链码服务"
    test_api "GET" "/init" "" "初始化链码服务"
    
    # 2. 测试WriteRecord功能
    log_info "步骤 2: 测试WriteRecord功能"
    test_api "POST" "/write-record" "key=test_key_001&value=test_value_001" "写入记录 test_key_001"
    test_api "POST" "/write-record" "key=test_key_002&value={\"name\":\"test\",\"value\":123}" "写入JSON记录 test_key_002"
    test_api "POST" "/write-record" "key=test_key_003&value=hello_world" "写入简单记录 test_key_003"
    
    # 3. 测试ReadRecord功能
    log_info "步骤 3: 测试ReadRecord功能"
    test_api "GET" "/read-record?key=test_key_001" "" "读取记录 test_key_001"
    test_api "GET" "/read-record?key=test_key_002" "" "读取JSON记录 test_key_002"
    test_api "GET" "/read-record?key=test_key_003" "" "读取简单记录 test_key_003"
    
    # 4. 测试资产管理功能
    log_info "步骤 4: 测试资产管理功能"
    test_api "POST" "/create-asset" "assetID=asset_001&color=red&size=10&owner=alice&appraisedValue=1000" "创建资产 asset_001"
    test_api "GET" "/read-asset?assetID=asset_001" "" "读取资产 asset_001"
    test_api "POST" "/transfer-asset" "assetID=asset_001&newOwner=bob" "转移资产 asset_001 给 bob"
    test_api "GET" "/read-asset?assetID=asset_001" "" "再次读取资产 asset_001"
    
    # 5. 测试批量操作
    log_info "步骤 5: 测试批量操作"
    test_api "POST" "/create-asset" "assetID=asset_002&color=blue&size=20&owner=charlie&appraisedValue=2000" "创建资产 asset_002"
    test_api "POST" "/create-asset" "assetID=asset_003&color=green&size=15&owner=david&appraisedValue=1500" "创建资产 asset_003"
    test_api "GET" "/get-all-assets" "" "获取所有资产"
    
    # 6. 测试错误处理
    log_info "步骤 6: 测试错误处理"
    test_api "GET" "/read-record?key=nonexistent_key" "" "读取不存在的记录"
    test_api "GET" "/read-asset?assetID=nonexistent_asset" "" "读取不存在的资产"
    
    # 7. 完整功能测试
    log_info "步骤 7: 完整功能测试"
    test_api "GET" "/test" "" "执行完整功能测试"
    
    echo "=================================="
    log_success "所有测试完成！"
    echo "=================================="
    
    # 显示API文档
    echo ""
    log_info "API接口文档:"
    echo "=================================="
    echo "初始化服务: GET/POST $API_BASE/init"
    echo "写入记录: POST $API_BASE/write-record (key, value)"
    echo "读取记录: GET $API_BASE/read-record?key=<key>"
    echo "创建资产: POST $API_BASE/create-asset (assetID, color, size, owner, appraisedValue)"
    echo "读取资产: GET $API_BASE/read-asset?assetID=<assetID>"
    echo "转移资产: POST $API_BASE/transfer-asset (assetID, newOwner)"
    echo "获取所有资产: GET $API_BASE/get-all-assets"
    echo "测试功能: GET $API_BASE/test"
    echo "关闭服务: GET/POST $API_BASE/close"
    echo "=================================="
}

# 错误处理
trap 'log_error "测试被中断"; exit 1' INT TERM

# 运行主测试
main "$@" 