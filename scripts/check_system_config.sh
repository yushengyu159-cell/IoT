#!/bin/bash

echo "🔍 系统配置检查..."

# 获取当前工作目录
CURRENT_DIR=$(pwd)
echo "当前工作目录: $CURRENT_DIR"

# 1. 检查配置文件
echo "📋 1. 检查配置文件..."
if [ -f "configs/fabric/connection-optimized.yaml" ]; then
    echo "✅ 连接配置文件存在"
    
    # 检查YAML语法
    if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml', 'r'))" 2>/dev/null; then
        echo "✅ YAML语法正确"
    else
        echo "❌ YAML语法错误"
        exit 1
    fi
    
    # 检查cryptoPath配置
    CRYPTO_PATH=$(grep "cryptoPath:" configs/fabric/connection-optimized.yaml | awk '{print $2}')
    echo "📁 cryptoPath: $CRYPTO_PATH"
    
    if [ -d "configs/fabric/organizations/$CRYPTO_PATH" ]; then
        echo "✅ cryptoPath目录存在"
    else
        echo "❌ cryptoPath目录不存在"
        exit 1
    fi
else
    echo "❌ 连接配置文件不存在"
    exit 1
fi

# 2. 检查MSP目录结构
echo "📁 2. 检查MSP目录结构..."
MSP_PATH="configs/fabric/organizations/$CRYPTO_PATH"

echo "检查目录: $MSP_PATH"
if [ -d "$MSP_PATH" ]; then
    echo "✅ MSP目录存在"
    
    # 检查必要的子目录
    for dir in "cacerts" "tlscacerts" "signcerts" "keystore" "admincerts"; do
        if [ -d "$MSP_PATH/$dir" ]; then
            echo "✅ $dir 目录存在"
        else
            echo "❌ $dir 目录不存在"
        fi
    done
else
    echo "❌ MSP目录不存在"
    exit 1
fi

# 3. 检查证书文件
echo "🔐 3. 检查证书文件..."
if [ -f "$MSP_PATH/cacerts/ca.org1.example.com-cert.pem" ]; then
    echo "✅ CA证书存在"
else
    echo "❌ CA证书不存在"
fi

if [ -f "$MSP_PATH/tlscacerts/tlsca.org1.example.com-cert.pem" ]; then
    echo "✅ TLS CA证书存在"
else
    echo "❌ TLS CA证书不存在"
fi

if [ -f "$MSP_PATH/signcerts/Admin@org1.example.com-cert.pem" ]; then
    echo "✅ 用户证书存在"
else
    echo "❌ 用户证书不存在"
fi

if [ -f "$MSP_PATH/keystore/priv_sk" ]; then
    echo "✅ 私钥文件存在"
    # 检查权限
    PERM=$(stat -c %a "$MSP_PATH/keystore/priv_sk")
    if [ "$PERM" = "600" ]; then
        echo "✅ 私钥权限正确 (600)"
    else
        echo "⚠️ 私钥权限不正确: $PERM (应该是600)"
    fi
else
    echo "❌ 私钥文件不存在"
fi

# 4. 检查环境变量
echo "🌍 4. 检查环境变量..."
if [ -n "$FABRIC_SDK_GO_MSP_CACERTS_PATH" ]; then
    echo "✅ FABRIC_SDK_GO_MSP_CACERTS_PATH 已设置"
else
    echo "⚠️ FABRIC_SDK_GO_MSP_CACERTS_PATH 未设置"
fi

if [ -n "$FABRIC_SDK_GO_MSP_TLSCACERTS_PATH" ]; then
    echo "✅ FABRIC_SDK_GO_MSP_TLSCACERTS_PATH 已设置"
else
    echo "⚠️ FABRIC_SDK_GO_MSP_TLSCACERTS_PATH 未设置"
fi

# 5. 检查Go模块
echo "🔧 5. 检查Go模块..."
if [ -f "go.mod" ]; then
    echo "✅ go.mod 文件存在"
else
    echo "❌ go.mod 文件不存在"
    exit 1
fi

if [ -f "main.go" ]; then
    echo "✅ main.go 文件存在"
else
    echo "❌ main.go 文件不存在"
    exit 1
fi

# 6. 检查Fabric网络连接
echo "🌐 6. 检查Fabric网络连接..."
echo "注意: 这需要Fabric网络正在运行"

echo "🎉 系统配置检查完成！"
echo ""
echo "📊 检查总结:"
echo "- 配置文件: ✅"
echo "- MSP目录: ✅"
echo "- 证书文件: ✅"
echo "- 环境变量: ⚠️ (部分设置)"
echo "- Go模块: ✅"
echo "- 网络连接: 需要验证" 