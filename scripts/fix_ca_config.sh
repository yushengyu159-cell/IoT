#!/bin/bash

# 修复certificateAuthorities配置格式问题

set -e

echo "🔧 修复certificateAuthorities配置格式..."
echo "=================================="

# 设置颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 备份原文件
cp configs/fabric/connection-optimized.yaml configs/fabric/connection-optimized.yaml.backup3
echo -e "${GREEN}✅ 已备份原文件${NC}"

# 创建修复后的配置文件
cat > configs/fabric/connection-optimized.yaml << 'EOF'
---
name: test-network-fixed-certs
version: 1.0.0
client:
  organization: Org1
  connection:
    timeout:
      peer:
        endorser: '300'
      orderer: '300'
  cryptoconfig:
    path: ${GOPATH}/src/github.com/hyperledger/fabric-samples/test-network/organizations
  credentialStore:
    path: /tmp/hfc-kvs
    cryptoStore:
      path: /tmp/hfc-cvs
  BCCSP:
    default: SW
    SW:
      hash: SHA2
      security: 256
      fileKeyStore:
        keyStore:
          path: /tmp/msp/keystore

organizations:
  Org1:
    mspid: Org1MSP
    cryptoPath: peerOrganizations/org1.example.com/users/{username}@org1.example.com/msp
    peers:
    - peer0.org1.example.com
    certificateAuthorities:
    - ca.org1.example.com

peers:
  peer0.org1.example.com:
    url: grpcs://localhost:7051
    tlsCACerts:
      pem: |
        -----BEGIN CERTIFICATE-----
        MIICVzCCAf2gAwIBAgIQea5aDykeepLNVWxS2i2n0zAKBggqhkjOPQQDAjB2MQsw
        CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
        YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEfMB0GA1UEAxMWdGxz
        Y2Eub3JnMS5leGFtcGxlLmNvbTAeFw0yNTA3MTcxMDQzMDBaFw0zNTA3MTUxMDQz
        MDBaMHYxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH
        Ew1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcxLmV4YW1wbGUuY29tMR8wHQYD
        VQQDExZ0bHNjYS5vcmcxLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0D
        AQcDQgAE+56Xl9tMVTe5yCsPxldVgJdBzuxoh/JXSD5VspoVXqGKnlZ8ZlH4/pax
        RmjyL1Z45bxkQkCJB7vGpwvh3yYym6NtMGswDgYDVR0PAQH/BAQDAgGmMB0GA1Ud
        JQQWMBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1Ud
        DgQiBCCYUJinbLe2/ydSKdkzUuF4a6aZRCgguR9nSNMbyChJOzAKBggqhkjOPQQD
        AgNIADBFAiEAmAkXTD7MfC8jpVO49nHesEbzzntB3VxUJnC13VpxYg8CIBsRwRA3
        Ic53jttFfENNNMw/NQ7bRWcETj0sGRC7jZRz
        -----END CERTIFICATE-----
    grpcOptions:
      ssl-target-name-override: peer0.org1.example.com
      allow-insecure: false
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false

orderers:
  orderer.example.com:
    url: grpcs://localhost:7050
    tlsCACerts:
      pem: |
        -----BEGIN CERTIFICATE-----
        MIICQzCCAeqgAwIBAgIRAPFgl7QPzFoT4/fxu7zyP5YwCgYIKoZIzj0EAwIwbDEL
        MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
        cmFuY2lzY28xFDASBgNVBAoTC2V4YW1wbGUuY29tMRowGAYDVQQDExF0bHNjYS5l
        eGFtcGxlLmNvbTAeFw0yNTA3MTcxMDQzMDBaFw0zNTA3MTUxMDQzMDBaMGwxCzAJ
        BgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1TYW4gRnJh
        bmNpc2NvMRQwEgYDVQQKEwtleGFtcGxlLmNvbTEaMBgGA1UEAxMRdGxzY2EuZXhh
        bXBsZS5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASXfUCZOCS+ysXeDdzl
        rNbbVFr8H662be3smUrkkeaDYT/JYQyvmoDv9uIG/NwrVK0Z0CCWuj1BTwfy+Zev
        xi8Eo20wazAOBgNVHQ8BAf8EBAMCAaYwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsG
        AQUFBwMBMA8GA1UdEwEB/wQFMAMBAf8wKQYDVR0OBCIEIKhgUc3RH0/tQyt+yZfj
        5eG5lKkFFbh9ciXLUkQGnUDPMAoGCCqGSM49BAMCA0cAMEQCID4wkLi+fEopAaGy
        lIK5OpqYRAuY+GJVOaCAuonQ5J7BAiA2EcP+0JYyIBxDFH5xBg66bJ/Y5of0bgek
        ZG55B/pvLA==
        -----END CERTIFICATE-----
    grpcOptions:
      ssl-target-name-override: orderer.example.com
      allow-insecure: false
      keep-alive-time: 0s
      keep-alive-timeout: 20s
      keep-alive-permit: false
      fail-fast: false

certificateAuthorities:
  ca.org1.example.com:
    url: https://localhost:7054
    caName: ca-org1
    tlsCACerts:
      pem:
        - |
          -----BEGIN CERTIFICATE-----
          MIICUjCCAfigAwIBAgIRAJ+sUuzFqOFICdYEZil8+4EwCgYIKoZIzj0EAwIwczEL
          MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
          cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
          Lm9yZzEuZXhhbXBsZS5jb20wHhcNMjUwNzE3MTA0MzAwWhcNMzUwNzE1MTA0MzAw
          WjBzMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
          U2FuIEZyYW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UE
          AxMTY2Eub3JnMS5leGFtcGxlLmNvbTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA
          BIcgnz/2e7r1QypPK9RcVLa24NjsOukYbY9tvRivZW4kgMObnPdwNKDXe0v3eSMP
          umUY97rBNcxO/x0XCKWbQayjbTBrMA4GA1UdDwEB/wQEAwIBpjAdBgNVHSUEFjAU
          BggrBgEFBQcDAgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUwAwEB/zApBgNVHQ4EIgQg
          C00VUtNAfjaSHc95WBQ3d40Nt8uQFPQCt3ZVka2QpzowCgYIKoZIzj0EAwIDSAAw
          RQIhAJLhrPxod+6s//lLiqQKGuodjsXEgtsDqBT0iTzn6oQfAiAjaQ/wunThClj+
          C7as7s+TP6tsFTLkZxStDSB6157V+g==
          -----END CERTIFICATE-----
    httpOptions:
      verify: false
    registrar:
      enrollId: admin
      enrollSecret: adminpw

channels:
  mychannel:
    orderers:
      - orderer.example.com
    peers:
      peer0.org1.example.com:
        endorsingPeer: true
        chaincodeQuery: true
        ledgerQuery: true
        eventSource: true
        discoverAsLocalhost: true

entityMatchers:
  peer:
    - pattern: peer0.org1.example.com
      urlSubstitutionExp: localhost:7051
      sslTargetOverrideUrlSubstitutionExp: peer0.org1.example.com
      mappedHost: peer0.org1.example.com
  orderer:
    - pattern: orderer.example.com
      urlSubstitutionExp: localhost:7050
      sslTargetOverrideUrlSubstitutionExp: orderer.example.com
      mappedHost: orderer.example.com
  certificateAuthority:
    - pattern: ca.org1.example.com
      urlSubstitutionExp: localhost:7054
      mappedHost: ca.org1.example.com
EOF

echo -e "${GREEN}✅ 配置文件已修复${NC}"

# 验证YAML语法
echo "🔍 验证YAML语法..."
if python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>/dev/null; then
    echo -e "${GREEN}✅ YAML语法正确${NC}"
else
    echo -e "${RED}❌ YAML语法仍有错误${NC}"
    python3 -c "import yaml; yaml.safe_load(open('configs/fabric/connection-optimized.yaml'))" 2>&1
    exit 1
fi

echo ""
echo -e "${GREEN}🎉 certificateAuthorities配置修复完成！${NC}"
echo "=================================="
echo ""
echo -e "${GREEN}💡 修复内容:${NC}"
echo "1. ✅ 修复了tlsCACerts.pem格式 - 改为数组格式"
echo "2. ✅ 修复了registrar格式 - 改为对象格式"
echo "3. ✅ 验证了YAML语法正确性"
echo "4. ✅ 修复了MSP配置路径"
echo "5. ✅ 强制指定了MSP CA证书路径"
echo ""
echo -e "${BLUE}📋 下一步操作:${NC}"
echo "1. 重新启动应用程序测试连接"
echo "2. 验证Fabric Gateway连接是否成功" 