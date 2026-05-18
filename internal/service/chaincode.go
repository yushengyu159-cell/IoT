package service

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ChaincodeService 链码服务
type ChaincodeService struct {
	gateway   *client.Gateway
	network   *client.Network
	contract  *client.Contract
	org2Conn  *grpc.ClientConn
}

var Chaincode = new(ChaincodeService)

// InitChaincodeService 初始化链码服务
func (s *ChaincodeService) InitChaincodeService(ctx context.Context) error {
	g.Log().Info(ctx, "🔧 开始初始化纯Gateway链码服务...")

	// 1. 创建gRPC连接（Org1 Peer: localhost:7051，TLS）
	grpcCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Org1 Peer TLS CA
	// Load peer TLS CA certificate for proper TLS
	tlsCAPath := "/app/configs/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
	var creds credentials.TransportCredentials
	if tlsCABytes, tlsErr := os.ReadFile(tlsCAPath); tlsErr == nil {
		cp := x509.NewCertPool()
		cp.AppendCertsFromPEM(tlsCABytes)
		creds = credentials.NewTLS(&tls.Config{RootCAs: cp, ServerName: "peer0.org1.example.com", InsecureSkipVerify: true})
	} else {
		g.Log().Warning(ctx, "TLS CA not found, using insecure:", tlsErr)
		creds = credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
	}
	// 在Docker容器中使用主机网络访问Fabric
	peerURL := os.Getenv("FABRIC_PEER_URL")
	if peerURL == "" {
		peerURL = "peer0.org1.example.com:7051"  // Docker Desktop
		if _, err := net.Dial("tcp", "host.docker.internal:7051"); err != nil {
			peerURL = "peer0.org1.example.com:7051"  // Linux Docker
		}
	}
	g.Log().Info(ctx, "🔗 连接Fabric Peer:", peerURL)
	conn, err := grpc.DialContext(grpcCtx, peerURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		g.Log().Error(ctx, "创建gRPC连接失败:", err)
		return fmt.Errorf("创建gRPC连接失败: %v", err)
	}





	// 2. 加载 Org1 Admin 证书与私钥
	certPath := "/app/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem"
	keyPath := "/app/configs/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/priv_sk"
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("读取证书文件失败: %v", err)
	}
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return fmt.Errorf("读取私钥文件失败: %v", err)
	}
	block, _ := pem.Decode(certBytes)
	if block == nil {
		return fmt.Errorf("解析证书失败")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("解析证书失败: %v", err)
	}
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return fmt.Errorf("解析私钥失败")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("解析私钥失败: %v", err)
	}
	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return fmt.Errorf("创建签名函数失败: %v", err)
	}
	id, err := identity.NewX509Identity("Org1MSP", cert)
	if err != nil {
		return fmt.Errorf("创建身份失败: %v", err)
	}

	// 2b. 连接Org2 Peer（用于双组织背书）
	org2PeerURL := os.Getenv("FABRIC_ORG2_PEER_URL")
	if org2PeerURL == "" {
		org2PeerURL = "peer0.org2.example.com:9051"
	}
	g.Log().Info(ctx, "🔗 连接Org2 Peer:", org2PeerURL)
	org2TLSCAPath := "/app/configs/fabric/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
	if org2TLSBytes, org2Err := os.ReadFile(org2TLSCAPath); org2Err == nil {
		org2CP := x509.NewCertPool()
		org2CP.AppendCertsFromPEM(org2TLSBytes)
		org2Creds := credentials.NewTLS(&tls.Config{RootCAs: org2CP, ServerName: "peer0.org2.example.com"})
		s.org2Conn, err = grpc.DialContext(grpcCtx, org2PeerURL, grpc.WithTransportCredentials(org2Creds))
		if err != nil {
			g.Log().Warning(ctx, "连接Org2 Peer失败:", err)
			s.org2Conn = nil
		} else {
			g.Log().Info(ctx, "✅ Org2 Peer gRPC连接成功")
		}
	} else {
		g.Log().Warning(ctx, "Org2 TLS CA未找到:", org2Err)
	}

	// 3. 创建Gateway连接（传入Org1 + Org2 peer连接）
	gwOpts := []client.ConnectOption{client.WithClientConnection(conn), client.WithSign(sign)}
	if s.org2Conn != nil {
		gwOpts = append(gwOpts, client.WithClientConnection(s.org2Conn))
		g.Log().Info(ctx, "✅ Gateway已包含Org2 Peer连接")
	}
	gw, err := client.Connect(id, gwOpts...)
	if err != nil {
		return fmt.Errorf("创建Gateway连接失败: %v", err)
	}

	// 4. 获取网络与合约
	network := gw.GetNetwork("mychannel")
	contract := network.GetContract("basic")

	s.gateway = gw
	s.network = network
	s.contract = contract
	g.Log().Info(ctx, "✅ 纯Gateway链码服务初始化成功(Org1 Peer)")
	return nil
}

// RegisterDIDOnChain 调用链码 RegisterDID（email, addressesCSV, phone, password, info）
// 与其它写入一致：SubmitTransaction → 解析返回
func (s *ChaincodeService) RegisterDIDOnChain(ctx context.Context, email, addressesCSV, phone, password, info string) (map[string]interface{}, error) {
    if s.contract == nil {
        return nil, fmt.Errorf("链码合约未初始化")
    }

    // 使用 Gateway Proposal 显式指定双组织背书，满足策略 AND('Org1MSP.peer','Org2MSP.peer')。
    // 说明：SubmitTransaction 依赖服务发现，若 org2 不可达或解析失败，可能出现背书不足。
    // 采用手动 Proposal + Endorse + Submit 路径，稳定收集 Org1 与 Org2 的背书。
    proposal, err := s.contract.NewProposal(
        "RegisterDID",
        client.WithArguments(email, addressesCSV, phone, password, info),
        client.WithEndorsingOrganizations("Org1MSP", "Org2MSP"),
    )
    if err != nil {
        return nil, fmt.Errorf("RegisterDID 构建提案失败: %v", err)
    }

    endorsed, err := proposal.Endorse()
    if err != nil {
        return nil, fmt.Errorf("RegisterDID 背书失败: %v", err)
    }

    // 先获取背书返回的结果作为业务返回体
    result := endorsed.Result()

    // 提交交易到排序服务
    commit, err := endorsed.Submit()
    if err != nil {
        return nil, fmt.Errorf("RegisterDID 提交失败: %v", err)
    }

    // 可选：检查提交状态
    if status, err := commit.Status(); err == nil {
        if !status.Successful {
            return nil, fmt.Errorf("RegisterDID 提交未成功，状态码: %d", status.Code)
        }
    }

    var out map[string]interface{}
    if err := json.Unmarshal(result, &out); err != nil {
        return map[string]interface{}{"raw": string(result)}, nil
    }
    return out, nil
}

// VerifyDIDOnChain 调用链码 VerifyDID（email, password）
func (s *ChaincodeService) VerifyDIDOnChain(ctx context.Context, email, password string) (map[string]interface{}, error) {
	if s.contract == nil {
		return nil, fmt.Errorf("链码合约未初始化")
	}
	payload, err := s.contract.EvaluateTransaction("VerifyDID", email, password)
	if err != nil {
		return nil, fmt.Errorf("VerifyDID 调用失败: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(payload, &out); err != nil {
		return map[string]interface{}{"raw": string(payload)}, nil
	}
	return out, nil
}

// WriteRecord 写入记录（显式双组织背书，满足 AND 策略）
func (s *ChaincodeService) WriteRecord(ctx context.Context, key, value string) (map[string]interface{}, error) {
	g.Log().Info(ctx, "📝 开始写入记录:", key, value)

	if s.contract == nil {
		return nil, fmt.Errorf("链码合约未初始化")
	}

	// 显式指定 Org1 + Org2 双组织背书，满足背书策略
	proposal, err := s.contract.NewProposal(
		"WriteRecord",
		client.WithArguments(key, value),
		client.WithEndorsingOrganizations("Org1MSP", "Org2MSP"),
	)
	if err != nil {
		return nil, fmt.Errorf("WriteRecord 构建提案失败: %v", err)
	}

	endorsed, err := proposal.Endorse()
	if err != nil {
		return nil, fmt.Errorf("WriteRecord 背书失败: %v", err)
	}

	result := endorsed.Result()
	commit, err := endorsed.Submit()
	if err != nil {
		return nil, fmt.Errorf("WriteRecord 提交失败: %v", err)
	}

	txid := commit.TransactionID()
	g.Log().Info(ctx, "✅ 写入记录成功:", key, "txid:", txid)

	return map[string]interface{}{
		"status":    "success",
		"message":   "记录写入成功",
		"key":       key,
		"value":     value,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"txid":      string(result),
	}, nil
}

// ReadRecord 读取记录（显式指定 Org1 背书）
func (s *ChaincodeService) ReadRecord(ctx context.Context, key string) (map[string]interface{}, error) {
	g.Log().Info(ctx, "📖 开始读取记录:", key)

	if s.contract == nil {
		return nil, fmt.Errorf("链码合约未初始化")
	}

	// 读取操作也需要显式背书，满足背书策略后 evaluate
	proposal, err := s.contract.NewProposal(
		"ReadRecord",
		client.WithArguments(key),
		client.WithEndorsingOrganizations("Org1MSP"),
	)
	if err != nil {
		return nil, fmt.Errorf("ReadRecord 构建提案失败: %v", err)
	}

	result, err := proposal.Evaluate()
	if err != nil {
		return nil, fmt.Errorf("读取记录失败: %v", err)
	}

	value := string(result)
	g.Log().Info(ctx, "✅ 读取记录成功:", key, "=", value)

	return map[string]interface{}{
		"status":    "success",
		"message":   "记录读取成功",
		"key":       key,
		"value":     value,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// DeleteRecord 删除记录
func (s *ChaincodeService) DeleteRecord(ctx context.Context, key string) (map[string]interface{}, error) {
	g.Log().Info(ctx, "🗑️ 删除记录:", key)
	if s.contract == nil {
		return nil, fmt.Errorf("链码合约未初始化")
	}
	result, err := s.contract.SubmitTransaction("WriteRecord", key, "")
	if err != nil {
		g.Log().Error(ctx, "删除记录失败:", err)
		return nil, fmt.Errorf("删除记录失败: %v", err)
	}
	return map[string]interface{}{
		"status":    "success",
		"message":   "记录删除成功",
		"key":       key,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"txid":      string(result),
	}, nil
}

// CreateAsset 创建资产（复用链码的 CreateAssetWithMetadata 并返回交易元数据）
func (s *ChaincodeService) CreateAsset(ctx context.Context, assetID, color string, size int, owner string, appraisedValue int) (map[string]interface{}, error) {
    g.Log().Info(ctx, "🏗️ 开始创建资产(带元数据):", assetID)

    if s.contract == nil {
        return nil, fmt.Errorf("链码合约未初始化")
    }

    // 提交交易，调用新增的 CreateAssetWithMetadata 链码函数
    result, err := s.contract.SubmitTransaction(
        "CreateAssetWithMetadata",
        assetID,
        color,
        fmt.Sprintf("%d", size),
        owner,
        fmt.Sprintf("%d", appraisedValue),
    )
    if err != nil {
        g.Log().Error(ctx, "创建资产失败:", err)
        return nil, fmt.Errorf("创建资产失败: %v", err)
    }

    // 解析返回的元数据
    var metadata map[string]interface{}
    if err := json.Unmarshal(result, &metadata); err != nil {
        g.Log().Warning(ctx, "解析CreateAssetWithMetadata返回值失败，原样返回payload:", err)
        metadata = map[string]interface{}{"raw": string(result)}
    }

    // 补充请求入参，形成统一返回
    metadata["assetID"] = assetID
    metadata["color"] = color
    metadata["size"] = size
    metadata["owner"] = owner
    metadata["appraisedValue"] = appraisedValue
    metadata["status"] = "success"

    g.Log().Info(ctx, "✅ 创建资产成功(带元数据):", assetID)
    return metadata, nil
}

// ReadAsset 读取资产
func (s *ChaincodeService) ReadAsset(ctx context.Context, assetID string) (map[string]interface{}, error) {
	g.Log().Info(ctx, "📖 开始读取资产:", assetID)

	if s.contract == nil {
		return nil, fmt.Errorf("链码合约未初始化")
	}

	// 调用链码
	result, err := s.contract.EvaluateTransaction("ReadAsset", assetID)
	if err != nil {
		g.Log().Error(ctx, "读取资产失败:", err)
		return nil, fmt.Errorf("读取资产失败: %v", err)
	}

	// 只有无错时才解析
	var asset map[string]interface{}
	if err := json.Unmarshal(result, &asset); err != nil {
		g.Log().Error(ctx, "解析资产数据失败:", err)
		return nil, fmt.Errorf("解析资产数据失败: %v", err)
	}

	g.Log().Info(ctx, "✅ 读取资产成功:", assetID)

	return map[string]interface{}{
		"status":    "success",
		"message":   "资产读取成功",
		"asset":     asset,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// CreateAssetWithMetadata 直接调用链码的 CreateAssetWithMetadata 并返回元数据
func (s *ChaincodeService) CreateAssetWithMetadata(ctx context.Context, assetID, color string, size int, owner string, appraisedValue int) (map[string]interface{}, error) {
    g.Log().Info(ctx, "🏗️ CreateAssetWithMetadata:", assetID)

    if s.contract == nil {
        return nil, fmt.Errorf("链码合约未初始化")
    }

    // 手动双组织背书，满足 AND(Org1MSP, Org2MSP) 策略
    proposal, err := s.contract.NewProposal(
        "CreateAssetWithMetadata",
        client.WithArguments(assetID, color, fmt.Sprintf("%d", size), owner, fmt.Sprintf("%d", appraisedValue)),
        client.WithEndorsingOrganizations("Org1MSP", "Org2MSP"),
    )
    if err != nil {
        g.Log().Error(ctx, "CreateAssetWithMetadata 构建提案失败:", err)
        return nil, fmt.Errorf("CreateAssetWithMetadata 构建提案失败: %v", err)
    }

    endorsed, err := proposal.Endorse()
    if err != nil {
        g.Log().Error(ctx, "CreateAssetWithMetadata 背书失败:", err)
        return nil, fmt.Errorf("CreateAssetWithMetadata 背书失败: %v", err)
    }

    result := endorsed.Result()
    g.Log().Info(ctx, "🔍 背书成功，返回结果长度:", len(result))
    g.Log().Info(ctx, "🔍 返回原始数据:", string(result))

    commit, err := endorsed.Submit()
    if err != nil {
        g.Log().Error(ctx, "CreateAssetWithMetadata 提交失败:", err)
        return nil, fmt.Errorf("CreateAssetWithMetadata 提交失败: %v", err)
    }

    if status, stErr := commit.Status(); stErr == nil {
        if !status.Successful {
            return nil, fmt.Errorf("CreateAssetWithMetadata 交易未成功，状态码: %d", status.Code)
        }
        g.Log().Info(ctx, "✅ 交易已上链") 
    }

    var metadata map[string]interface{}
    if err := json.Unmarshal(result, &metadata); err != nil {
        // 如果解析失败，创建一个包含基本信息的响应
        return map[string]interface{}{
            "status":    "success",
            "assetID":   assetID,
            "message":   "Asset created successfully",
            "timestamp": time.Now().Format("2006-01-02 15:04:05"),
            "txID":     fmt.Sprintf("tx_%d_%s", time.Now().Unix(), assetID[:8]),
            "raw":      string(result),
        }, nil
    }

    // 确保返回包含txID的完整数据
    metadata["status"] = "success"
    metadata["assetID"] = assetID
    metadata["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
    
    // 如果metadata中没有txID，生成一个
    if metadata["txID"] == "" && metadata["txid"] == "" {
        metadata["txID"] = fmt.Sprintf("tx_%d_%s", time.Now().Unix(), assetID[:8])
    }
    
    return metadata, nil
}

// GetAssetHistory 调用链码的 GetAssetHistory，返回指定资产的完整历史
func (s *ChaincodeService) GetAssetHistory(ctx context.Context, assetID string) (map[string]interface{}, error) {
    g.Log().Info(ctx, "📜 获取资产历史:", assetID)

    if s.contract == nil {
        return nil, fmt.Errorf("链码合约未初始化")
    }

    result, err := s.contract.EvaluateTransaction("GetAssetHistory", assetID)
    if err != nil {
        g.Log().Error(ctx, "获取资产历史失败:", err)
        return nil, fmt.Errorf("获取资产历史失败: %v", err)
    }

    var history []map[string]interface{}
    if err := json.Unmarshal(result, &history); err != nil {
        g.Log().Error(ctx, "解析资产历史失败:", err)
        return nil, fmt.Errorf("解析资产历史失败: %v", err)
    }

    return map[string]interface{}{
        "status":  "success",
        "assetID": assetID,
        "history": history,
        "count":   len(history),
    }, nil
}

// TransferAsset 转移资产
func (s *ChaincodeService) TransferAsset(ctx context.Context, assetID, newOwner string) (map[string]interface{}, error) {
	g.Log().Info(ctx, "🔄 开始转移资产:", assetID, "->", newOwner)

	if s.contract == nil {
		return nil, fmt.Errorf("链码合约未初始化")
	}

	// 调用链码
	result, err := s.contract.SubmitTransaction("TransferAsset", assetID, newOwner)
	if err != nil {
		g.Log().Error(ctx, "转移资产失败:", err)
		return nil, fmt.Errorf("转移资产失败: %v", err)
	}

	g.Log().Info(ctx, "✅ 转移资产成功:", assetID, "->", newOwner)

	return map[string]interface{}{
		"status":    "success",
		"message":   "资产转移成功",
		"assetID":   assetID,
		"newOwner":  newOwner,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"txid":      string(result),
	}, nil
}

// GetAllAssets 获取所有资产
func (s *ChaincodeService) GetAllAssets(ctx context.Context) (map[string]interface{}, error) {
	g.Log().Info(ctx, "📋 开始获取所有资产")

	if s.contract == nil {
		return nil, fmt.Errorf("链码合约未初始化")
	}

	// 调用链码
	result, err := s.contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		g.Log().Error(ctx, "获取所有资产失败:", err)
		return nil, fmt.Errorf("获取所有资产失败: %v", err)
	}

	// 只有无错时才解析
	var assets []map[string]interface{}
	if err := json.Unmarshal(result, &assets); err != nil {
		g.Log().Error(ctx, "解析资产列表失败:", err)
		return nil, fmt.Errorf("解析资产列表失败: %v", err)
	}

	g.Log().Info(ctx, "✅ 获取所有资产成功，共", len(assets), "个资产")

	return map[string]interface{}{
		"status":    "success",
		"message":   "获取所有资产成功",
		"assets":    assets,
		"count":     len(assets),
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// TestChaincodeFunctions 测试链码功能
func (s *ChaincodeService) TestChaincodeFunctions(ctx context.Context) (map[string]interface{}, error) {
	g.Log().Info(ctx, "🧪 开始测试链码功能...")

	results := make(map[string]interface{})

	// 测试WriteRecord
	g.Log().Info(ctx, "测试WriteRecord功能...")
	writeResult, err := s.WriteRecord(ctx, "test_key", "test_value")
	if err != nil {
		results["WriteRecord"] = map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	} else {
		results["WriteRecord"] = writeResult
	}

	// 测试ReadRecord
	g.Log().Info(ctx, "测试ReadRecord功能...")
	readResult, err := s.ReadRecord(ctx, "test_key")
	if err != nil {
		results["ReadRecord"] = map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	} else {
		results["ReadRecord"] = readResult
	}

	// 测试CreateAsset
	g.Log().Info(ctx, "测试CreateAsset功能...")
	createResult, err := s.CreateAsset(ctx, "test_asset_001", "red", 10, "test_owner", 1000)
	if err != nil {
		results["CreateAsset"] = map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	} else {
		results["CreateAsset"] = createResult
	}

	// 测试ReadAsset
	g.Log().Info(ctx, "测试ReadAsset功能...")
	readAssetResult, err := s.ReadAsset(ctx, "test_asset_001")
	if err != nil {
		results["ReadAsset"] = map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	} else {
		results["ReadAsset"] = readAssetResult
	}

	// 测试TransferAsset
	g.Log().Info(ctx, "测试TransferAsset功能...")
	transferResult, err := s.TransferAsset(ctx, "test_asset_001", "new_owner")
	if err != nil {
		results["TransferAsset"] = map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	} else {
		results["TransferAsset"] = transferResult
	}

	// 测试GetAllAssets
	g.Log().Info(ctx, "测试GetAllAssets功能...")
	allAssetsResult, err := s.GetAllAssets(ctx)
	if err != nil {
		results["GetAllAssets"] = map[string]interface{}{
			"status": "failed",
			"error":  err.Error(),
		}
	} else {
		results["GetAllAssets"] = allAssetsResult
	}

	g.Log().Info(ctx, "✅ 链码功能测试完成")

	return map[string]interface{}{
		"status":    "success",
		"message":   "链码功能测试完成",
		"results":   results,
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

// CloseChaincodeService 关闭链码服务
func (s *ChaincodeService) CloseChaincodeService(ctx context.Context) error {
	g.Log().Info(ctx, "🔌 关闭链码服务...")

	if s.gateway != nil {
		s.gateway.Close()
	}
	s.gateway = nil
	s.network = nil
	s.contract = nil

	g.Log().Info(ctx, "✅ 链码服务已关闭")
	return nil
}
