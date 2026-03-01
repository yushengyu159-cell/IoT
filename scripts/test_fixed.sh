#!/bin/bash

# 修复后的功能测试脚本
echo "========== 修复后的Fabric SDK测试 =========="

BASE_URL="http://localhost:8199"

# 等待服务启动
echo "等待服务启动..."
sleep 3

# 测试1: 健康检查
echo -e "\n1. 测试健康检查..."
curl -s "${BASE_URL}/health/" | jq '.'

# 测试2: 就绪检查
echo -e "\n2. 测试就绪检查..."
curl -s "${BASE_URL}/health/ready" | jq '.'

# 测试3: 存活检查
echo -e "\n3. 测试存活检查..."
curl -s "${BASE_URL}/health/live" | jq '.'

# 测试4: 指标接口
echo -e "\n4. 测试指标接口..."
curl -s "${BASE_URL}/health/metrics" | jq '.'

# 测试5: Fabric连接状态
echo -e "\n5. 测试Fabric连接状态..."
curl -s "${BASE_URL}/api/fabric/status" | jq '.'

# 测试6: ESG文件上传
echo -e "\n6. 测试ESG文件上传..."
UPLOAD_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/esg/upload" \
  -H "Content-Type: application/json" \
  -d '{
    "fileName": "测试.pdf",
    "fileContent": "修复后的测试文件内容",
    "fileType": "pdf",
    "description": "修复后的功能验证测试"
  }')

echo "$UPLOAD_RESPONSE" | jq '.'
FILE_ID=$(echo "$UPLOAD_RESPONSE" | jq -r '.data.fileId')

# 测试7: 文件查询
echo -e "\n7. 测试文件查询..."
curl -s "${BASE_URL}/api/esg/query?fileId=${FILE_ID}" | jq '.'

# 测试8: 文件列表
echo -e "\n8. 测试文件列表..."
curl -s "${BASE_URL}/api/esg/list?page=1&pageSize=10" | jq '.'

# 测试9: 统计信息
echo -e "\n9. 测试统计信息..."
curl -s "${BASE_URL}/api/esg/stats" | jq '.'

# 测试10: 文件删除
echo -e "\n10. 测试文件删除..."
curl -s -X POST "${BASE_URL}/api/esg/delete" \
  -H "Content-Type: application/json" \
  -d "{\"fileId\": \"${FILE_ID}\"}" | jq '.'

# 测试11: 验证删除后的状态
echo -e "\n11. 验证删除后的状态..."
curl -s "${BASE_URL}/api/esg/query?fileId=${FILE_ID}" | jq '.'

# 测试12: 最终健康状态
echo -e "\n12. 最终健康状态..."
curl -s "${BASE_URL}/health/" | jq '.'

echo -e "\n========== 测试完成 =========="
echo "如果所有测试都通过，说明修复成功！" 