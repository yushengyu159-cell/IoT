package service

import (
	"fmt"
	"fabric-sdk/internal/model"
	"gorm.io/gorm"
)

// DIDStorageService DID数据库存储服务
type DIDStorageService struct {
	db *gorm.DB
}

var DIDStorage *DIDStorageService

// NewDIDStorageService 创建DID存储服务实例
func NewDIDStorageService(db *gorm.DB) *DIDStorageService {
	return &DIDStorageService{db: db}
}

// SaveDIDToDatabase 将完整的DID信息保存到数据库
func (s *DIDStorageService) SaveDIDToDatabase(did *model.DID) error {
	// 检查是否已存在
	var existingDID model.DID
	result := s.db.Where("email = ?", did.Email).First(&existingDID)
	
	if result.Error == nil {
		// 更新现有记录
		updateData := map[string]interface{}{
			"full_name":      did.FullName,
			"role":           did.Role,
			"building_name":  did.BuildingName,
			"building_addr":  did.BuildingAddr,
			"building_type":  did.BuildingType,
			"property_name":  did.PropertyName,
			"occupation":     did.Occupation,
			"institution":    did.Institution,
			"status":         did.Status,
			"did":            did.DID,
		}
		
		if err := s.db.Model(&existingDID).Updates(updateData).Error; err != nil {
			return fmt.Errorf("更新DID记录失败: %v", err)
		}
		
		return nil
	}
	
	// 创建新记录
	if err := s.db.Create(did).Error; err != nil {
		return fmt.Errorf("创建DID记录失败: %v", err)
	}
	
	return nil
}

// GetDIDByEmail 根据邮箱获取DID信息
func (s *DIDStorageService) GetDIDByEmail(email string) (*model.DID, error) {
	var did model.DID
	result := s.db.Where("email = ?", email).First(&did)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("DID记录不存在")
		}
		return nil, fmt.Errorf("查询DID记录失败: %v", result.Error)
	}
	
	return &did, nil
}

// GetDIDByDID 根据DID标识符获取DID信息
func (s *DIDStorageService) GetDIDByDID(didString string) (*model.DID, error) {
	var did model.DID
	result := s.db.Where("did = ?", didString).First(&did)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("DID记录不存在")
		}
		return nil, fmt.Errorf("查询DID记录失败: %v", result.Error)
	}
	
	return &did, nil
}

// UpdateDIDStatus 更新DID状态
func (s *DIDStorageService) UpdateDIDStatus(email, status string) error {
	result := s.db.Model(&model.DID{}).Where("email = ?", email).Update("status", status)
	
	if result.Error != nil {
		return fmt.Errorf("更新DID状态失败: %v", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到要更新的DID记录")
	}
	
	return nil
}

// ListDIDsByRole 根据角色列出DID
func (s *DIDStorageService) ListDIDsByRole(role string) ([]*model.DID, error) {
	var dids []*model.DID
	result := s.db.Where("role = ?", role).Find(&dids)
	
	if result.Error != nil {
		return nil, fmt.Errorf("查询角色DID列表失败: %v", result.Error)
	}
	
	return dids, nil
}

// SearchDIDs 搜索DID记录
func (s *DIDStorageService) SearchDIDs(query string) ([]*model.DID, error) {
	var dids []*model.DID
	
	// 支持多字段搜索
	searchQuery := fmt.Sprintf("%%%s%%", query)
	result := s.db.Where(
		"email LIKE ? OR full_name LIKE ? OR building_name LIKE ? OR institution LIKE ?",
		searchQuery, searchQuery, searchQuery, searchQuery,
	).Find(&dids)
	
	if result.Error != nil {
		return nil, fmt.Errorf("搜索DID记录失败: %v", result.Error)
	}
	
	return dids, nil
}

// DeleteDID 删除DID记录（软删除）
func (s *DIDStorageService) DeleteDID(email string) error {
	result := s.db.Model(&model.DID{}).Where("email = ?", email).Update("status", "deleted")
	
	if result.Error != nil {
		return fmt.Errorf("删除DID记录失败: %v", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到要删除的DID记录")
	}
	
	return nil
}

// GetDIDStatistics 获取DID统计信息
func (s *DIDStorageService) GetDIDStatistics() (map[string]interface{}, error) {
	var stats struct {
		TotalCount     int64
		OwnerCount     int64
		ManagerCount   int64
		InstitutionCount int64
		CompletedCount int64
		PendingCount   int64
	}
	
	// 总数量
	if err := s.db.Model(&model.DID{}).Count(&stats.TotalCount).Error; err != nil {
		return nil, fmt.Errorf("统计总数失败: %v", err)
	}
	
	// 按角色统计
	if err := s.db.Model(&model.DID{}).Where("role = ?", "owner").Count(&stats.OwnerCount).Error; err != nil {
		return nil, fmt.Errorf("统计Owner数量失败: %v", err)
	}
	
	if err := s.db.Model(&model.DID{}).Where("role = ?", "property_manager").Count(&stats.ManagerCount).Error; err != nil {
		return nil, fmt.Errorf("统计Property Manager数量失败: %v", err)
	}
	
	if err := s.db.Model(&model.DID{}).Where("role = ?", "institution").Count(&stats.InstitutionCount).Error; err != nil {
		return nil, fmt.Errorf("统计Institution数量失败: %v", err)
	}
	
	// 按状态统计
	if err := s.db.Model(&model.DID{}).Where("status = ?", "completed").Count(&stats.CompletedCount).Error; err != nil {
		return nil, fmt.Errorf("统计完成数量失败: %v", err)
	}
	
	if err := s.db.Model(&model.DID{}).Where("status = ?", "pending").Count(&stats.PendingCount).Error; err != nil {
		return nil, fmt.Errorf("统计待处理数量失败: %v", err)
	}
	
	return map[string]interface{}{
		"total_count":      stats.TotalCount,
		"owner_count":      stats.OwnerCount,
		"manager_count":    stats.ManagerCount,
		"institution_count": stats.InstitutionCount,
		"completed_count":  stats.CompletedCount,
		"pending_count":    stats.PendingCount,
	}, nil
}
