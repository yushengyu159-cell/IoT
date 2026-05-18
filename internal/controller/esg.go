package controller

import (
	"context"
	"strings"
	"fabric-sdk/internal/middleware"
	"fabric-sdk/internal/service"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

type ESGController struct{}

var ESG = new(ESGController)
var DB *gorm.DB

// GenerateReport simulates AI report generation
// POST /api/esg/report/generate
func (c *ESGController) GenerateReport(r *ghttp.Request) {
	// Simulate AI processing delay in frontend via progress bar, 
	// here we just return the success signal
	reportType := r.Get("type").String()
	period := r.Get("period").String()
	
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Report Generation Started",
		Data: g.Map{
			"reportId":   "RPT-" + time.Now().Format("20060102150405"),
			"status":     "processing",
			"type":       reportType,
			"period":     period,
			"estimated":  "30s",
		},
	})
}

// GetAIAnalysis returns AI analysis results
// GET /api/esg/analysis
func (c *ESGController) GetAIAnalysis(r *ghttp.Request) {
	// Mock AI Analysis Data
	data := g.Map{
		"rating": g.Map{
			"overall": "AA",
			"score":   88,
			"trend":   "stable",
			"dimensions": []g.Map{
				{"name": "Environmental", "score": 92, "grade": "AAA"},
				{"name": "Social", "score": 85, "grade": "AA"},
				{"name": "Governance", "score": 87, "grade": "AA"},
			},
		},
		"prediction": g.Map{
			"nextQuarter": "AA+",
			"confidence":  "92%",
			"trend": []g.Map{
				{"date": "2024-Q1", "score": 85},
				{"date": "2024-Q2", "score": 86},
				{"date": "2024-Q3", "score": 88},
				{"date": "2024-Q4", "score": 89}, // Prediction
			},
		},
		"comparison": g.Map{
			"labels": []string{"Energy", "Water", "Waste", "Community", "Compliance"},
			"myBuilding": []float64{90, 85, 88, 82, 95},
			"industryAvg": []float64{75, 70, 72, 75, 80},
		},
		"suggestions": []g.Map{
			{
				"title": "Upgrade HVAC System",
				"impact": "High",
				"category": "Environmental",
				"description": "Replace aging chiller units to improve energy efficiency by 15%.",
			},
			{
				"title": "Community Engagement Program",
				"impact": "Medium",
				"category": "Social",
				"description": "Launch quarterly community workshops to boost social score.",
			},
		},
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    data,
	})
}

// GetReportHistory returns list of generated reports
// GET /api/esg/reports
func (c *ESGController) GetReportHistory(r *ghttp.Request) {
	reports := []g.Map{
		{
			"id": "RPT-20241201001",
			"name": "2024 Annual ESG Report",
			"type": "Annual",
			"date": "2024-12-01",
			"status": "Completed",
			"downloadUrl": "#",
		},
		{
			"id": "RPT-20241001002",
			"name": "Q3 2024 Sustainability Brief",
			"type": "Quarterly",
			"date": "2024-10-01",
			"status": "Completed",
			"downloadUrl": "#",
		},
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    reports,
	})
}

// POST /api/esg/upload
func (c *ESGController) UploadFile(r *ghttp.Request) {
	g.Log().Infof(r.Context(), "===== [ESG UPLOAD] 收到上传请求 =====")
	file := r.GetUploadFile("file")
	desc := r.Get("desc").String()
	// Get uploader from token (trusted), fallback to form param
	uploader := ""
	if email, valid := middleware.ValidateToken(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")); valid {
		uploader = email
	}
	if uploader == "" {
		uploader = r.Get("uploader").String()
	}
	g.Log().Infof(r.Context(), "[ESG UPLOAD] file=%v desc=%s uploader=%s", file, desc, uploader)
	if file == nil {
		g.Log().Warningf(r.Context(), "[ESG UPLOAD] file is nil!")
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "未上传文件"})
		return
	}
	g.Log().Infof(r.Context(), "[ESG UPLOAD] filename=%s size=%d", file.Filename, file.Size)
	f, err := file.Open()
	if err != nil {
		g.Log().Errorf(r.Context(), "[ESG UPLOAD] 文件打开失败: %v", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "文件打开失败: " + err.Error()})
		return
	}
	defer f.Close()
	cid, _, err := service.UploadESGFile(context.Background(), f, file.Filename, desc, uploader, file.Size)
	if err != nil {
		g.Log().Errorf(r.Context(), "[ESG UPLOAD] 上传失败: %v", err)
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: err.Error()})
		return
	}
	g.Log().Infof(r.Context(), "[ESG UPLOAD] 上传成功 CID=%s", cid)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "文件已成功存储到区块链",
		Data: map[string]interface{}{
			"cid":     cid,
			"onChain": true,
		},
	})
}

// GET /api/esg/query?cid=xxx
func (c *ESGController) QueryFile(r *ghttp.Request) {
	cid := r.Get("cid").String()
	if cid == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "cid参数不能为空"})
		return
	}
	file, err := service.QueryESGFileFromDB(r.Context(), cid)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: err.Error()})
		return
	}
	
	// 增强文件详情：补充上传者的 DID 和姓名信息
	enhancedFile := map[string]interface{}{
		"ID":              file.ID,
		"CID":             file.CID,
		"Filename":        file.Filename,
		"Desc":            file.Desc,
		"Uploader":        file.Uploader,
		"UploadAt":        file.UploadAt,
		"Txid":            file.Txid,
		"ChunkCount":      file.ChunkCount,
		"ChunkSize":       file.ChunkSize,
		"AllCIDs":         file.AllCIDs,
		"EncryptionKey":   file.EncryptionKey,
		"IV":              file.IV,
		"CipherSample":    file.CipherSample,
		"FileSize":        file.FileSize,
		"UploadStartTime": file.UploadStartTime,
		"UploadEndTime":   file.UploadEndTime,
		"TotalTime":       file.TotalTime,
	}

	// 直接从数据库关联查询上传者的 DID 和姓名信息
	if file.Uploader != "" {
		g.Log().Infof(r.Context(), "🔍 开始查询上传者信息: %s", file.Uploader)
		
		uploaderInfo, err := service.GetUploaderInfoFromDB(file.Uploader)
		if err == nil && uploaderInfo != nil {
			g.Log().Infof(r.Context(), "✅ 成功获取上传者信息: %+v", uploaderInfo)
			enhancedFile["uploaderDid"] = uploaderInfo.DID
			enhancedFile["uploaderName"] = uploaderInfo.FullName
			enhancedFile["uploaderEmail"] = uploaderInfo.Email
			enhancedFile["uploaderRole"] = uploaderInfo.Role
		} else {
			// 如果查询失败，记录日志但不影响主要功能
			g.Log().Warningf(r.Context(), "❌ 无法获取上传者信息: %s, 错误: %v", file.Uploader, err)
			
			// 尝试直接查询数据库，看看是否能找到用户信息
			if service.DB != nil {
				var didRecord struct {
					DID      string `gorm:"column:did"`
					FullName string `gorm:"column:full_name"`
					Email    string `gorm:"column:email"`
					Role     string `gorm:"column:role"`
				}
				// 尝试按 DID 查询
				if err := service.DB.Table("dids").Where("did = ?", file.Uploader).First(&didRecord).Error; err == nil {
					g.Log().Infof(r.Context(), "✅ 通过 DID 直接查询成功: %+v", didRecord)
					enhancedFile["uploaderDid"] = didRecord.DID
					enhancedFile["uploaderName"] = didRecord.FullName
					enhancedFile["uploaderEmail"] = didRecord.Email
					enhancedFile["uploaderRole"] = didRecord.Role
				} else {
					g.Log().Warningf(r.Context(), "❌ 通过 DID 直接查询也失败: %v", err)
				}
			}
		}
	} else {
		g.Log().Warningf(r.Context(), "⚠️ 文件没有上传者信息: %s", file.CID)
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "查询成功", Data: enhancedFile})
}

// GET /api/esg/list
func (c *ESGController) ListFiles(r *ghttp.Request) {
	userEmail := r.Get("userEmail").String()
	var metas ([]service.ESGFileMeta)
	var err error
	if userEmail == "" || userEmail == "all" {
		metas, err = service.ListESGFilesFromDB(r.Context())
	} else {
		metas, err = service.ListESGFilesFromDBByUser(r.Context(), userEmail)
	}
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: err.Error()})
		return
	}
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "批量查询成功", Data: metas})
}
// POST /api/esg/batch-list
func (c *ESGController) BatchListFiles(r *ghttp.Request) {
	var req struct {
		Cids       []string `json:"cids"`
		UserEmail  string   `json:"userEmail"`  // 添加用户邮箱参数
	}
	if err := r.Parse(&req); err != nil || len(req.Cids) == 0 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "参数cids不能为空"})
		return
	}
	if req.UserEmail == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "缺少用户邮箱参数"})
		return
	}
	
	metas, _ := service.BatchQueryESGFilesByUser(r.Context(), req.Cids, req.UserEmail)
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "本地批量查询成功", Data: metas})
}

// POST /api/esg/upload-encrypted
func (c *ESGController) UploadEncrypted(r *ghttp.Request) {
	startTime := time.Now()

	file := r.GetUploadFile("file")
	desc := r.Get("desc").String()
	uploader := r.Get("uploader").String()
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

	// 调用服务层处理加密上传
	meta, key, cipherSample, chunkCount, chunkUrls, fabricResult, err := service.EncryptAndUploadESGFile(r.Context(), f, file.Filename, desc, uploader)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "加密分片上传链上存证失败: " + err.Error()})
		return
	}

	// 计算总耗时
	totalTime := time.Since(startTime)

	// 时间统计信息
	timeStats := map[string]interface{}{
		"totalTime":   totalTime.String(),
		"totalTimeMs": totalTime.Milliseconds(),
		"startTime":   startTime.Format("2006-01-02 15:04:05.000"),
		"endTime":     time.Now().Format("2006-01-02 15:04:05.000"),
	}

	resp := map[string]interface{}{
		"meta":         meta,
		"key":          gconv.String(key),
		"onChain":      true,
		"cipherSample": cipherSample,
		"chunkCount":   chunkCount,
		"chunkUrls":    chunkUrls,
		"fabricResult": fabricResult,
		"timeStats":    timeStats,
	}
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 200, Message: "加密分片上传并链上存证成功", Data: resp})
}

// POST /api/esg/download-encrypted
func (c *ESGController) DownloadEncrypted(r *ghttp.Request) {
	var req struct {
		Cid       string `json:"cid"`
		KeyBase64 string `json:"key"`
	}
	if err := r.Parse(&req); err != nil || req.Cid == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "参数错误"})
		return
	}
	key := gconv.Bytes(req.KeyBase64)
	if len(key) == 0 {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 400, Message: "密钥解码失败"})
		return
	}
	data, err := service.DownloadAndDecryptESGFile(r.Context(), req.Cid, key)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{Code: 500, Message: "下载解密失败: " + err.Error()})
		return
	}
	r.Response.Write(data)
}

// 其它如 DeleteFile、GetFileStats 等接口如需保留请重构为直接操作 service 层和现有模型，否则建议注释或删除。

// ========== ESG Analytics & Reporting 新功能 ==========

// GetReportDetail returns detailed report content
// GET /api/esg/report/:id
func (c *ESGController) GetReportDetail(r *ghttp.Request) {
	reportId := r.Get("id").String()
	if reportId == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "报告ID不能为空",
		})
		return
	}

	// Mock report content
	report := g.Map{
		"id":          reportId,
		"name":        "2024 Annual ESG Report",
		"type":        "Annual",
		"date":        "2024-12-01",
		"status":      "Completed",
		"sections": []g.Map{
			{
				"title":   "执行摘要",
				"content": "本报告总结了建筑在环境、社会和治理方面的综合表现...",
			},
			{
				"title":   "环境表现分析",
				"content": "碳排放量较去年下降15%，能源效率提升12%...",
			},
			{
				"title":   "社会表现分析",
				"content": "社区参与度提升，员工满意度达到85%...",
			},
			{
				"title":   "治理表现分析",
				"content": "治理结构完善，合规性达到100%...",
			},
			{
				"title":   "数据核证",
				"content": "所有数据已通过区块链验证，确保真实性和可追溯性...",
			},
			{
				"title":   "改进建议",
				"content": "建议升级HVAC系统，加强社区 engagement...",
			},
			{
				"title":   "未来展望",
				"content": "预计下一年度ESG评级将提升至AAA...",
			},
		},
	}

	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Success",
		Data:    report,
	})
}

// ExportReport exports report in specified format
// GET /api/esg/report/:id/export?format=pdf
func (c *ESGController) ExportReport(r *ghttp.Request) {
	reportId := r.Get("id").String()
	format := r.Get("format").String()
	if format == "" {
		format = "pdf"
	}

	if reportId == "" {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    400,
			Message: "报告ID不能为空",
		})
		return
	}

	// Mock export response
	r.Response.WriteJson(ghttp.DefaultHandlerResponse{
		Code:    200,
		Message: "Export started",
		Data: g.Map{
			"reportId":    reportId,
			"format":      format,
			"downloadUrl": "/api/esg/report/" + reportId + "/download?format=" + format,
			"status":      "ready",
		},
	})
}
