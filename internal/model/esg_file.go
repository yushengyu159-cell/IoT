package model

type ESGFile struct {
	ID           uint   `gorm:"primaryKey"`
	CID          string `gorm:"uniqueIndex;size:128"`           // 主文件CID（第一个分片）
	Filename     string `gorm:"size:255"`                       // 文件名
	Desc         string `gorm:"size:500"`                       // 文件描述
	Uploader     string `gorm:"size:128"`                       // 上传者
	UploadAt     string `gorm:"size:32"`                        // 上传时间
	Txid         string `gorm:"size:128"`                       // 区块链交易哈希
	ChunkCount   int    `gorm:"default:0"`                      // 分片数量
	ChunkSize    int    `gorm:"default:0"`                      // 分片大小
	AllCIDs      string `gorm:"size:2000"`                      // 所有分片CIDs（JSON格式）
	EncryptionKey string `gorm:"size:500"`                      // 加密密钥（Base64）
	IV           string `gorm:"size:100"`                       // 初始化向量（Base64）
	CipherSample string `gorm:"size:500"`                       // 密文样本（Base64）
	FileSize     int64  `gorm:"default:0"`                      // 文件大小
	UploadStartTime string `gorm:"size:32"`                     // 上传开始时间
	UploadEndTime   string `gorm:"size:32"`                     // 上传结束时间
	TotalTime      string `gorm:"size:32"`                      // 总耗时
}
