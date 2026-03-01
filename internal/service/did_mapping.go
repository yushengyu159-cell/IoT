package service

import (
	"fmt"
	"fabric-sdk/internal/model"
	"strings"
)

// DIDMappingService DID数据映射服务
type DIDMappingService struct{}

var DIDMapping = new(DIDMappingService)

// MapRegistrationToDIDChain 将前端注册数据映射到链码DID请求
func (s *DIDMappingService) MapRegistrationToDIDChain(registration *model.RegisterRequest) *model.DIDChainRequest {
	// 构建地址信息字符串
	addresses := s.buildAddressString(registration)
	
	// 构建用户信息摘要
	info := s.buildInfoString(registration)
	
	return &model.DIDChainRequest{
		Email:     registration.Email,
		Addresses: addresses,
		Phone:     registration.Phone, // 使用前端收集的手机号
		Password:  registration.Password,
		Info:      info,
	}
}

// buildAddressString 构建地址信息字符串
func (s *DIDMappingService) buildAddressString(registration *model.RegisterRequest) string {
	var addressParts []string
	
	// 根据角色构建不同的地址信息
	switch registration.Role {
	case "owner":
		// Owner角色：建筑名称 + 建筑地址 + 建筑类型
		if registration.BuildingName != "" {
			addressParts = append(addressParts, fmt.Sprintf("Building:%s", registration.BuildingName))
		}
		if registration.BuildingAddr != "" {
			addressParts = append(addressParts, fmt.Sprintf("Address:%s", registration.BuildingAddr))
		}
		if registration.BuildingType != "" {
			addressParts = append(addressParts, fmt.Sprintf("Type:%s", registration.BuildingType))
		}
		
	case "property_manager":
		// Property Manager角色：物业名称 + 职业
		if registration.PropertyName != "" {
			addressParts = append(addressParts, fmt.Sprintf("Property:%s", registration.PropertyName))
		}
		if registration.Occupation != "" {
			addressParts = append(addressParts, fmt.Sprintf("Occupation:%s", registration.Occupation))
		}
		
	case "institution":
		// Institution角色：机构名称 + 职业
		if registration.Institution != "" {
			addressParts = append(addressParts, fmt.Sprintf("Institution:%s", registration.Institution))
		}
		if registration.Occupation != "" {
			addressParts = append(addressParts, fmt.Sprintf("Occupation:%s", registration.Occupation))
		}
	}
	
	// 如果没有地址信息，返回默认值
	if len(addressParts) == 0 {
		return "No address information"
	}
	
	return strings.Join(addressParts, " | ")
}

// buildInfoString 构建用户信息摘要
func (s *DIDMappingService) buildInfoString(registration *model.RegisterRequest) string {
	var infoParts []string
	
	// 基本信息
	if registration.FullName != "" {
		infoParts = append(infoParts, fmt.Sprintf("Name:%s", registration.FullName))
	}
	
	// 角色信息
	if registration.Role != "" {
		infoParts = append(infoParts, fmt.Sprintf("Role:%s", registration.Role))
	}
	
	// 根据角色添加特定信息
	switch registration.Role {
	case "owner":
		// Owner角色：建筑相关信息
		if registration.BuildingName != "" {
			infoParts = append(infoParts, fmt.Sprintf("Building:%s", registration.BuildingName))
		}
		if registration.BuildingType != "" {
			infoParts = append(infoParts, fmt.Sprintf("Type:%s", registration.BuildingType))
		}
		
	case "property_manager":
		// Property Manager角色：物业相关信息
		if registration.PropertyName != "" {
			infoParts = append(infoParts, fmt.Sprintf("Property:%s", registration.PropertyName))
		}
		if registration.Occupation != "" {
			infoParts = append(infoParts, fmt.Sprintf("Occupation:%s", registration.Occupation))
		}
		
	case "institution":
		// Institution角色：机构相关信息
		if registration.Institution != "" {
			infoParts = append(infoParts, fmt.Sprintf("Institution:%s", registration.Institution))
		}
		if registration.Occupation != "" {
			infoParts = append(infoParts, fmt.Sprintf("Occupation:%s", registration.Occupation))
		}
	}
	
	// 如果没有信息，返回默认值
	if len(infoParts) == 0 {
		return "Basic user information"
	}
	
	return strings.Join(infoParts, " | ")
}

// MapDIDChainToFullDID 将链码DID响应映射到完整的DID模型
func (s *DIDMappingService) MapDIDChainToFullDID(
	chainResponse *model.DIDChainResponse,
	registration *model.RegisterRequest,
) *model.DID {
	return &model.DID{
		Email:        registration.Email,
		Password:     registration.Password,
		Role:         registration.Role,
		FullName:     registration.FullName,
		BuildingName: registration.BuildingName,
		BuildingAddr: registration.BuildingAddr,
		BuildingType: registration.BuildingType,
		PropertyName: registration.PropertyName,
		Occupation:   registration.Occupation,
		Institution:  registration.Institution,
		Status:       "completed",
		DID:          chainResponse.DID,
		CreatedAt:    chainResponse.Timestamp,
	}
}

// ValidateRegistrationData 验证注册数据的完整性
func (s *DIDMappingService) ValidateRegistrationData(registration *model.RegisterRequest) error {
	// 基础字段验证
	if registration.Email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
	if registration.Password == "" {
		return fmt.Errorf("密码不能为空")
	}
	if registration.FullName == "" {
		return fmt.Errorf("姓名不能为空")
	}
	if registration.Role == "" {
		return fmt.Errorf("角色不能为空")
	}
	
	// 角色特定字段验证
	switch registration.Role {
	case "owner":
		if registration.BuildingName == "" {
			return fmt.Errorf("Owner角色必须填写建筑名称")
		}
		if registration.BuildingAddr == "" {
			return fmt.Errorf("Owner角色必须填写建筑地址")
		}
		if registration.BuildingType == "" {
			return fmt.Errorf("Owner角色必须选择建筑类型")
		}
		
	case "property_manager":
		// Property Manager的字段都是可选的
		
	case "institution":
		// Institution的字段都是可选的
		
	default:
		return fmt.Errorf("无效的用户角色: %s", registration.Role)
	}
	
	return nil
}
