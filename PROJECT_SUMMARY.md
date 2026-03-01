# 项目总结（go/fabric-sdk）

## 1. Fabric 网络连接方法

- 采用 `github.com/hyperledger/fabric-gateway` 官方 Go SDK，连接 Fabric 2.5.5 网络。
- 连接流程：
  1. 加载 Org1 Admin 的 X.509 证书和私钥（PEM 格式）。
  2. 加载 peer 节点的 TLS 根证书（如 peer0.org1.example.com 的 ca.crt）。
  3. 通过 `grpc.DialContext` 连接 peer 节点（支持 TLS/Insecure）。
  4. 使用 `identity.NewX509Identity` 和 `identity.NewPrivateKeySign` 创建身份与签名。
  5. 通过 `client.Connect` 创建 Gateway 客户端。
  6. 获取 Network（通道）和 Contract（链码合约）对象。
- 相关代码：`internal/service/chaincode.go` 的 `InitChaincodeService`。

## 2. 链码部署方法

- 推荐使用 test-network 脚本自动部署：
  ```bash
  cd ~/go/fabric-samples/test-network
  ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
  ```
- 手动部署流程：
  1. 打包链码
  2. 安装链码到各 peer
  3. 各组织批准链码定义
  4. 提交链码定义到通道
  5. 确认链码已就绪
- 链码升级需递增 sequence，所有组织批准后提交。

## 3. 链码调用方法

### 3.1 SDK 端调用（Go 代码）
- 写入记录：
  ```go
  result, err := contract.SubmitTransaction("WriteRecord", key, value)
  ```
- 读取记录：
  ```go
  result, err := contract.EvaluateTransaction("ReadRecord", key)
  ```
- 创建资产：
  ```go
  result, err := contract.SubmitTransaction("CreateAsset", assetID, color, size, owner, appraisedValue)
  ```
- 读取资产：
  ```go
  result, err := contract.EvaluateTransaction("ReadAsset", assetID)
  ```
- 批量查询：
  ```go
  result, err := contract.EvaluateTransaction("GetAllAssets")
  ```

### 3.2 HTTP API 路由
- `/api/chaincode/init`           POST  初始化链码服务
- `/api/chaincode/write-record`   POST  写入记录（参数：key, value）
- `/api/chaincode/read-record`    GET   读取记录（参数：key）
- `/api/chaincode/create-asset`   POST  创建资产（参数：assetID, color, size, owner, appraisedValue）
- `/api/chaincode/read-asset`     GET   读取资产（参数：assetID）
- `/api/chaincode/get-all-assets` GET   获取所有资产

### 3.3 CLI 调用（peer 命令）
- 写入：
  ```bash
  peer chaincode invoke -C mychannel -n basic -c '{"Args":["WriteRecord","test_key","test_value"]}' ...
  ```
- 读取：
  ```bash
  peer chaincode query -C mychannel -n basic -c '{"Args":["ReadRecord","test_key"]}'
  ```

## 4. 常见问题与最佳实践

- 初始化必须先调用 `/api/chaincode/init`，否则所有链码 API 都会提示“链码合约未初始化”。
- 账本批量接口（如 GetAllAssets）需保证链码实现返回合法 JSON，否则 SDK/CLI 都会报错。
- 证书、私钥、TLS 路径必须与 test-network 生成的实际文件完全一致。
- 链码升级需注意 sequence 递增，所有组织都需批准。
- 推荐开发环境下可跳过 TLS 校验（InsecureSkipVerify: true），生产环境务必开启。

## 5. 目录结构参考

- `internal/service/chaincode.go`：链码服务核心逻辑
- `internal/controller/chaincode.go`：HTTP API 控制器
- `main.go`、`internal/cmd/cmd.go`：服务启动与路由注册
- `go/fabric-samples/test-network/`：Fabric 网络与链码部署脚本

---

如需更详细的代码示例、配置文件说明或自动化脚本，请查阅源码或联系维护者。 

## 6. IPFS集成与大文件上传/下载全流程记录

### 6.1 环境准备
- 系统环境：Ubuntu 24.04 LTS
- IPFS版本：kubo v0.29.0
- GoFrame服务：go/fabric-sdk，端口8199

#### IPFS 安装与初始化
```bash
curl -L -o kubo_v0.29.0_linux-amd64.tar.gz https://github.com/ipfs/kubo/releases/download/v0.29.0/kubo_v0.29.0_linux-amd64.tar.gz
sudo tar -C /usr/local -xzf kubo_v0.29.0_linux-amd64.tar.gz
sudo mv /usr/local/kubo/ipfs /usr/local/bin/ipfs
ipfs --version  # 应输出 ipfs version 0.29.0
ipfs init
ipfs config --json Addresses.API '"/ip4/0.0.0.0/tcp/5001"'
ipfs config --json Addresses.Gateway '"/ip4/0.0.0.0/tcp/8080"'
nohup ipfs daemon > ~/ipfs.log 2>&1 &
```

### 6.2 GoFrame项目集成IPFS
- 零耦合：GoFrame 只需能访问 IPFS HTTP API（默认 http://127.0.0.1:5001），无需在 IPFS 侧做任何配置。
- 零依赖：仅用 Go 标准库 `net/http` 实现上传/下载，无需 go-ipfs-api 等第三方库。
- 环境变量支持：支持 `IPFS_API` 环境变量，未配置则默认 `http://127.0.0.1:5001`。

#### service 层
- `IpfsUpload(reader io.Reader, filename string) (string, error)`  
  构造 multipart/form-data，POST 到 `/api/v0/add`，解析返回 JSON，提取 CID。
- `IpfsDownload(cid string) (io.ReadCloser, error)`  
  GET `/api/v0/cat?arg=cid`，直接返回文件流。

#### controller 层
- `POST /api/ipfs/upload`  
  接收 form-data 文件，调用 `IpfsUpload`，返回 `{code:200, cid:xxx}`。
- `GET /api/ipfs/download?cid=xxx`  
  调用 `IpfsDownload`，流式返回文件内容。

#### 健壮性处理
- controller 层增加 panic 捕获、详细日志，任何异常都能返回明确 JSON。
- service 层对 IPFS 返回内容做健壮解析，CID 为空时返回错误。

### 6.3 测试与验证
- 上传：
```bash
cd ~/go/fabric-sdk
curl -F 'file=@测试.pdf' http://localhost:8199/api/ipfs/upload
```
- 下载：
```bash
curl -o 下载.pdf "http://localhost:8199/api/ipfs/download?cid=QmXETqAUKUfZ5cnHiu6LkuAykHCWPcr62qtkxM15UNtHLN"
```
- 直接调用IPFS原生API验证：
```bash
curl -X POST http://127.0.0.1:5001/api/v0/add -F 'file=@测试.pdf'
```

### 6.4 排障与修复
- 发现 go-ipfs-api 依赖因官方归档无法拉取，果断切换为标准库实现，彻底解决依赖问题。
- controller 层增加日志与 panic 捕获，保证接口无论成功失败都能返回明确 JSON。
- 多次用 curl 验证本地文件路径、API 可用性，确保链路全通。

### 6.5 最佳实践与建议
- 推荐用标准库直接访问 IPFS HTTP API，简单、健壮、无依赖。
- 所有上传/下载接口都应有详细日志和异常捕获，便于排障。
- 生产环境建议 IPFS 节点与 GoFrame 服务同机或内网，保证带宽与安全。
- 如需跨主机访问，记得开放 5001 端口并配置防火墙。

--- 