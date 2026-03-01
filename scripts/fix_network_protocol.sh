#!/bin/bash

echo "🎯 修复网络协议：从TLS切换到非TLS连接"

# 1. 检查配置文件中的协议设置
echo "🔐 1. 检查网络协议设置..."
echo "   peer0.org1.example.com URL"
if grep -q "url: grpc://localhost:7051" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 已使用非TLS协议 (grpc://)"
else
    echo "   ❌ 仍在使用TLS协议 (grpcs://)"
    exit 1
fi

echo "   orderer.example.com URL"
if grep -q "url: grpc://localhost:7050" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 已使用非TLS协议 (grpc://)"
else
    echo "   ❌ 仍在使用TLS协议 (grpcs://)"
    exit 1
fi

echo "   ca.org1.example.com URL"
if grep -q "url: http://localhost:7054" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 已使用非TLS协议 (http://)"
else
    echo "   ❌ 仍在使用TLS协议 (https://)"
    exit 1
fi

# 2. 检查TLS验证设置
echo "🔒 2. 检查TLS验证设置..."
if grep -q "verify: false" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 证书验证已关闭"
else
    echo "   ❌ 证书验证未关闭"
    exit 1
fi

if grep -q "allow-insecure: true" configs/fabric/connection-optimized.yaml; then
    echo "   ✅ 不安全连接已允许"
else
    echo "   ❌ 不安全连接未允许"
    exit 1
fi

# 3. 检查环境变量
echo "🌍 3. 检查环境变量..."
export FABRIC_SDK_GO_TLS_VERIFY=false
export FABRIC_SDK_GO_CERT_VERIFY=false
export FABRIC_SDK_GO_SYSTEM_CERT_POOL=false
export FABRIC_SDK_GO_ALLOW_INSECURE=true

echo "   FABRIC_SDK_GO_TLS_VERIFY: $FABRIC_SDK_GO_TLS_VERIFY"
echo "   FABRIC_SDK_GO_CERT_VERIFY: $FABRIC_SDK_GO_CERT_VERIFY"
echo "   FABRIC_SDK_GO_SYSTEM_CERT_POOL: $FABRIC_SDK_GO_SYSTEM_CERT_POOL"
echo "   FABRIC_SDK_GO_ALLOW_INSECURE: $FABRIC_SDK_GO_ALLOW_INSECURE"

# 4. 测试网络连接
echo "🌐 4. 测试网络连接..."
if timeout 5 bash -c "</dev/tcp/localhost/7051" 2>/dev/null; then
    echo "   ✅ peer0.org1:7051 连接正常"
else
    echo "   ❌ peer0.org1:7051 连接失败"
    exit 1
fi

if timeout 5 bash -c "</dev/tcp/localhost/7050" 2>/dev/null; then
    echo "   ✅ orderer:7050 连接正常"
else
    echo "   ❌ orderer:7050 连接失败"
    exit 1
fi

# 5. 检查配置文件语法
echo "🔍 5. 检查配置文件语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo "   ✅ YAML配置文件语法正确"
else
    echo "   ❌ YAML配置文件语法错误"
    exit 1
fi

# 6. 显示修复状态
echo "📊 6. 网络协议修复状态..."
echo "   - Peer协议: ✅ 已切换到 grpc:// (非TLS)"
echo "   - Orderer协议: ✅ 已切换到 grpc:// (非TLS)"
echo "   - CA协议: ✅ 已切换到 http:// (非TLS)"
echo "   - TLS验证: ✅ 已关闭"
echo "   - 不安全连接: ✅ 已允许"
echo "   - 网络连接: ✅ 所有端口可连接"
echo "   - 配置文件: ✅ 语法正确"

echo "🎉 网络协议修复完成！"
echo ""
echo "🚀 现在重启SDK服务进行测试："
echo "   pkill -f main.go"
echo "   go run main.go"
echo ""
echo "💡 修复说明："
echo "   - 将 grpcs:// 改为 grpc:// (移除TLS)"
echo "   - 将 https:// 改为 http:// (移除TLS)"
echo "   - 保持 allow-insecure: true 设置"
echo "   - 这样SDK将使用非TLS连接，避免证书验证问题" 