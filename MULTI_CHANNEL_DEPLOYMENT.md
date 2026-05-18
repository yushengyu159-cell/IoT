# 多通道分片方案部署文档

## 一、部署架构

```
┌───────────────────────────────────────────────────────────────┐
│                   4 Orderer (Raft 共识)                         │
│  orderer1:7050  orderer2:8050  orderer3:9050  orderer4:10050   │
└──────────────────────────┬────────────────────────────────────┘
                           │
    ┌──────────────────────┼──────────────────────────┐
    │                      │                          │
┌───┴───────────┐  ┌───────┴──────────┐  ┌───────────┴──────────┐
│access-channel │  │billing-channel   │  │maintain-channel      │  ┌──────────────┐
│  access_cc    │  │  billing_cc      │  │  maintain_cc         │  │ esg-channel  │
│  门禁权限      │  │  水电费结算       │  │  设备维护              │  │  esg_cc      │
└───────────────┘  └──────────────────┘  └──────────────────────┘  │  ESG存证     │
                                                                   └──────────────┘
┌───────────────────────────────────────────────────────────────┐
│  peer0.org1:7051                 peer0.org2:9051               │
│  (5个通道已加入)                   (5个通道已加入)               │
└───────────────────────────────────────────────────────────────┘
```

## 二、通道与链码映射

| 通道名 | 链码名 | 业务域 | 说明 |
|--------|--------|--------|------|
| `mychannel` | `basic` | 资产转移 | 原有通道，保留不变 |
| `access-channel` | `access_cc` | 门禁权限管理 | 门禁授权、撤销、出入记录 |
| `billing-channel` | `billing_cc` | 水电费结算 | 账单创建、缴费、抄表记录 |
| `maintain-channel` | `maintain_cc` | 设备维护记录 | 工单创建、派单、完成 |
| `esg-channel` | `esg_cc` | ESG存证 | 文件上传、验证、ESG报告 |

## 三、链码函数清单

### 3.1 access_cc（门禁权限）

| 函数 | 类型 | 参数 | 说明 |
|------|------|------|------|
| `GrantAccess` | 写入 | buildingID, doorID, userID, userName, permissionType, validFrom, validTo | 授权门禁权限 |
| `RevokeAccess` | 写入 | buildingID, doorID, userID | 撤销门禁权限 |
| `QueryPermission` | 查询 | buildingID, doorID, userID | 查询权限状态 |
| `LogAccess` | 写入 | buildingID, doorID, userID, userName, result | 记录通行日志 |
| `QueryAccessLogs` | 查询 | buildingID | 查询楼宇通行记录 |

### 3.2 billing_cc（水电费结算）

| 函数 | 类型 | 参数 | 说明 |
|------|------|------|------|
| `CreateBill` | 写入 | billID, buildingID, unitID, userID, billType, amount, period | 创建账单 |
| `PayBill` | 写入 | billID | 缴费 |
| `QueryBill` | 查询 | billID | 查询账单 |
| `QueryBillsByBuilding` | 查询 | buildingID | 按楼宇查询账单 |
| `RecordMeterReading` | 写入 | readingID, buildingID, unitID, meterType, prevReading, currReading | 抄表记录 |

### 3.3 maintain_cc（设备维护）

| 函数 | 类型 | 参数 | 说明 |
|------|------|------|------|
| `CreateMaintenance` | 写入 | recordID, buildingID, deviceID, deviceName, maintainType, description, reporterID, priority | 创建维护工单 |
| `AssignMaintenance` | 写入 | recordID, handlerID, scheduledDate | 派单 |
| `CompleteMaintenance` | 写入 | recordID | 完成维护 |
| `QueryMaintenance` | 查询 | recordID | 查询工单 |
| `QueryMaintenanceByBuilding` | 查询 | buildingID | 按楼宇查询工单 |

### 3.4 esg_cc（ESG存证）

| 函数 | 类型 | 参数 | 说明 |
|------|------|------|------|
| `UploadESGFile` | 写入 | recordID, buildingID, fileName, fileType, fileSize, ipfsCID, uploaderID, category, description | 上传ESG文件 |
| `VerifyESGFile` | 写入 | recordID, verification | 验证文件 |
| `QueryESGFile` | 查询 | recordID | 查询ESG文件 |
| `QueryESGFilesByBuilding` | 查询 | buildingID | 按楼宇查询ESG文件 |
| `SubmitESGReport` | 写入 | reportID, buildingID, period, carbonEmit, energyUsage, wasteRecycleRate, score | 提交ESG报告 |

## 四、API接口

### 4.1 通用多通道接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/api/multichannel/channels` | GET | 列出所有通道及映射 |
| `/api/multichannel/init` | GET/POST | 重新初始化多通道服务 |
| `/api/multichannel/invoke` | GET/POST | 写入交易（需要背书） |
| `/api/multichannel/query` | GET/POST | 查询数据（仅读） |

### 4.2 invoke/query 参数

- `channel`: 通道名称（如 `access-channel`）
- `fn`: 链码函数名（如 `GrantAccess`）
- `args`: JSON数组字符串（如 `["building_A","door_001","user_001"]`）

### 4.3 调用示例

查询所有通道：
```bash
curl http://localhost:8199/api/multichannel/channels
```

授权门禁（access-channel）：
```bash
curl "http://localhost:8199/api/multichannel/invoke?channel=access-channel&fn=GrantAccess&args=%5B%22building_A%22%2C%22door_001%22%2C%22user_001%22%2C%22Alice%22%2C%22full%22%2C%222026-01-01%22%2C%222027-01-01%22%5D"
```

查询门禁权限：
```bash
curl "http://localhost:8199/api/multichannel/query?channel=access-channel&fn=QueryPermission&args=%5B%22building_A%22%2C%22door_001%22%2C%22user_001%22%5D"
```

创建水电费账单（billing-channel）：
```bash
curl "http://localhost:8199/api/multichannel/invoke?channel=billing-channel&fn=CreateBill&args=%5B%22BILL_001%22%2C%22building_A%22%2C%22unit_101%22%2C%22user_001%22%2C%22water%22%2C%2252.5%22%2C%222026-Q1%22%5D"
```

创建维护工单（maintain-channel）：
```bash
curl "http://localhost:8199/api/multichannel/invoke?channel=maintain-channel&fn=CreateMaintenance&args=%5B%22MAINT_001%22%2C%22building_A%22%2C%22elevator_01%22%2C%22Elevator%22%2C%22repair%22%2C%22Noise%22%2C%22staff_001%22%2C%22high%22%5D"
```

上传ESG文件（esg-channel）：
```bash
curl "http://localhost:8199/api/multichannel/invoke?channel=esg-channel&fn=UploadESGFile&args=%5B%22ESG_001%22%2C%22building_A%22%2C%22report.pdf%22%2C%22pdf%22%2C%222048%22%2C%22QmXabc%22%2C%22admin%22%2C%22carbon%22%2C%22Annual%20report%22%5D"
```

## 五、CLI操作参考

### 环境变量设置

```bash
export PATH=/root/home/go/fabric-samples/bin:/usr/local/go/bin:$PATH
export FABRIC_CFG_PATH=/root/home/go/fabric-samples/test-network/peercfg
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_TLS_ROOTCERT_FILE=/tmp/peer0-org1-tls-ca.crt
export CORE_PEER_LOCALMSPID=Org1MSP
export CORE_PEER_MSPCONFIGPATH=/root/home/go/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export ORDERER_CA=/tmp/orderer1-tls-ca.crt
```

### 链码调用模板

写入交易：
```bash
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer1.example.com \
  --tls --cafile $ORDERER_CA \
  -C <通道名> -n <链码名> \
  --peerAddresses localhost:7051 --tlsRootCertFiles /tmp/peer0-org1-tls-ca.crt \
  --peerAddresses localhost:9051 --tlsRootCertFiles /tmp/peer0-org2-tls-ca.crt \
  -c '{"Args":["<函数名>","参数1","参数2"]}'
```

查询：
```bash
peer chaincode query -C <通道名> -n <链码名> -c '{"Args":["<函数名>","参数1"]}'
```

### 查看已安装/已提交的链码

```bash
# 查看已安装
peer lifecycle chaincode queryinstalled

# 查看通道已提交
peer lifecycle chaincode querycommitted --channelID <通道名>
```

## 六、文件路径索引

| 路径 | 说明 |
|------|------|
| `/root/home/go/fabric-samples/chaincodes/access_cc/` | 门禁权限链码源码 |
| `/root/home/go/fabric-samples/chaincodes/billing_cc/` | 水电费链码源码 |
| `/root/home/go/fabric-samples/chaincodes/maintain_cc/` | 设备维护链码源码 |
| `/root/home/go/fabric-samples/chaincodes/esg_cc/` | ESG存证链码源码 |
| `/tmp/*-channel.block` | 各通道创世块文件 |
| `/tmp/peer0-org1-tls-ca.crt` | peer0.org1 TLS CA证书 |
| `/tmp/peer0-org2-tls-ca.crt` | peer0.org2 TLS CA证书 |
| `/tmp/orderer*-tls-*.crt` | Orderer TLS证书 |
| `fabric-sdk/internal/service/multi_channel.go` | 多通道SDK服务 |
| `fabric-sdk/internal/controller/multi_channel.go` | 多通道控制器 |

## 七、数据结构

### 7.1 门禁权限（AccessPermission）

```json
{
  "docType": "access_permission",
  "permissionId": "PERM_building_A_door_001_user_001",
  "buildingId": "building_A",
  "doorId": "door_001",
  "userId": "user_001",
  "userName": "Alice",
  "permissionType": "full",
  "validFrom": "2026-01-01",
  "validTo": "2027-01-01",
  "status": "active",
  "createdAt": 1777304203
}
```

### 7.2 水电费账单（BillRecord）

```json
{
  "docType": "bill",
  "billId": "BILL_001",
  "buildingId": "building_A",
  "unitId": "unit_101",
  "userId": "user_001",
  "billType": "water",
  "amount": 52.5,
  "period": "2026-Q1",
  "status": "unpaid",
  "paidAt": 0,
  "createdAt": 1777304200
}
```

### 7.3 维护工单（MaintenanceRecord）

```json
{
  "docType": "maintenance",
  "recordId": "MAINT_001",
  "buildingId": "building_A",
  "deviceId": "elevator_01",
  "deviceName": "Elevator",
  "maintainType": "repair",
  "description": "Strange noise from motor",
  "reporterId": "staff_001",
  "handlerId": "",
  "status": "pending",
  "priority": "high",
  "scheduledDate": "",
  "completedAt": 0,
  "createdAt": 1777304200
}
```

### 7.4 ESG文件（ESGRecord）

```json
{
  "docType": "esg_file",
  "recordId": "ESG_001",
  "buildingId": "building_A",
  "fileName": "report_2026.pdf",
  "fileType": "pdf",
  "fileSize": 2048,
  "ipfsCid": "QmXabc123",
  "uploaderId": "admin_001",
  "category": "carbon_report",
  "description": "Annual carbon emission report",
  "verification": "pending",
  "createdAt": 1777304200
}
```

### 7.5 ESG报告（ESGReport）

```json
{
  "docType": "esg_report",
  "reportId": "RPT_001",
  "buildingId": "building_A",
  "period": "2026-Q1",
  "carbonEmit": 125.5,
  "energyUsage": 8500.0,
  "wasteRecycleRate": 0.72,
  "score": 85.5,
  "status": "submitted",
  "createdAt": 1777304200
}
```


## 八、部署信息

- 部署时间：2026-04-27（修复于 2026-04-28）
- Fabric版本：v3.0.0
- 部署状态：全部通道和链码已通过 CLI + SDK 验证
- SDK Gateway 多通道服务已全部通过验证（invoke/query 均正常）

### 修复记录

1. Admin 证书过期：Docker 镜像中证书来自旧 test-network，与当前 Fabric 网络不匹配。修复：将当前证书复制到 configs/fabric 挂载目录，代码改为从挂载路径加载。
2. Anchor Peers 缺失：所有通道均未配置 Anchor Peers，导致跨组织 Gossip 发现失败，Gateway 无法收集双组织背书。修复：为全部 5 个通道的 Org1 和 Org2 设置 Anchor Peers。
3. TLS CA 证书不匹配：configs/fabric 中的 TLS CA 证书与 peer 实际证书不一致。修复：从 peer 容器提取正确的 CA 证书。
