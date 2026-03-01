package model

type DID struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:64"`
	Phone        string `gorm:"size:32"`
	Email        string `gorm:"size:128"`
	Password     string `gorm:"size:128"`
	Role         string `gorm:"size:32"`
	Age          int
	CreatedAt    string `gorm:"size:32"`
	DID          string `gorm:"column:did;size:128;uniqueIndex"`
	
	// 新增字段 - 支持前端注册的完整信息
	FullName     string `gorm:"size:64"`        // 全名（替代info）
	BuildingName string `gorm:"size:128"`      // 建筑名称
	BuildingAddr string `gorm:"size:256"`      // 建筑地址
	BuildingType string `gorm:"size:64"`       // 建筑类型
	PropertyName string `gorm:"size:128"`      // 物业名称
	Occupation   string `gorm:"size:64"`       // 职业
	Institution  string `gorm:"size:128"`      // 机构名称
	Status       string `gorm:"size:32;default:'pending'"` // 注册状态：pending, verified, completed
	Language     string `gorm:"size:16;default:'zh'"`     // 用户语言偏好：en, zh, zh-TW
}

// TableName 显式指定表名为 dids
func (DID) TableName() string {
	return "dids"
}

// DIDChainRequest 链码DID登记请求结构（保持与现有链码兼容）
type DIDChainRequest struct {
	Email     string `json:"email"`
	Addresses string `json:"addresses"` // 逗号分隔的地址信息
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Info      string `json:"info"`      // 用户信息摘要
}

// DIDChainResponse 链码DID登记响应结构
type DIDChainResponse struct {
	DID       string `json:"did"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	TxID      string `json:"txID"`
}

// RegisterRequest 前端注册请求结构
type RegisterRequest struct {
	Email        string `json:"email" v:"required|email"`
	Password     string `json:"password" v:"required|min:6"`
	FullName     string `json:"full_name" v:"required"`
	Phone        string `json:"phone,omitempty"`    // 手机号
	Role         string `json:"role" v:"required|in:owner,property_manager,institution"`
	BuildingName string `json:"building_name"`      // Owner角色专用
	BuildingAddr string `json:"building_addr"`      // Owner角色专用
	BuildingType string `json:"building_type"`      // Owner角色专用
	PropertyName string `json:"property_name"`      // Property Manager角色专用
	Occupation   string `json:"occupation"`         // Property Manager和Institution角色专用
	Institution  string `json:"institution"`        // Institution角色专用
}

// RegisterResponse 注册响应结构
type RegisterResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	UserID       string                 `json:"user_id,omitempty"`
	RedirectURL  string                 `json:"redirect_url,omitempty"`
	ErrorDetails string                 `json:"error_details,omitempty"`
	ChaincodeData map[string]interface{} `json:"chaincode_data,omitempty"` // 链码返回的数据
	RegisterTime string                 `json:"register_time,omitempty"`   // 用户注册时间
}

// EmailListResponse 获取所有注册成功邮箱的响应结构
type EmailListResponse struct {
	Total  int      `json:"total"`  // 总邮箱数
	Emails []string `json:"emails"` // 邮箱列表
}

// UserDetailResponse 根据邮箱获取用户详细信息的响应结构
type UserDetailResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	Role         string `json:"role"`
	Phone        string `json:"phone,omitempty"`
	Age          int    `json:"age,omitempty"`
	DID          string `json:"did,omitempty"`
	BuildingName string `json:"building_name,omitempty"`
	BuildingAddr string `json:"building_addr,omitempty"`
	BuildingType string `json:"building_type,omitempty"`
	PropertyName string `json:"property_name,omitempty"`
	Occupation   string `json:"occupation,omitempty"`
	Institution  string `json:"institution,omitempty"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
