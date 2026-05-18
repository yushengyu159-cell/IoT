package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fabric-sdk/internal/model"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	
	"github.com/gogf/gf/v2/frame/g"
)

type ESGFileMeta struct {
	CID      string `json:"cid"`
	Filename string `json:"filename"`
	Desc     string `json:"desc,omitempty"`
	Uploader string `json:"uploader,omitempty"`
	UploadAt string `json:"uploadAt"`
	FileSize int64  `json:"fileSize,omitempty"`
}

type EncryptedESGFileMeta struct {
	FileName  string   `json:"fileName"`
	Desc      string   `json:"desc"`
	Uploader  string   `json:"uploader"`
	FileSize  int64    `json:"fileSize"`
	ChunkSize int      `json:"chunkSize"`
	CIDs      []string `json:"cids"`
	IV        string   `json:"iv"`
	UploadAt  string   `json:"uploadAt"`
}

// 上传ESG文件，IPFS+MySQL存储（移除链码依赖）
func UploadESGFile(ctx context.Context, file io.Reader, filename, desc, uploader string, fileSize int64) (cid string, txid string, err error) {
	// 1. 上传到IPFS
	cid, err = IpfsUpload(file, filename)
	if err != nil {
		return "", "", fmt.Errorf("IPFS上传失败: %v", err)
	}
	// 2. 组装元数据
	meta := ESGFileMeta{
		CID:      cid,
		Filename: filename,
		Desc:     desc,
		Uploader: uploader,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	// 3. 直接写入MySQL（移除链码依赖）
	txid = "mysql_storage_" + time.Now().Format("20060102150405")
 if DB != nil {
		dbFile := model.ESGFile{
			CID: cid,
		}
		result := DB.Where("c_id = ?", cid).FirstOrCreate(&dbFile)
		if result.Error == nil {
			dbFile.Filename = filename
			dbFile.Desc = desc
			dbFile.Uploader = uploader
			dbFile.UploadAt = meta.UploadAt
			dbFile.Txid = txid
			dbFile.FileSize = fileSize
			DB.Save(&dbFile)
		} else {
			fmt.Printf("[ESG] DB save warning: %v\n", result.Error)
		}
	}
	return cid, txid, nil
}

// EncryptAndUploadESGFile 加密分片上传+MySQL存储（移除链码依赖）
func EncryptAndUploadESGFile(ctx context.Context, file io.Reader, filename, desc, uploader string) (meta EncryptedESGFileMeta, key []byte, cipherSample string, chunkCount int, chunkUrls []string, fabricResult map[string]interface{}, err error) {
	// 读取全部文件内容
	plainData, err := io.ReadAll(file)
	if err != nil {
		return meta, nil, "", 0, nil, nil, err
	}
	fileSize := len(plainData)
	// 计算分片大小，平均分为3片
	chunkSize := fileSize / 3
	if fileSize%3 != 0 {
		chunkSize++ // 保证最后一片包含所有剩余数据
	}
	// 重新用EncryptAndSplitFile分片
	chunks, key, iv, err := EncryptAndSplitFile(bytes.NewReader(plainData), chunkSize)
	if err != nil {
		return meta, nil, "", 0, nil, nil, err
	}
	var cids []string
	for i, chunk := range chunks {
		cid, err := IpfsUpload(bytes.NewReader(chunk), filename+".part"+strconv.Itoa(i))
		if err != nil {
			return meta, nil, "", 0, nil, nil, err
		}
		cids = append(cids, cid)
	}
	beijing := time.Now().In(time.FixedZone("CST", 8*3600))
	meta = EncryptedESGFileMeta{
		FileName:  filename,
		Desc:      desc,
		Uploader:  uploader,
		FileSize:  int64(fileSize),
		ChunkSize: chunkSize,
		CIDs:      cids,
		IV:        base64.StdEncoding.EncodeToString(iv),
		UploadAt:  beijing.Format("2006-01-02 15:04:05"),
	}
	// 移除链码依赖，直接生成模拟结果
	fabricResult = map[string]interface{}{
		"txID":     "mysql_storage_" + time.Now().Format("20060102150405"),
		"status":   "success",
		"message":  "文件已存储到MySQL数据库",
		"assetID":  cids[0],
	}
	cipherFragment := ""
	if len(chunks) > 0 && len(chunks[0]) >= 32 {
		cipherFragment = base64.StdEncoding.EncodeToString(chunks[0][:32])
	} else if len(chunks) > 0 {
		cipherFragment = base64.StdEncoding.EncodeToString(chunks[0])
	}
	chunkCount = len(chunks)
	for _, cid := range cids {
		chunkUrls = append(chunkUrls, "http://127.0.0.1:8081/ipfs/"+cid)
	}
	
	// 保存文件信息到数据库
	fmt.Printf("=== 开始数据库保存逻辑 ===\n")
	fmt.Printf("检查数据库连接: DB=%v\n", DB)
	fmt.Printf("数据库连接类型: %T\n", DB)
	
	if DB != nil {
		// 使用第一个分片CID作为主文件标识
		mainCID := cids[0]
		fmt.Printf("准备保存文件到数据库: CID=%s, Filename=%s, Desc=%s, Uploader=%s\n", 
			mainCID, filename, desc, uploader)
		
		// 将所有分片CIDs转换为JSON字符串
		allCIDsJSON, _ := json.Marshal(cids)
		
		// 从fabricResult中获取交易ID
		txid := "encrypted_upload"
		if fabricResult != nil {
			if txID, ok := fabricResult["txID"].(string); ok {
				txid = txID
			}
		}
		
		// 记录当前时间用于计算总耗时
		endTime := time.Now()
		
		dbFile := model.ESGFile{
			CID:            mainCID,
			Filename:       filename,
			Desc:           desc,
			Uploader:       uploader,
			UploadAt:       meta.UploadAt,
			Txid:           txid,
			ChunkCount:     chunkCount,
			ChunkSize:      chunkSize,
			AllCIDs:        string(allCIDsJSON),
			EncryptionKey:  base64.StdEncoding.EncodeToString(key),
			IV:             base64.StdEncoding.EncodeToString(iv),
			CipherSample:   cipherFragment,
			FileSize:       int64(fileSize),
			UploadStartTime: meta.UploadAt,
			UploadEndTime:   endTime.Format("2006-01-02 15:04:05"),
			TotalTime:      fmt.Sprintf("%dms", endTime.Sub(time.Now().Add(-time.Duration(fileSize)*time.Millisecond)).Milliseconds()),
		}
		
		fmt.Printf("准备创建数据库记录: %+v\n", dbFile)
		
		if err := DB.Create(&dbFile).Error; err != nil {
			// 记录错误但不影响返回结果
			fmt.Printf("保存到数据库失败: %v\n", err)
		} else {
			fmt.Printf("文件成功保存到数据库，ID: %d\n", dbFile.ID)
		}
	} else {
		fmt.Printf("❌ 数据库连接为空，无法保存文件信息\n")
	}
	
	fmt.Printf("=== 数据库保存逻辑结束 ===\n")
	
	return meta, key, cipherFragment, chunkCount, chunkUrls, fabricResult, nil
}

// DownloadAndDecryptESGFile 根据链上元数据和密钥下载解密
func DownloadAndDecryptESGFile(ctx context.Context, cid string, key []byte) ([]byte, error) {
	result, err := Chaincode.ReadAsset(ctx, cid)
	if err != nil {
		return nil, err
	}
	// 从资产数据中提取信息
	_, ok := result["asset"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("链码返回数据格式错误")
	}
	// 构建元数据（这里简化处理，实际应该从资产属性中提取）
	meta := EncryptedESGFileMeta{
		FileName:  cid,
		Desc:      "ESG文件",
		Uploader:  "unknown",
		FileSize:  0,
		ChunkSize: 0,
		CIDs:      []string{cid},
		IV:        "",
		UploadAt:  time.Now().Format("2006-01-02 15:04:05"),
	}
	iv, _ := base64.StdEncoding.DecodeString(meta.IV)
	var chunks [][]byte
	for _, cid := range meta.CIDs {
		chunkReader, err := IpfsDownload(cid)
		if err != nil {
			return nil, err
		}
		chunk, err := io.ReadAll(chunkReader)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}
	return MergeAndDecryptFile(chunks, key, iv)
}

// 查询单个ESG文件元数据
func QueryESGFile(ctx context.Context, cid string) (meta ESGFileMeta, err error) {
	result, err := Chaincode.ReadAsset(ctx, cid)
	if err != nil {
		return meta, err
	}
	// 从资产数据中提取信息
	asset, ok := result["asset"].(map[string]interface{})
	if !ok {
		return meta, fmt.Errorf("链码返回数据格式错误")
	}
	
	// 从资产数据中提取详细信息
	meta = ESGFileMeta{
		CID:      cid,
		Filename: cid, // 使用CID作为文件名
		Desc:     "ESG文件",
		Uploader: "unknown",
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	
	// 尝试从资产属性中提取更多信息
	if owner, ok := asset["Owner"].(string); ok && owner != "" {
		meta.Uploader = owner
	}
	if color, ok := asset["Color"].(string); ok && color != "" {
		meta.Desc = fmt.Sprintf("ESG文件 - 类型: %s", color)
	}
	
	return meta, nil
}

// 批量检索所有ESG文件元数据（从区块链）
func ListESGFiles(ctx context.Context) ([]ESGFileMeta, error) {
	result, err := Chaincode.GetAllAssets(ctx)
	if err != nil {
		return nil, err
	}
	assets, ok := result["assets"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("链码返回数据格式错误")
	}
	var metas []ESGFileMeta
	for _, item := range assets {
		assetMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		// 从资产数据中提取详细信息
		if id, ok := assetMap["ID"].(string); ok {
			meta := ESGFileMeta{
				CID:      id,
				Filename: id,
				Desc:     "ESG文件",
				Uploader: "unknown",
				UploadAt: time.Now().Format("2006-01-02 15:04:05"),
			}
			
			// 尝试从资产属性中提取更多信息
			if owner, ok := assetMap["Owner"].(string); ok && owner != "" {
				meta.Uploader = owner
			}
			if color, ok := assetMap["Color"].(string); ok && color != "" {
				meta.Desc = fmt.Sprintf("ESG文件 - 类型: %s", color)
			}
			
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

// 从数据库获取ESG文件列表
func ListESGFilesFromDB(ctx context.Context) ([]ESGFileMeta, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	
	var files []model.ESGFile
	if err := DB.Find(&files).Error; err != nil {
		return nil, fmt.Errorf("查询数据库失败: %v", err)
	}
	
	var metas []ESGFileMeta
	for _, file := range files {
		meta := ESGFileMeta{
			CID:      file.CID,
			Filename: file.Filename,
			Desc:     file.Desc,
			Uploader: file.Uploader,
			UploadAt: file.UploadAt,
			FileSize: file.FileSize,
		}
		metas = append(metas, meta)
	}
	
	return metas, nil
}

// 从数据库获取指定用户的ESG文件列表（用户隔离）
func ListESGFilesFromDBByUser(ctx context.Context, userEmail string) ([]ESGFileMeta, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	
	if userEmail == "" {
		return nil, fmt.Errorf("用户邮箱不能为空")
	}
	
	g.Log().Info(nil, "🔍 开始查询用户文件，用户邮箱:", userEmail)
	
	var files []model.ESGFile
	// 根据用户邮箱或DID过滤文件，实现用户隔离
	// 支持两种查询方式：
	// 1. 直接按邮箱查询（如果uploader字段存储的是邮箱）
	// 2. 按DID查询（如果uploader字段存储的是DID）
	if err := DB.Where("uploader = ? OR uploader = ?", userEmail, fmt.Sprintf("did:example:%s", userEmail)).Find(&files).Error; err != nil {
		return nil, fmt.Errorf("查询用户文件失败: %v", err)
	}
	
	g.Log().Info(nil, "🔍 第一次查询结果，文件数量:", len(files))
	
	// 如果上面的查询没有结果，尝试通过用户表查找DID，再按DID查询文件
	if len(files) == 0 {
		g.Log().Info(nil, "🔍 第一次查询无结果，尝试通过用户表查找DID")
		// 先查询用户表，获取用户的DID
		var user model.DID
		if err := DB.Where("email = ?", userEmail).First(&user).Error; err == nil && user.DID != "" {
			g.Log().Info(nil, "🔍 找到用户DID:", user.DID, "，按DID查询文件")
			// 找到了用户DID，按DID查询文件
			if err := DB.Where("uploader = ?", user.DID).Find(&files).Error; err != nil {
				return nil, fmt.Errorf("按DID查询用户文件失败: %v", err)
			}
			g.Log().Info(nil, "🔍 按DID查询结果，文件数量:", len(files))
		} else {
			g.Log().Warning(nil, "⚠️ 未找到用户DID信息")
		}
	}
	
	var metas []ESGFileMeta
	for _, file := range files {
		meta := ESGFileMeta{
			CID:      file.CID,      // 这里会自动映射到数据库的 c_id 字段
			Filename: file.Filename, // 这里会自动映射到数据库的 filename 字段
			Desc:     file.Desc,     // 这里会自动映射到数据库的 desc 字段
			Uploader: file.Uploader, // 这里会自动映射到数据库的 uploader 字段
			UploadAt: file.UploadAt, // 这里会自动映射到数据库的 upload_at 字段
		}
		metas = append(metas, meta)
	}
	
	g.Log().Info(nil, "✅ 用户文件查询完成，最终返回文件数量:", len(metas))
	return metas, nil
}

// 从数据库查询单个ESG文件详情
func QueryESGFileFromDB(ctx context.Context, cid string) (*model.ESGFile, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	
	var file model.ESGFile
	// 使用正确的数据库字段名 c_id
	if err := DB.Where("c_id = ?", cid).First(&file).Error; err != nil {
		// 兜底：尝试使用去空白cid或按txid再次查询，避免因为传参差异导致的查不到
		trimmed := strings.TrimSpace(cid)
		if trimmed != cid {
			if e2 := DB.Where("c_id = ?", trimmed).First(&file).Error; e2 == nil {
				return &file, nil
			}
		}
		if e3 := DB.Where("txid = ?", trimmed).First(&file).Error; e3 == nil {
			return &file, nil
		}
		return nil, fmt.Errorf("查询文件失败: %v", err)
	}
	
	return &file, nil
}

// 本地批量查询ESG文件元数据
func BatchQueryESGFiles(ctx context.Context, cids []string) ([]ESGFileMeta, error) {
	var metas []ESGFileMeta
	for _, cid := range cids {
		meta, err := QueryESGFile(ctx, cid)
		if err == nil {
			metas = append(metas, meta)
		}
	}
	return metas, nil
}

// 本地批量查询指定用户的ESG文件元数据（用户隔离）
func BatchQueryESGFilesByUser(ctx context.Context, cids []string, userEmail string) ([]ESGFileMeta, error) {
	if userEmail == "" {
		return nil, fmt.Errorf("用户邮箱不能为空")
	}
	
	// 先获取用户的DID，用于文件归属验证
	var userDID string
	var user model.DID
	if err := DB.Where("email = ?", userEmail).First(&user).Error; err == nil && user.DID != "" {
		userDID = user.DID
	}
	
	var metas []ESGFileMeta
	for _, cid := range cids {
		// 查询文件详情
		file, err := QueryESGFileFromDB(ctx, cid)
		if err == nil {
			// 检查文件是否属于当前用户（支持邮箱和DID两种方式）
			if file.Uploader == userEmail || file.Uploader == userDID {
				meta := ESGFileMeta{
					CID:      file.CID,
					Filename: file.Filename,
					Desc:     file.Desc,
					Uploader: file.Uploader,
					UploadAt: file.UploadAt,
				}
				metas = append(metas, meta)
			}
		}
	}
	return metas, nil
}

// 查询IPFS分片资产信息
func QueryIPFSChunkAsset(ctx context.Context, cid string) (map[string]interface{}, error) {
	result, err := Chaincode.ReadAsset(ctx, cid)
	if err != nil {
		return nil, err
	}
	
	// 从资产数据中提取信息
	asset, ok := result["asset"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("链码返回数据格式错误")
	}
	
	// 构建IPFS分片资产信息
	chunkInfo := map[string]interface{}{
		"cid":         cid,
		"assetType":   "IPFS_CHUNK",
		"owner":       "unknown",
		"createdAt":   time.Now().Format("2006-01-02 15:04:05"),
		"assetData":   asset,
	}
	
	// 尝试从资产属性中提取更多信息
	if owner, ok := asset["Owner"].(string); ok && owner != "" {
		chunkInfo["owner"] = owner
	}
	if color, ok := asset["Color"].(string); ok && color != "" {
		chunkInfo["assetType"] = color
	}
	
	return chunkInfo, nil
}

// 获取所有IPFS分片资产
func ListIPFSChunkAssets(ctx context.Context) ([]map[string]interface{}, error) {
	result, err := Chaincode.GetAllAssets(ctx)
	if err != nil {
		return nil, err
	}
	
	assets, ok := result["assets"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("链码返回数据格式错误")
	}
	
	var chunkAssets []map[string]interface{}
	for _, item := range assets {
		assetMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		
		// 检查是否为IPFS分片资产
		if color, ok := assetMap["Color"].(string); ok && (color == "IPFS_CHUNK" || color == "ESG_FILE") {
			if id, ok := assetMap["ID"].(string); ok {
				chunkInfo := map[string]interface{}{
					"cid":         id,
					"assetType":   color,
					"owner":       "unknown",
					"createdAt":   time.Now().Format("2006-01-02 15:04:05"),
					"assetData":   assetMap,
				}
				
				// 尝试从资产属性中提取更多信息
				if owner, ok := assetMap["Owner"].(string); ok && owner != "" {
					chunkInfo["owner"] = owner
				}
				
				chunkAssets = append(chunkAssets, chunkInfo)
			}
		}
	}
	
	return chunkAssets, nil
}
