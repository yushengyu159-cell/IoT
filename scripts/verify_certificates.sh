#!/bin/bash

# 证书验证和配置优化脚本
# 检查证书是否被重新生成，验证证书链完整性

set -e

echo "🔍 开始证书验证和配置优化..."
echo "=================================="

# 设置颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查证书文件时间戳
echo -e "${BLUE}📅 检查证书文件时间戳...${NC}"

# 检查关键证书文件
CERT_FILES=(
    "configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
    "configs/fabric/organizations/peerOrganizations/org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"
    "configs/fabric/organizations/peerOrganizations/org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem"
    "configs/fabric/organizations/ordererOrganizations/example.com/msp/cacerts/ca.example.com-cert.pem"
    "configs/fabric/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
)

for cert_file in "${CERT_FILES[@]}"; do
    if [ -f "$cert_file" ]; then
        timestamp=$(stat -c %Y "$cert_file")
        date_str=$(date -d @$timestamp '+%Y-%m-%d %H:%M:%S')
        echo -e "${GREEN}✅ $cert_file - $date_str${NC}"
    else
        echo -e "${RED}❌ $cert_file - 文件不存在${NC}"
    fi
done

echo ""

# 验证证书链
echo -e "${BLUE}🔗 验证证书链完整性...${NC}"

# 检查根CA证书
ROOT_CA="configs/fabric/organizations/ordererOrganizations/example.com/msp/cacerts/ca.example.com-cert.pem"
if [ -f "$ROOT_CA" ]; then
    echo -e "${GREEN}✅ 根CA证书存在${NC}"
    # 验证证书格式
    if openssl x509 -in "$ROOT_CA" -text -noout > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 根CA证书格式正确${NC}"
    else
        echo -e "${RED}❌ 根CA证书格式错误${NC}"
    fi
else
    echo -e "${RED}❌ 根CA证书不存在${NC}"
fi

# 检查TLS CA证书
TLS_CA="configs/fabric/organizations/ordererOrganizations/example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
if [ -f "$TLS_CA" ]; then
    echo -e "${GREEN}✅ TLS CA证书存在${NC}"
    # 验证证书格式
    if openssl x509 -in "$TLS_CA" -text -noout > /dev/null 2>&1; then
        echo -e "${GREEN}✅ TLS CA证书格式正确${NC}"
    else
        echo -e "${RED}❌ TLS CA证书格式错误${NC}"
    fi
else
    echo -e "${RED}❌ TLS CA证书不存在${NC}"
fi

# 检查用户证书
USER_CERT="configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
if [ -f "$USER_CERT" ]; then
    echo -e "${GREEN}✅ 用户证书存在${NC}"
    # 验证证书格式
    if openssl x509 -in "$USER_CERT" -text -noout > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 用户证书格式正确${NC}"
    else
        echo -e "${RED}❌ 用户证书格式错误${NC}"
    fi
else
    echo -e "${RED}❌ 用户证书不存在${NC}"
fi

echo ""

# 检查私钥文件
echo -e "${BLUE}🔐 检查私钥文件...${NC}"

KEY_DIR="configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore"
if [ -d "$KEY_DIR" ]; then
    key_files=$(find "$KEY_DIR" -name "*.pem" -o -name "*_sk")
    if [ -n "$key_files" ]; then
        echo -e "${GREEN}✅ 私钥文件存在${NC}"
        for key_file in $key_files; do
            echo -e "${GREEN}  - $key_file${NC}"
        done
    else
        echo -e "${RED}❌ 私钥文件不存在${NC}"
    fi
else
    echo -e "${RED}❌ 私钥目录不存在${NC}"
fi

echo ""

# 验证证书链关系
echo -e "${BLUE}🔗 验证证书链关系...${NC}"

# 检查用户证书和Org1 CA证书
USER_CERT="configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
ORG1_CA="configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/cacerts/ca.org1.example.com-cert.pem"

if [ -f "$USER_CERT" ] && [ -f "$ORG1_CA" ]; then
    # 提取用户证书的颁发者
    user_issuer=$(openssl x509 -in "$USER_CERT" -noout -issuer 2>/dev/null | sed 's/issuer=//')
    # 提取Org1 CA证书的主题
    org1_subject=$(openssl x509 -in "$ORG1_CA" -noout -subject 2>/dev/null | sed 's/subject=//')
    
    if [ "$user_issuer" = "$org1_subject" ]; then
        echo -e "${GREEN}✅ 用户证书由Org1 CA正确签发${NC}"
    else
        echo -e "${RED}❌ 用户证书不是由Org1 CA签发${NC}"
        echo -e "${YELLOW}   用户证书颁发者: $user_issuer${NC}"
        echo -e "${YELLOW}   Org1 CA主题: $org1_subject${NC}"
    fi
fi

# 检查根CA证书
if [ -f "$USER_CERT" ] && [ -f "$ROOT_CA" ]; then
    # 提取用户证书的颁发者
    user_issuer=$(openssl x509 -in "$USER_CERT" -noout -issuer 2>/dev/null | sed 's/issuer=//')
    # 提取根CA证书的主题
    root_subject=$(openssl x509 -in "$ROOT_CA" -noout -subject 2>/dev/null | sed 's/subject=//')
    
    if [ "$user_issuer" = "$root_subject" ]; then
        echo -e "${GREEN}✅ 用户证书由根CA签发${NC}"
    else
        echo -e "${YELLOW}⚠️ 用户证书可能不是由根CA签发${NC}"
        echo -e "${YELLOW}   用户证书颁发者: $user_issuer${NC}"
        echo -e "${YELLOW}   根CA主题: $root_subject${NC}"
    fi
fi

echo ""

# 检查配置文件中的证书内容
echo -e "${BLUE}📋 检查配置文件中的证书内容...${NC}"

CONFIG_FILE="configs/fabric/connection-optimized.yaml"
if [ -f "$CONFIG_FILE" ]; then
    echo -e "${GREEN}✅ 优化配置文件存在${NC}"
    
    # 检查是否包含完整的PEM证书
    if grep -q "-----BEGIN CERTIFICATE-----" "$CONFIG_FILE"; then
        echo -e "${GREEN}✅ 配置文件包含PEM证书${NC}"
    else
        echo -e "${RED}❌ 配置文件不包含PEM证书${NC}"
    fi
    
    # 检查MSP配置
    if grep -q "caCerts:" "$CONFIG_FILE"; then
        echo -e "${GREEN}✅ 配置文件包含MSP caCerts配置${NC}"
    else
        echo -e "${RED}❌ 配置文件不包含MSP caCerts配置${NC}"
    fi
    
    if grep -q "tlscacerts:" "$CONFIG_FILE"; then
        echo -e "${GREEN}✅ 配置文件包含MSP tlscacerts配置${NC}"
    else
        echo -e "${RED}❌ 配置文件不包含MSP tlscacerts配置${NC}"
    fi
else
    echo -e "${RED}❌ 优化配置文件不存在${NC}"
fi

echo ""

# 创建证书备份
echo -e "${BLUE}💾 创建证书备份...${NC}"

BACKUP_DIR="configs/fabric/cert_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$BACKUP_DIR"

# 备份关键证书
cp -r configs/fabric/organizations "$BACKUP_DIR/"
echo -e "${GREEN}✅ 证书备份完成: $BACKUP_DIR${NC}"

echo ""

# 生成证书验证报告
echo -e "${BLUE}📊 生成证书验证报告...${NC}"

REPORT_FILE="cert_verification_report_$(date +%Y%m%d_%H%M%S).txt"
{
    echo "证书验证报告 - $(date)"
    echo "=================================="
    echo ""
    echo "证书文件状态:"
    for cert_file in "${CERT_FILES[@]}"; do
        if [ -f "$cert_file" ]; then
            timestamp=$(stat -c %Y "$cert_file")
            date_str=$(date -d @$timestamp '+%Y-%m-%d %H:%M:%S')
            echo "✅ $cert_file - $date_str"
        else
            echo "❌ $cert_file - 文件不存在"
        fi
    done
    echo ""
    echo "证书链验证:"
    if [ -f "$USER_CERT" ] && [ -f "$ROOT_CA" ]; then
        user_issuer=$(openssl x509 -in "$USER_CERT" -noout -issuer 2>/dev/null | sed 's/issuer=//')
        root_subject=$(openssl x509 -in "$ROOT_CA" -noout -subject 2>/dev/null | sed 's/subject=//')
        if [ "$user_issuer" = "$root_subject" ]; then
            echo "✅ 用户证书由根CA签发"
        else
            echo "⚠️ 用户证书可能不是由根CA签发"
            echo "   用户证书颁发者: $user_issuer"
            echo "   根CA主题: $root_subject"
        fi
    fi
} > "$REPORT_FILE"

echo -e "${GREEN}✅ 验证报告已生成: $REPORT_FILE${NC}"

echo ""
echo -e "${GREEN}🎉 证书验证完成！${NC}"
echo "=================================="

# 提供建议
echo -e "${YELLOW}💡 建议:${NC}"
echo "1. 如果证书最近被重新生成，请更新配置文件中的证书内容"
echo "2. 确保所有证书文件的时间戳一致"
echo "3. 验证证书链的完整性"
echo "4. 检查私钥文件是否存在且可读"
echo "5. 确保配置文件中的证书内容是最新的"
echo "6. 检查MSP配置路径是否正确"
echo "7. 验证Fabric SDK环境变量设置"

echo ""
echo -e "${BLUE}📋 下一步操作:${NC}"
echo "1. 检查验证报告: $REPORT_FILE"
echo "2. 查看证书备份: $BACKUP_DIR"
echo "3. 如果需要，运行证书更新脚本" 