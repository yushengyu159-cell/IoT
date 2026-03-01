#!/usr/bin/env bash
set -e

MSP_BASE="${1:-/home/ubuntu/go/fabric-sdk/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp}"

echo "🔍 Checking MSP: $MSP_BASE"
echo "--------------------------------------------------"

# 1. 目录结构
echo "📁 Directory structure:"
ls -l "$MSP_BASE" || { echo "❌ MSP dir not found"; exit 1; }

# 2. CA 证书存在
CA_CERT="$MSP_BASE/../msp/cacerts/ca.org1.example.com-cert.pem"
if [[ -f "$CA_CERT" ]]; then
  echo "✅ CA cert found: $CA_CERT"
else
  echo "❌ CA cert missing: $CA_CERT"; exit 1
fi

# 3. 用户证书存在
USER_CERT="$MSP_BASE/signcerts/Admin@org1.example.com-cert.pem"
if [[ -f "$USER_CERT" ]]; then
  echo "✅ User cert found: $USER_CERT"
else
  echo "❌ User cert missing: $USER_CERT"; exit 1
fi

# 4. 证书链验证
if openssl verify -CAfile "$CA_CERT" "$USER_CERT" >/dev/null 2>&1; then
  echo "✅ Certificate chain valid"
else
  echo "❌ Certificate chain invalid"; exit 1
fi

echo "🎉 MSP configuration looks good!" 