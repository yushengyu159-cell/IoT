package model

import (
	"time"
	"gorm.io/gorm"
)

// FailedRegistration 失败注册记录表
type FailedRegistration struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Email       string         `gorm:"uniqueIndex;size:255;not null" json:"email"` // 邮箱地址
	DID         string         `gorm:"size:255" json:"did"`                        // DID标识符
	Reason      string         `gorm:"type:text" json:"reason"`                     // 失败原因
	Status      string         `gorm:"size:50;default:'failed'" json:"status"`     // 状态: failed, retry, success
	FailedAt    time.Time      `gorm:"not null" json:"failed_at"`                  // 失败时间
	RetryCount  int            `gorm:"default:0" json:"retry_count"`               // 重试次数
	LastRetryAt *time.Time     `json:"last_retry_at"`                              // 最后重试时间
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// TableName 指定表名
func (FailedRegistration) TableName() string {
	return "failed_registrations"
}

// FailedRegistrationRequest 失败注册请求
type FailedRegistrationRequest struct {
	Email  string `json:"email" binding:"required,email"`
	Reason string `json:"reason"`
}

// FailedRegistrationResponse 失败注册响应
type FailedRegistrationResponse struct {
	ID          uint      `json:"id"`
	Email       string    `json:"email"`
	DID         string    `json:"did"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"`
	FailedAt    time.Time `json:"failed_at"`
	RetryCount  int       `json:"retry_count"`
	LastRetryAt *time.Time `json:"last_retry_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// FailedRegistrationListResponse 失败注册列表响应
type FailedRegistrationListResponse struct {
	Total int                        `json:"total"`
	List  []FailedRegistrationResponse `json:"list"`
}

