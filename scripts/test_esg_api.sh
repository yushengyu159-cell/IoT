#!/bin/bash

# ESG API测试脚本
BASE_URL="http://localhost:8199"

echo "=== ESG API 测试脚本 ==="
echo "基础URL: $BASE_URL"
echo ""

# 1. 测试文件上传
echo "1. 测试文件上传..."
# 读取测试.pdf文件内容
PDF_CONTENT=$(cat ../测试.pdf | base64 -w 0)
UPLOAD_RESPONSE=$(curl -s -X POST "$BASE_URL/api/esg/upload" \
  -H "Content-Type: application/json" \
  -d "{
    \"fileName\": \"测试.pdf\",
    \"fileContent\": \"$PDF_CONTENT\",
    \"fileType\": \"pdf\",
    \"description\": \"测试PDF文档\"
  }")

echo "上传响应:"
echo "$UPLOAD_RESPONSE" | jq .
echo ""

# 提取文件ID和哈希
FILE_ID=$(echo "$UPLOAD_RESPONSE" | jq -r '.data.id')
FILE_HASH=$(echo "$UPLOAD_RESPONSE" | jq -r '.data.fileHash')

echo "文件ID: $FILE_ID"
echo "文件哈希: $FILE_HASH"
echo ""

# 如果上传失败，使用已知的文件ID进行测试
if [ "$FILE_ID" = "null" ] || [ -z "$FILE_ID" ]; then
    echo "上传失败，使用已知文件ID进行测试..."
    FILE_ID="2bfddf94-87d2-4d78-bb2f-c9ff77f202a0"
    FILE_HASH="855a3b4d416e91bccbfbbd0baf323fd57fe99007f60dd092324fee6baadb06da"
    echo "使用测试文件ID: $FILE_ID"
    echo "使用测试文件哈希: $FILE_HASH"
    echo ""
fi

# 2. 测试文件查询（通过ID）
echo "2. 测试文件查询（通过ID）..."
QUERY_RESPONSE=$(curl -s -X GET "$BASE_URL/api/esg/query?id=$FILE_ID")
echo "查询响应:"
echo "$QUERY_RESPONSE" | jq .
echo ""

# 3. 测试文件查询（通过哈希）
echo "3. 测试文件查询（通过哈希）..."
QUERY_HASH_RESPONSE=$(curl -s -X GET "$BASE_URL/api/esg/query?fileHash=$FILE_HASH")
echo "哈希查询响应:"
echo "$QUERY_HASH_RESPONSE" | jq .
echo ""

# 4. 测试文件列表
echo "4. 测试文件列表..."
LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/esg/list?page=1&pageSize=10")
echo "列表响应:"
echo "$LIST_RESPONSE" | jq .
echo ""

# 5. 测试文件统计
echo "5. 测试文件统计..."
STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/esg/stats")
echo "统计响应:"
echo "$STATS_RESPONSE" | jq .
echo ""

# 6. 上传第二个文件
echo "6. 上传第二个文件..."
# 使用JSON格式的测试内容
JSON_CONTENT='{"company": "Test Corp", "year": 2024, "sustainability_score": 85, "environmental_impact": "low", "social_responsibility": "high"}'
JSON_CONTENT_B64=$(echo "$JSON_CONTENT" | base64 -w 0)
UPLOAD_RESPONSE2=$(curl -s -X POST "$BASE_URL/api/esg/upload" \
  -H "Content-Type: application/json" \
  -d "{
    \"fileName\": \"sustainability_report.json\",
    \"fileContent\": \"$JSON_CONTENT_B64\",
    \"fileType\": \"json\",
    \"description\": \"可持续发展报告\"
  }")

echo "第二个文件上传响应:"
echo "$UPLOAD_RESPONSE2" | jq .
echo ""

# 7. 再次测试文件列表
echo "7. 再次测试文件列表..."
LIST_RESPONSE2=$(curl -s -X GET "$BASE_URL/api/esg/list?page=1&pageSize=10")
echo "更新后的列表响应:"
echo "$LIST_RESPONSE2" | jq .
echo ""

# 8. 测试文件删除
echo "8. 测试文件删除..."
DELETE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/esg/delete" \
  -H "Content-Type: application/json" \
  -d "{\"id\": \"$FILE_ID\"}")

echo "删除响应:"
echo "$DELETE_RESPONSE" | jq .
echo ""

# 9. 测试删除后的查询
echo "9. 测试删除后的查询..."
QUERY_DELETED_RESPONSE=$(curl -s -X GET "$BASE_URL/api/esg/query?id=$FILE_ID")
echo "删除后查询响应:"
echo "$QUERY_DELETED_RESPONSE" | jq .
echo ""

# 10. 最终统计
echo "10. 最终统计..."
FINAL_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/esg/stats")
echo "最终统计响应:"
echo "$FINAL_STATS_RESPONSE" | jq .
echo ""

echo "=== 测试完成 ===" 