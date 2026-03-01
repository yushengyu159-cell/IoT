#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
文件上传接口测试脚本
测试账号: 18027473816@163.com
接口地址: https://app.esgvisa.com/api/esg/upload-encrypted

使用方法:
    python3 test_upload.py [文件路径] [描述]
    或
    python3 test_upload.py  # 使用默认test.pdf
"""

import requests
import json
import sys
import os
from datetime import datetime

# 配置信息
API_URL = "https://app.esgvisa.com/api/esg/upload-encrypted"
UPLOADER = "18027473816@163.com"
DEFAULT_TEST_FILE = "test.pdf"

def upload_file(file_path, desc="自动化测试文件上传", uploader=UPLOADER):
    """
    上传文件到服务器
    
    Args:
        file_path: 文件路径
        desc: 文件描述
        uploader: 上传者账号
    
    Returns:
        dict: 响应数据
    """
    print("=" * 50)
    print("文件上传接口测试")
    print("=" * 50)
    print(f"接口地址: {API_URL}")
    print(f"测试账号: {uploader}")
    print(f"测试文件: {file_path}")
    print("=" * 50)
    print()
    
    # 检查文件是否存在
    if not os.path.exists(file_path):
        print(f"❌ 错误: 文件 {file_path} 不存在")
        return None
    
    # 获取文件大小
    file_size = os.path.getsize(file_path)
    print(f"📄 文件信息:")
    print(f"  文件名: {os.path.basename(file_path)}")
    print(f"  文件大小: {file_size:,} 字节 ({file_size / 1024 / 1024:.2f} MB)")
    print()
    
    # 准备上传数据
    try:
        with open(file_path, 'rb') as f:
            files = {
                'file': (os.path.basename(file_path), f, 'application/octet-stream')
            }
            data = {
                'desc': f"{desc} - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}",
                'uploader': uploader
            }
            
            print("📤 开始上传文件...")
            print()
            
            # 发送请求
            response = requests.post(
                API_URL,
                files=files,
                data=data,
                timeout=300  # 5分钟超时
            )
    except FileNotFoundError:
        print(f"❌ 错误: 文件 {file_path} 不存在")
        return None
    
    # 检查HTTP状态码
    print(f"HTTP状态码: {response.status_code}")
    print()
    
    # 解析响应
    try:
        result = response.json()
    except json.JSONDecodeError:
        print("❌ 错误: 响应不是有效的JSON格式")
        print(f"响应内容: {response.text[:500]}")
        return None
    
    print("响应内容:")
    print(json.dumps(result, indent=2, ensure_ascii=False))
    print()
    
    # 检查业务状态码
    if response.status_code == 200 and result.get('code') == 200:
        print("✅ 上传成功!")
        print()
        
        # 提取关键信息
        data_info = result.get('data', {})
        meta = data_info.get('meta', {})
        
        print("📋 文件信息:")
        print(f"  文件名: {meta.get('fileName', 'N/A')}")
        print(f"  文件大小: {meta.get('fileSize', 0):,} 字节")
        print(f"  分片数量: {data_info.get('chunkCount', 0)}")
        print(f"  上传时间: {meta.get('uploadAt', 'N/A')}")
        
        time_stats = data_info.get('timeStats', {})
        if time_stats:
            print(f"  上传耗时: {time_stats.get('totalTime', 'N/A')}")
        
        print()
        print("🔑 关键数据:")
        cids = meta.get('cids', [])
        if cids:
            print(f"  主CID: {cids[0]}")
        print(f"  加密密钥: {data_info.get('key', 'N/A')}")
        print()
        print("⚠️  请妥善保管加密密钥，用于后续文件下载")
        
        # 保存结果到文件
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        result_file = f"upload_result_{timestamp}.json"
        with open(result_file, 'w', encoding='utf-8') as f:
            json.dump(result, f, indent=2, ensure_ascii=False)
        print()
        print(f"✅ 响应已保存到文件: {result_file}")
        
        return result
    else:
        print("❌ 上传失败:")
        print(f"  错误码: {result.get('code', 'N/A')}")
        print(f"  错误消息: {result.get('message', 'N/A')}")
        return None

if __name__ == "__main__":
    # 检查命令行参数
    if len(sys.argv) > 1:
        file_path = sys.argv[1]
        desc = sys.argv[2] if len(sys.argv) > 2 else "命令行测试上传"
        result = upload_file(file_path, desc)
        sys.exit(0 if result else 1)
    else:
        # 使用默认测试文件
        if os.path.exists(DEFAULT_TEST_FILE):
            result = upload_file(DEFAULT_TEST_FILE)
            sys.exit(0 if result else 1)
        else:
            print(f"❌ 错误: 默认测试文件 {DEFAULT_TEST_FILE} 不存在")
            print()
            print("使用方法:")
            print(f"  python3 {sys.argv[0]} [文件路径] [描述]")
            print("  或创建测试文件: test.pdf")
            sys.exit(1)




