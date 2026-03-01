package controller

import (
	"bytes"
	"encoding/base64"
	"fabric-sdk/internal/service"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

type IpfsController struct{}

var Ipfs = new(IpfsController)

// POST /api/ipfs/upload
func (c *IpfsController) Upload(r *ghttp.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("[IPFS Upload] panic: %v", rec)
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "IPFS上传接口异常: panic"})
		}
	}()
	file := r.GetUploadFile("file")
	if file == nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "未上传文件"})
		return
	}
	f, err := file.Open()
	if err != nil {
		log.Printf("[IPFS Upload] 文件打开失败: %v", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "文件打开失败: " + err.Error()})
		return
	}
	defer f.Close()
	cid, err := service.IpfsUpload(f, file.Filename)
	if err != nil {
		log.Printf("[IPFS Upload] IPFS上传失败: %v", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "IPFS上传失败: " + err.Error()})
		return
	}
	if cid == "" {
		log.Printf("[IPFS Upload] IPFS返回空CID")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "IPFS返回空CID"})
		return
	}
	log.Printf("[IPFS Upload] 上传成功, CID: %s", cid)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "上传成功", Data: map[string]interface{}{"cid": cid}})
}

// GET /api/ipfs/download?cid=xxx
func (c *IpfsController) Download(r *ghttp.Request) {
	cid := r.Get("cid").String()
	if cid == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "cid参数不能为空"})
		return
	}
	reader, err := service.IpfsDownload(cid)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "IPFS下载失败: " + err.Error()})
		return
	}
	defer reader.Close()
	// 直接流式返回文件内容
	r.Response.Header().Set("Content-Disposition", "attachment; filename="+cid)
	r.Response.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(r.Response.Writer, reader)
}

// POST /api/ipfs/upload-encrypted
func (c *IpfsController) UploadEncrypted(r *ghttp.Request) {
	startTime := time.Now()

	file := r.GetUploadFile("file")
	if file == nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "未上传文件"})
		return
	}
	f, err := file.Open()
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "文件打开失败: " + err.Error()})
		return
	}
	defer f.Close()

	// 读取文件时间统计
	readStart := time.Now()
	plainData, err := io.ReadAll(f)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "文件读取失败: " + err.Error()})
		return
	}
	readTime := time.Since(readStart)

	fileSize := len(plainData)
	chunkSize := fileSize / 4  // 改为4个分片
	if fileSize%4 != 0 {
		chunkSize++
	}

	// 加密分片时间统计
	encryptStart := time.Now()
	chunks, key, iv, err := service.EncryptAndSplitFile(bytes.NewReader(plainData), chunkSize)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "加密分片失败: " + err.Error()})
		return
	}
	encryptTime := time.Since(encryptStart)

	// IPFS上传时间统计
	uploadStart := time.Now()
	var cids []string
	var chunkUrls []string
	for i, chunk := range chunks {
		cid, err := service.IpfsUpload(bytes.NewReader(chunk), file.Filename+".part"+strconv.Itoa(i))
		if err != nil {
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "IPFS分片上传失败: " + err.Error()})
			return
		}
		cids = append(cids, cid)
		chunkUrls = append(chunkUrls, "http://127.0.0.1:8081/ipfs/"+cid)
	}
	uploadTime := time.Since(uploadStart)

	// 区块链写入时间统计

	// 区块链写入时间统计
	blockchainStart := time.Now()
	var fabricResult map[string]interface{}
	var txHash string
	onChain := false
	if len(cids) > 0 {
		// 使用第一个分片CID作为资产ID，创建IPFS分片资产
		fabricResult, err = service.Chaincode.CreateAssetWithMetadata(r.Context(), cids[0], "IPFS_CHUNK", 0, "uploader", 0)
		if err == nil {
			onChain = true
			// 从链码返回结果中提取交易哈希
			if txID, ok := fabricResult["txID"]; ok {
				txHash = txID.(string)
				fmt.Printf("✅ 区块链写入成功，交易哈希: %s\n", txHash)
			} else if txid, ok := fabricResult["txid"]; ok {
				txHash = txid.(string)
				fmt.Printf("✅ 区块链写入成功，交易哈希: %s\n", txHash)
			} else {
				txHash = "unknown"
				fmt.Printf("⚠️ 区块链写入成功，但未获取到交易哈希\n")
			}
		} else {
			// 记录错误但不中断流程
			fmt.Printf("❌ 区块链写入失败: %v\n", err)
			fmt.Printf("🔍 错误详情: %+v\n", err)
		}
	}
	blockchainTime := time.Since(blockchainStart)

	cipherSample := ""
	if len(chunks) > 0 && len(chunks[0]) >= 32 {
		cipherSample = base64.StdEncoding.EncodeToString(chunks[0][:32])
	} else if len(chunks) > 0 {
		cipherSample = base64.StdEncoding.EncodeToString(chunks[0])
	}

	// 计算总耗时
	totalTime := time.Since(startTime)

	// 时间统计信息
	timeStats := map[string]interface{}{
		"totalTime":      totalTime.String(),
		"totalTimeMs":    totalTime.Milliseconds(),
		"readTime":       readTime.String(),
		"readTimeMs":     readTime.Milliseconds(),
		"encryptTime":    encryptTime.String(),
		"encryptTimeMs":  encryptTime.Milliseconds(),
		"uploadTime":     uploadTime.String(),
		"uploadTimeMs":   uploadTime.Milliseconds(),
		"blockchainTime": blockchainTime.String(),
		"blockchainTimeMs": blockchainTime.Milliseconds(),
		"startTime":      startTime.Format("2006-01-02 15:04:05.000"),
		"endTime":        time.Now().Format("2006-01-02 15:04:05.000"),
	}

	resp := map[string]interface{}{
		"cids":         cids,
		"key":          base64.StdEncoding.EncodeToString(key),
		"iv":           base64.StdEncoding.EncodeToString(iv),
		"chunkSize":    chunkSize,
		"fileName":     file.Filename,
		"onChain":      onChain,
		"txHash":       txHash,           // 添加区块链交易哈希
		"blockchainStatus": map[string]interface{}{
			"success":    onChain,
			"txHash":     txHash,
			"assetID":    func() string { 
				if len(cids) > 0 { 
					return cids[0] 
				} 
				return "" 
			}(),
			"timestamp":  time.Now().Format("2006-01-02 15:04:05"),
		},
		"cipherSample": cipherSample,
		"chunkCount":   len(chunks),
		"chunkUrls":    chunkUrls,
		"timeStats":    timeStats,
	}

	if onChain {
		resp["fabricResult"] = fabricResult
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "加密分片上传并链上存证成功", Data: resp})
	} else {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "加密分片上传成功（链上存证失败）", Data: resp})
	}
}

// POST /api/ipfs/download-decrypted
func (c *IpfsController) DownloadDecrypted(r *ghttp.Request) {
	var req struct {
		Cids      []string `json:"cids"`
		KeyBase64 string   `json:"key"`
		IVBase64  string   `json:"iv"`
	}
	if err := r.Parse(&req); err != nil || len(req.Cids) == 0 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "参数错误"})
		return
	}
	key, err := base64.StdEncoding.DecodeString(req.KeyBase64)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "密钥解码失败"})
		return
	}
	iv, err := base64.StdEncoding.DecodeString(req.IVBase64)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "IV解码失败"})
		return
	}
	var chunks [][]byte
	for i, cid := range req.Cids {
		log.Printf("[Download] 正在下载分片 %d/%d: %s", i+1, len(req.Cids), cid)
		chunkReader, err := service.IpfsDownload(cid)
		if err != nil {
			log.Printf("[Download] 分片 %s 下载失败: %v", cid, err)
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "IPFS分片下载失败: " + err.Error()})
			return
		}
		chunk, err := io.ReadAll(chunkReader)
		if err != nil {
			log.Printf("[Download] 分片 %s 读取失败: %v", cid, err)
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "分片读取失败: " + err.Error()})
			return
		}
		log.Printf("[Download] 分片 %s 下载成功，大小: %d 字节", cid, len(chunk))
		chunks = append(chunks, chunk)
	}
	log.Printf("[Download] 所有分片下载完成，总数: %d", len(chunks))
	data, err := service.MergeAndDecryptFile(chunks, key, iv)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "解密失败: " + err.Error()})
		return
	}
	// 直接返回文件内容
	r.Response.Write(data)
}
