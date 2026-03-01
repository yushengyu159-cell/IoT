#!/bin/bash

# 文件上传接口测试脚本
# 测试账号: 18027473816@163.com
# 接口地址: https://app.esgvisa.com/api/esg/upload-encrypted

API_URL="https://app.esgvisa.com/api/esg/upload-encrypted"
UPLOADER="18027473816@163.com"
TEST_FILE="${1:-test.pdf}"

echo "=========================================="
echo "文件上传接口测试"
echo "=========================================="
echo "接口地址: $API_URL"
echo "测试账号: $UPLOADER"
echo "测试文件: $TEST_FILE"
echo "=========================================="
echo ""

# 检查测试文件是否存在
if [ ! -f "$TEST_FILE" ]; then
    echo "❌ 错误: 测试文件 $TEST_FILE 不存在"
    echo "请先创建一个测试文件"
    echo ""
    echo "使用方法:"
    echo "  ./test_upload.sh [文件路径]"
    echo "  或"
    echo "  bash test_upload.sh [文件路径]"
    exit 1
fi

# 获取文件大小
FILE_SIZE=$(stat -f%z "$TEST_FILE" 2>/dev/null || stat -c%s "$TEST_FILE" 2>/dev/null)
FILE_SIZE_MB=$(echo "scale=2; $FILE_SIZE / 1024 / 1024" | bc)

echo "📄 文件信息:"
echo "  文件名: $(basename "$TEST_FILE")"
echo "  文件大小: $FILE_SIZE 字节 ($FILE_SIZE_MB MB)"
echo ""

echo "📤 开始上传文件..."
echo ""

# 执行上传
RESPONSE=$(curl -X POST "$API_URL" \
  -F "file=@$TEST_FILE" \
  -F "desc=自动化测试文件上传 - $(date '+%Y-%m-%d %H:%M:%S')" \
  -F "uploader=$UPLOADER" \
  -w "\nHTTP_CODE:%{http_code}" \
  -s)

# 提取HTTP状态码
HTTP_CODE=$(echo "$RESPONSE" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed 's/HTTP_CODE:[0-9]*$//')

echo "HTTP状态码: $HTTP_CODE"
echo ""
echo "响应内容:"

# 尝试使用jq格式化JSON，如果没有jq则直接输出
if command -v jq &> /dev/null; then
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
    echo "$BODY"
    echo ""
    echo "💡 提示: 安装jq可以更好地格式化JSON输出 (apt-get install jq 或 brew install jq)"
fi

echo ""

# 检查上传结果
if [ "$HTTP_CODE" = "200" ]; then
    # 尝试提取code字段
    if command -v jq &> /dev/null; then
        CODE=$(echo "$BODY" | jq -r '.code' 2>/dev/null)
        MESSAGE=$(echo "$BODY" | jq -r '.message' 2>/dev/null)
    else
        # 简单的grep提取（如果jq不可用）
        CODE=$(echo "$BODY" | grep -o '"code":[0-9]*' | cut -d: -f2)
        MESSAGE=$(echo "$BODY" | grep -o '"message":"[^"]*"' | cut -d'"' -f4)
    fi
    
    if [ "$CODE" = "200" ]; then
        echo "✅ 上传成功!"
        echo ""
        
        # 提取关键信息（如果jq可用）
        if command -v jq &> /dev/null; then
            CID=$(echo "$BODY" | jq -r '.data.meta.cids[0]' 2>/dev/null)
            KEY=$(echo "$BODY" | jq -r '.data.key' 2>/dev/null)
            FILE_NAME=$(echo "$BODY" | jq -r '.data.meta.fileName' 2>/dev/null)
            FILE_SIZE_RESP=$(echo "$BODY" | jq -r '.data.meta.fileSize' 2>/dev/null)
            CHUNK_COUNT=$(echo "$BODY" | jq -r '.data.chunkCount' 2>/dev/null)
            TOTAL_TIME=$(echo "$BODY" | jq -r '.data.timeStats.totalTime' 2>/dev/null)
            
            echo "📋 文件信息:"
            echo "  文件名: $FILE_NAME"
            echo "  文件大小: $FILE_SIZE_RESP 字节"
            echo "  分片数量: $CHUNK_COUNT"
            echo "  上传耗时: $TOTAL_TIME"
            echo ""
            echo "🔑 关键数据:"
            echo "  CID: $CID"
            echo "  加密密钥: $KEY"
            echo ""
            echo "⚠️  请妥善保管加密密钥，用于后续文件下载"
        else
            echo "📋 文件信息已在上方响应中显示"
            echo "⚠️  请妥善保管返回的加密密钥，用于后续文件下载"
        fi
        
        # 保存结果到文件
        TIMESTAMP=$(date +%Y%m%d_%H%M%S)
        RESULT_FILE="upload_result_${TIMESTAMP}.json"
        echo "$BODY" > "$RESULT_FILE"
        echo ""
        echo "✅ 响应已保存到文件: $RESULT_FILE"
    else
        echo "❌ 上传失败:"
        echo "  错误码: $CODE"
        echo "  错误消息: $MESSAGE"
        exit 1
    fi
else
    echo "❌ HTTP错误: $HTTP_CODE"
    echo "响应内容: $BODY"
    exit 1
fi

echo ""
echo "=========================================="
echo "测试完成"
echo "=========================================="




