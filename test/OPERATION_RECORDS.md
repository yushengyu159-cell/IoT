# Fabric网络操作记录

**记录日期**: 2026-03-25  
**服务器**: 47.238.159.234

---

## 🔧 快速操作命令

### 1. Fabric网络重启
```bash
# 关闭Fabric网络
cd /root/home/go/fabric-samples/test-network
./network.sh down

# 启动网络并创建通道
./network.sh up createChannel
```

### 2. 部署链码
```bash
# 部署Go语言链码
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go

# 或部署JavaScript链码
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-javascript -ccl javascript
```

### 3. 启动其他容器
```bash
# 启动MySQL
docker start mysql-db

# 启动IPFS
docker start ipfs-daemon

# 启动fabric-sdk应用
cd /root/home/go/fabric-sdk
./deploy.sh
```

---

## 📝 完整重启流程

```bash
# 1. 关闭Fabric网络
cd /root/home/go/fabric-samples/test-network
./network.sh down

# 2. 启动网络并创建通道
./network.sh up createChannel

# 3. 部署链码
./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go

# 4. 启动其他容器
docker start mysql-db
docker start ipfs-daemon

# 5. 启动应用
cd /root/home/go/fabric-sdk
./deploy.sh
```

---

## ✅ 验证服务状态

```bash
# 查看所有容器
docker ps

# 查看Fabric容器
docker ps | grep hyperledger

# 测试应用
curl http://localhost:8199/api/index
```

---

**最后更新**: 2026-03-25
