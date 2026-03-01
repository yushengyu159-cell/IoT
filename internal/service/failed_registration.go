package service

import (
	"fabric-sdk/internal/model"
	"time"

	"gorm.io/gorm"
)

// AddFailedRegistration 添加失败注册记录
func AddFailedRegistration(email, did, reason string) error {
	// 检查是否已存在失败记录
	var existing model.FailedRegistration
	err := DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		// 已存在，更新重试次数和状态
		existing.RetryCount++
		now := time.Now()
		existing.LastRetryAt = &now
		existing.Status = "retry"
		existing.Reason = reason
		existing.UpdatedAt = time.Now()
		
		return DB.Save(&existing).Error
	} else if err == gorm.ErrRecordNotFound {
		// 不存在，创建新记录
		failedReg := model.FailedRegistration{
			Email:      email,
			DID:        did,
			Reason:     reason,
			Status:     "failed",
			FailedAt:   time.Now(),
			RetryCount: 0,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		
		return DB.Create(&failedReg).Error
	}
	
	return err
}

// GetFailedRegistration 获取失败注册记录
func GetFailedRegistration(email string) (*model.FailedRegistration, error) {
	var failedReg model.FailedRegistration
	err := DB.Where("email = ?", email).First(&failedReg).Error
	if err != nil {
		return nil, err
	}
	return &failedReg, nil
}

// UpdateFailedRegistrationStatus 更新失败注册状态
func UpdateFailedRegistrationStatus(email, status string) error {
	return DB.Model(&model.FailedRegistration{}).
		Where("email = ?", email).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

// DeleteFailedRegistration 删除失败注册记录（注册成功后）
func DeleteFailedRegistration(email string) error {
	return DB.Where("email = ?", email).Delete(&model.FailedRegistration{}).Error
}

// GetAllFailedRegistrations 获取所有失败注册记录
func GetAllFailedRegistrations() ([]model.FailedRegistrationResponse, error) {
	var failedRegs []model.FailedRegistration
	err := DB.Order("failed_at DESC").Find(&failedRegs).Error
	if err != nil {
		return nil, err
	}

	var responses []model.FailedRegistrationResponse
	for _, reg := range failedRegs {
		responses = append(responses, model.FailedRegistrationResponse{
			ID:          reg.ID,
			Email:       reg.Email,
			DID:         reg.DID,
			Reason:      reg.Reason,
			Status:      reg.Status,
			FailedAt:    reg.FailedAt,
			RetryCount:  reg.RetryCount,
			LastRetryAt: reg.LastRetryAt,
			CreatedAt:   reg.CreatedAt,
		})
	}

	return responses, nil
}

// IsEmailInFailedList 检查邮箱是否在失败列表中
func IsEmailInFailedList(email string) bool {
	var count int64
	DB.Model(&model.FailedRegistration{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// GetFailedRegistrationStats 获取失败注册统计
func GetFailedRegistrationStats() (map[string]interface{}, error) {
	var totalCount, retryCount, successCount int64
	
	// 总失败数
	DB.Model(&model.FailedRegistration{}).Count(&totalCount)
	
	// 重试中数量
	DB.Model(&model.FailedRegistration{}).Where("status = ?", "retry").Count(&retryCount)
	
	// 成功数量（已删除的记录，这里统计为0，因为成功后会被删除）
	successCount = 0
	
	stats := map[string]interface{}{
		"total_failed":   totalCount,
		"retry_count":    retryCount,
		"success_count":  successCount,
		"failure_rate":   float64(totalCount) / float64(totalCount+successCount) * 100,
	}
	
	return stats, nil
}
