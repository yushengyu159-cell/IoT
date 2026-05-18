package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"fabric-sdk/internal/model"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
)

// 使用 model 包中定义的类型
// RegisterRequest 和 RegisterResponse 已在 model/did.go 中定义

// UserRegistration 用户注册记录
type UserRegistration struct {
	ID           uint      `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"-"` // 不返回密码
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	Phone        string    `json:"phone,omitempty"` // 手机号
	Age          int       `json:"age,omitempty"`   // 年龄
	DID          string    `json:"did,omitempty"`   // DID标识符
	BuildingName string    `json:"building_name,omitempty"`
	BuildingAddr string    `json:"building_addr,omitempty"`
	BuildingType string    `json:"building_type,omitempty"`
	PropertyName string    `json:"property_name,omitempty"`
	Occupation   string    `json:"occupation,omitempty"`
	Institution  string    `json:"institution,omitempty"`
	Status       string    `json:"status"` // pending, verified, completed
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// 内存存储用户注册信息（生产环境建议使用数据库）
var userRegistrations = make(map[string]*UserRegistration)

// per-email 互斥锁，避免同一邮箱并发注册
var emailLocks sync.Map // map[string]*sync.Mutex

func acquireEmailLock(email string) func() {
	val, _ := emailLocks.LoadOrStore(email, &sync.Mutex{})
	mu := val.(*sync.Mutex)
	mu.Lock()
	return func() { mu.Unlock() }
}

// RegisterUser 处理用户注册
func RegisterUser(req *model.RegisterRequest) (*model.RegisterResponse, error) {
	// 0. 同邮箱注册互斥，序列化处理
	unlock := acquireEmailLock(req.Email)
	defer unlock()
    // 1. 检查邮箱是否已注册（不早退，由后续分支执行存在即更新逻辑）

	// 2. 验证密码强度
	if !ValidatePassword(req.Password) {
		return &model.RegisterResponse{
			Success:      false,
			Message:      "密码长度必须在8-16位之间",
			ErrorDetails: "INVALID_PASSWORD",
		}, nil
	}

	// 3. 验证注册数据完整性
	if err := DIDMapping.ValidateRegistrationData(req); err != nil {
		return &model.RegisterResponse{
			Success:      false,
			Message:      err.Error(),
			ErrorDetails: "INVALID_REGISTRATION_DATA",
		}, nil
	}

	// 4. 映射到链码DID请求
	didChainReq := DIDMapping.MapRegistrationToDIDChain(req)

    // 5. 调用链码进行DID登记（改为单背书：资产创建式登记，与 /api/did/register 一致）
    chaincodeResult, err := Chaincode.CreateAssetWithMetadata(context.Background(),
		didChainReq.Email, // assetID：使用邮箱
		"USER_DID",        // color：复用为类型标识
		0,                 // size
		didChainReq.Email, // owner：邮箱
		0,                 // appraisedValue
	)

	// 添加调试日志
	if err != nil {
		g.Log().Error(nil, "❌ DID链码登记失败:", err)
	} else {
		g.Log().Info(nil, "✅ DID链码登记成功，返回结果:", chaincodeResult)
	}

    if err != nil {
        // 链码失败不再中断流程：记录错误，继续将资料写入数据库
        g.Log().Warning(nil, "⚠️ DID链码登记失败，继续落库并置为已完成:", err)
    }

    // 无论链码结果如何，统一执行资料合并更新（使用通用更新路径，带UPSERT兜底）
    _ = UpdateUserRegistration(req.Email, map[string]interface{}{
        "fullName":     req.FullName,
        "role":         req.Role,
        "phone":        req.Phone,
        "buildingName": req.BuildingName,
        "buildingAddr": req.BuildingAddr,
        "buildingType": req.BuildingType,
    })

	// 6. 创建用户注册记录
	user := &UserRegistration{
		Email:        req.Email,
		Password:     HashPassword(req.Password),
		FullName:     req.FullName,
		Role:         req.Role,
		Phone:        req.Phone,                                      // 手机号
		Age:          20,                                             // 默认年龄20岁
		DID:          extractDIDFromChaincodeResult(chaincodeResult), // 使用链码返回的DID
		BuildingName: req.BuildingName,
		BuildingAddr: req.BuildingAddr,
		BuildingType: req.BuildingType,
		PropertyName: req.PropertyName,
		Occupation:   req.Occupation,
		Institution:  req.Institution,
		Status:       "completed", // 注册成功后直接完成，无需审核
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// 7. 根据角色处理不同的注册流程
	response, err := ProcessRoleBasedRegistration(user)
	if err != nil {
		return &model.RegisterResponse{
			Success:      false,
			Message:      "注册处理失败",
			ErrorDetails: err.Error(),
		}, nil
	}

	// 5. 存储/更新用户信息到MySQL数据库
	if IsEmailRegistered(req.Email) {
		g.Log().Info(nil, "ℹ️ 邮箱已存在，执行资料更新并置为已完成:", req.Email)
		if DB == nil {
			return &model.RegisterResponse{Success: false, Message: "数据库未初始化"}, nil
		}
		updates := map[string]interface{}{
			"full_name":         user.FullName,
			"role":              user.Role,
			"phone":             user.Phone,
			"building_name":     user.BuildingName,
			"building_addr":     user.BuildingAddr,
			"building_type":     user.BuildingType,
			"property_name":     user.PropertyName,
			"occupation":        user.Occupation,
			"institution":       user.Institution,
			"status":            "completed",
		}
		if err := DB.Table("dids").Where("email = ?", req.Email).Updates(updates).Error; err != nil {
			g.Log().Error(nil, "❌ 更新用户信息失败:", err)
			return &model.RegisterResponse{
				Success:      false,
				Message:      "更新用户信息失败",
				ErrorDetails: err.Error(),
			}, nil
		} else {
			RefreshUserCache(req.Email)
		}
		if chaincodeResult != nil {
			response.ChaincodeData = chaincodeResult
		}
		return response, nil
	}
	
	g.Log().Info(nil, "🔍 邮箱不存在，准备创建新用户:", req.Email)
	if err := SaveUserToDatabase(user); err != nil {
		g.Log().Error(nil, "❌ 保存用户到数据库失败:", err)
		return &model.RegisterResponse{
			Success:      false,
			Message:      "用户信息保存失败",
			ErrorDetails: err.Error(),
		}, nil
	}
	
	g.Log().Info(nil, "✅ 用户信息保存成功，继续后续流程:", req.Email)

	// 6. 同时保存到内存缓存（用于快速访问）
	userRegistrations[req.Email] = user

	// 7. 记录注册日志
	g.Log().Info(nil, "✅ 用户注册成功:", req.Email, "角色:", req.Role)

	// 8. 将链码返回的数据添加到响应中
	if chaincodeResult != nil {
		response.ChaincodeData = chaincodeResult
	}

	// 9. 设置真实的注册时间
	response.RegisterTime = user.CreatedAt.Format("2006-01-02 15:04:05")

	return response, nil
}

// extractDIDFromChaincodeResult 从链码返回结果中提取DID
func extractDIDFromChaincodeResult(chaincodeResult map[string]interface{}) string {
	if chaincodeResult == nil {
		g.Log().Warning(nil, "⚠️ 链码返回结果为空，使用本地生成DID")
		return GenerateDID("fallback")
	}

	// 尝试从链码返回结果中提取DID
	if did, ok := chaincodeResult["did"].(string); ok && did != "" {
		g.Log().Info(nil, "✅ 使用链码返回的DID:", did)
		return did
	}

	// 如果链码没有返回DID，降级到本地生成
	g.Log().Warning(nil, "⚠️ 链码未返回DID，使用本地生成")
	return GenerateDID("fallback")
}

// ProcessRoleBasedRegistration 根据角色处理注册流程
func ProcessRoleBasedRegistration(user *UserRegistration) (*model.RegisterResponse, error) {
	switch user.Role {
	case "owner":
		// Owner需要完整的建筑信息
		if user.BuildingName == "" || user.BuildingAddr == "" || user.BuildingType == "" {
			return &model.RegisterResponse{
				Success:      false,
				Message:      "Owner角色需要完整的建筑信息",
				ErrorDetails: "INCOMPLETE_BUILDING_INFO",
			}, nil
		}
		// Owner注册完成，跳转到仪表板界面
		return &model.RegisterResponse{
			Success:     true,
			Message:     "Owner注册完成",
			UserID:      fmt.Sprintf("user_%d", time.Now().Unix()),
			RedirectURL: "/static/esg-dashboard.html",
		}, nil

	case "property_manager":
		// Property Manager注册完成，跳转到仪表板界面
		user.Status = "completed"
		return &model.RegisterResponse{
			Success:     true,
			Message:     "Property Manager注册完成",
			UserID:      fmt.Sprintf("user_%d", time.Now().Unix()),
			RedirectURL: "/static/esg-dashboard.html",
		}, nil

	case "institution":
		// Institution注册完成，跳转到仪表板界面
		user.Status = "completed"
		return &model.RegisterResponse{
			Success:     true,
			Message:     "Institution注册完成",
			UserID:      fmt.Sprintf("user_%d", time.Now().Unix()),
			RedirectURL: "/static/esg-dashboard.html",
		}, nil

	default:
		return &model.RegisterResponse{
			Success:      false,
			Message:      "无效的用户角色",
			ErrorDetails: "INVALID_ROLE",
		}, nil
	}
}

// IsEmailVerified 检查邮箱是否已验证
func IsEmailVerified(email string) bool {
	// 调用邮箱服务的验证状态检查
	verification, exists := GetVerificationStatus(email)
	if !exists {
		return false
	}
	return verification.Used && !time.Now().After(verification.ExpiresAt)
}

// IsEmailRegistered 检查邮箱是否已注册
func IsEmailRegistered(email string) bool {
	// 优先检查数据库，确保数据一致性
	if DB != nil {
		var count int64
		var dbErr error

		// 强制刷新数据库连接，避免缓存问题
		if dbErr = DB.Raw("SELECT COUNT(*) FROM `dids` WHERE email = ?", email).Scan(&count).Error; dbErr == nil {
			g.Log().Infof(nil, "🔍 数据库查询结果: email=%s, count=%d", email, count)

			if count > 0 {
				// 数据库中存在，同步到内存缓存
				if _, exists := userRegistrations[email]; !exists {
					// 如果内存缓存中没有，从数据库加载
					if user, _ := GetUserRegistration(email); user != nil {
						userRegistrations[email] = user
					}
				}
				return true
			} else {
				// 数据库中不存在，清理内存缓存
				delete(userRegistrations, email)
				return false
			}
		}
		// 如果数据库查询出错，记录日志
		g.Log().Errorf(nil, "数据库查询失败: %v", dbErr)
	}

	// 数据库不可用时，使用内存缓存作为备选
	if _, exists := userRegistrations[email]; exists {
		return true
	}

	return false
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) bool {
	return len(password) >= 8 && len(password) <= 16
}

// HashPassword 密码哈希
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// ValidatePasswordHash 验证密码哈希
func ValidatePasswordHash(password, hashedPassword string) bool {
	inputHash := HashPassword(password)
	return inputHash == hashedPassword
}



// GetUserRegistration 获取用户注册信息
func GetUserRegistration(email string) (*UserRegistration, bool) {
	// 首先从数据库获取
	if DB != nil {
		var didRecord model.DID
		if err := DB.Where("email = ?", email).First(&didRecord).Error; err == nil {
			// 添加调试日志
			g.Log().Info(nil, "🔍 数据库查询结果:")
			g.Log().Info(nil, "   - Email:", didRecord.Email)
			g.Log().Info(nil, "   - Phone:", didRecord.Phone)
			g.Log().Info(nil, "   - CreatedAt:", didRecord.CreatedAt)
			g.Log().Info(nil, "   - Role:", didRecord.Role)

			// 转换为UserRegistration结构
			user := &UserRegistration{
				Email:        didRecord.Email,
				Password:     didRecord.Password,
				FullName:     didRecord.FullName,
				Role:         didRecord.Role,
				Phone:        didRecord.Phone,
				Age:          int(didRecord.Age),
				DID:          didRecord.DID,
				BuildingName: didRecord.BuildingName,
				BuildingAddr: didRecord.BuildingAddr,
				BuildingType: didRecord.BuildingType,
				PropertyName: didRecord.PropertyName,
				Occupation:   didRecord.Occupation,
				Institution:  didRecord.Institution,
				Status:       didRecord.Status,
				CreatedAt:    parseCreatedAt(didRecord.CreatedAt), // 解析数据库中的创建时间
				UpdatedAt:    time.Now(),
			}

			// 同时更新内存缓存
			userRegistrations[email] = user
			return user, true
		}
	}

	// 如果数据库中没有，从内存缓存获取
	if user, exists := userRegistrations[email]; exists {
		return user, true
	}

	return nil, false
}

// GetUserRegistrationByDID 根据DID获取用户注册信息
func GetUserRegistrationByDID(did string) (*UserRegistration, bool) {
	if DB != nil {
		var didRecord model.DID
		if err := DB.Where("did = ?", did).First(&didRecord).Error; err == nil {
			user := &UserRegistration{
				Email:        didRecord.Email,
				Password:     didRecord.Password,
				FullName:     didRecord.FullName,
				Role:         didRecord.Role,
				Phone:        didRecord.Phone,
				Age:          int(didRecord.Age),
				DID:          didRecord.DID,
				BuildingName: didRecord.BuildingName,
				BuildingAddr: didRecord.BuildingAddr,
				BuildingType: didRecord.BuildingType,
				PropertyName: didRecord.PropertyName,
				Occupation:   didRecord.Occupation,
				Institution:  didRecord.Institution,
				Status:       didRecord.Status,
				CreatedAt:    parseCreatedAt(didRecord.CreatedAt), // 解析数据库中的创建时间
				UpdatedAt:    time.Now(),
			}
			// 同步缓存（按邮箱键）
			userRegistrations[didRecord.Email] = user
			return user, true
		}
	}
	return nil, false
}

// GetUserByFullName 根据姓名查询用户注册信息
func GetUserByFullName(fullName string) (*UserRegistration, bool) {
	if DB == nil {
		return nil, false
	}

	var didRecord model.DID
	if err := DB.Where("full_name = ?", fullName).First(&didRecord).Error; err == nil {
		// 转换为UserRegistration结构
		user := &UserRegistration{
			Email:        didRecord.Email,
			Password:     didRecord.Password,
			FullName:     didRecord.FullName,
			Role:         didRecord.Role,
			Phone:        didRecord.Phone,
			Age:          int(didRecord.Age),
			DID:          didRecord.DID,
			BuildingName: didRecord.BuildingName,
			BuildingAddr: didRecord.BuildingAddr,
			BuildingType: didRecord.BuildingType,
			PropertyName: didRecord.PropertyName,
			Occupation:   didRecord.Occupation,
			Institution:  didRecord.Institution,
			Status:       didRecord.Status,
			CreatedAt:    parseCreatedAt(didRecord.CreatedAt), // 解析数据库中的创建时间
			UpdatedAt:    time.Now(),
		}

		// 同时更新内存缓存
		userRegistrations[didRecord.Email] = user
		return user, true
	}

	return nil, false
}

// ClearUserFromCache 从内存缓存中清除用户
func ClearUserFromCache(email string) {
	delete(userRegistrations, email)
	g.Log().Infof(nil, "✅ 已从内存缓存清除用户: %s", email)
}

// ForceRefreshDatabase 强制刷新数据库连接，清理可能的缓存
func ForceRefreshDatabase() {
	if DB != nil {
		// 执行一个简单的查询来刷新连接
		var result int64
		if err := DB.Raw("SELECT 1").Scan(&result).Error; err == nil {
			g.Log().Infof(nil, "✅ 数据库连接已刷新")
		} else {
			g.Log().Errorf(nil, "数据库连接刷新失败: %v", err)
		}
	}
}

// RefreshUserCache 刷新用户缓存，确保与数据库同步
func RefreshUserCache(email string) {
	// 清除内存缓存
	ClearUserFromCache(email)

	// 强制刷新数据库查询，避免缓存问题
	if DB != nil {
		var count int64
		if err := DB.Raw("SELECT COUNT(*) FROM `dids` WHERE email = ?", email).Scan(&count).Error; err == nil {
			g.Log().Infof(nil, "🔄 强制刷新数据库查询: email=%s, count=%d", email, count)

			if count > 0 {
				if user, _ := GetUserRegistration(email); user != nil {
					userRegistrations[email] = user
					g.Log().Infof(nil, "✅ 已刷新用户缓存: %s", email)
				}
			} else {
				g.Log().Infof(nil, "✅ 数据库中不存在用户: %s", email)
			}
		} else {
			g.Log().Errorf(nil, "强制刷新数据库查询失败: %v", err)
		}
	}
}

// UpdateUserRegistration 更新用户注册信息
func UpdateUserRegistration(email string, updates map[string]interface{}) error {
	// 1) 内存缓存更新（向后兼容）
	if user, exists := userRegistrations[email]; exists {
		if buildingName, ok := updates["buildingName"].(string); ok {
			user.BuildingName = buildingName
		}
		if buildingAddr, ok := updates["buildingAddr"].(string); ok {
			user.BuildingAddr = buildingAddr
		}
		if buildingType, ok := updates["buildingType"].(string); ok {
			user.BuildingType = buildingType
		}
		if phone, ok := updates["phone"].(string); ok {
			user.Phone = phone
		}
		user.UpdatedAt = time.Now()
		userRegistrations[email] = user
	}

	// 2) 持久化到 MySQL（关键修复）
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	dbUpdates := map[string]interface{}{}
	if v, ok := updates["buildingName"].(string); ok {
		dbUpdates["building_name"] = v
	}
	if v, ok := updates["buildingAddr"].(string); ok {
		dbUpdates["building_addr"] = v
	}
	if v, ok := updates["buildingType"].(string); ok {
		dbUpdates["building_type"] = v
	}
	if v, ok := updates["phone"].(string); ok {
		dbUpdates["phone"] = v
	}
	// Removed updated_timestamp - field doesn't exist in DID model

	if len(dbUpdates) > 0 {
		tx := DB.Table("dids").Where("email = ?", email).Updates(dbUpdates)
		if tx.Error != nil {
			return fmt.Errorf("更新用户信息失败: %v", tx.Error)
		}
		// 若没有任何行被更新，则尝试插入一条新记录（UPSERT兜底）
		if tx.RowsAffected == 0 {
			g.Log().Warning(nil, "⚠️ 未找到用户记录，执行插入兜底:", email)
			did := ""
			if u, ok := userRegistrations[email]; ok {
				did = u.DID
			}
			record := &model.DID{
				Email:        email,
				FullName:     fmt.Sprintf("%v", updates["fullName"]),
				Role:         "owner",
				Phone:        fmt.Sprintf("%v", updates["phone"]),
				DID:          did,
				BuildingName: fmt.Sprintf("%v", updates["buildingName"]),
				BuildingAddr: fmt.Sprintf("%v", updates["buildingAddr"]),
				BuildingType: fmt.Sprintf("%v", updates["buildingType"]),
				Status:       "pending",
				CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
			}
			if err := DB.Create(record).Error; err != nil {
				return fmt.Errorf("插入用户记录失败: %v", err)
			}
			g.Log().Info(nil, "✅ 已为缺失用户插入记录:", email)
		}
	}

	// 3) 刷新缓存
	RefreshUserCache(email)
	g.Log().Info(nil, "✅ 用户信息已持久化更新:", email)
	return nil
}

// CompleteUserRegistration 完成用户注册（已废弃，注册流程已改为直接完成）
func CompleteUserRegistration(email string) error {
	// 内存状态更新
	if user, exists := userRegistrations[email]; exists {
		user.Status = "completed"
		user.UpdatedAt = time.Now()
		userRegistrations[email] = user
	}

	// 数据库状态更新
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	
	// 先检查用户是否存在
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM `dids` WHERE email = ?", email).Scan(&count).Error; err != nil {
		g.Log().Error(nil, "❌ 检查用户存在性失败:", err)
		return fmt.Errorf("检查用户存在性失败: %v", err)
	}
	
	if count == 0 {
		g.Log().Error(nil, "❌ 用户不存在，无法完成注册:", email)
		return fmt.Errorf("用户不存在，无法完成注册: %s", email)
	}
	
	g.Log().Info(nil, "🔍 用户存在，准备更新状态:", email)
	if err := DB.Table("dids").Where("email = ?", email).Updates(map[string]interface{}{
		"status": "completed",
	}).Error; err != nil {
		g.Log().Error(nil, "❌ 更新注册状态失败:", err)
		return fmt.Errorf("更新注册状态失败: %v", err)
	}

	RefreshUserCache(email)
	g.Log().Info(nil, "✅ 用户注册完成并已持久化:", email)
	return nil
}

// SetUserStatusPendingReview 将用户状态置为待审核
func SetUserStatusPendingReview(email string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}
    if err := DB.Table("dids").Where("email = ?", email).Updates(map[string]interface{}{
        "status": "pending_review",
    }).Error; err != nil {
		return fmt.Errorf("设置待审核失败: %v", err)
	}
	// 同步内存缓存
	if u, ok := userRegistrations[email]; ok {
		u.Status = "pending_review"
		u.UpdatedAt = time.Now()
		userRegistrations[email] = u
	}
	return nil
}

// AdminApproveUser 审核通过并发送通知
func AdminApproveUser(email string) error {
	return AdminApproveUserWithLanguage(email, "", "")
}

// AdminApproveUserWithLanguage 审核通过并发送通知（支持指定语言）
func AdminApproveUserWithLanguage(email, reason, language string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
    if err := DB.Table("dids").Where("email = ?", email).Updates(map[string]interface{}{
        "status": "completed",
    }).Error; err != nil {
		return fmt.Errorf("审核通过失败: %v", err)
	}
	
	// 删除失败注册记录（如果存在）
	DeleteFailedRegistration(email)
	
	RefreshUserCache(email)
	_ = SendAuditResultEmailWithLanguage(email, true, reason, language)
	return nil
}

// AdminRejectUser 审核拒绝并发送通知
func AdminRejectUser(email, reason string) error {
	return AdminRejectUserWithLanguage(email, reason, "")
}

// AdminRejectUserWithLanguage 审核拒绝并发送通知（支持指定语言）
func AdminRejectUserWithLanguage(email, reason, language string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	
	// 获取用户信息用于邮件通知
	var did model.DID
	if err := DB.Where("email = ?", email).First(&did).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}
	
	// 发送拒绝邮件通知
	_ = SendAuditResultEmailWithLanguage(email, false, reason, language)
	
	// 从数据库中删除被拒绝的用户
	if err := DB.Where("email = ?", email).Delete(&model.DID{}).Error; err != nil {
		return fmt.Errorf("删除用户失败: %v", err)
	}
	
	// 从缓存中移除用户
	RefreshUserCache(email)
	
	// 添加失败注册记录，允许用户重新注册
	AddFailedRegistration(email, did.DID, reason)
	
	g.Log().Info(nil, "✅ 用户审核拒绝，已删除用户:", email)
	return nil
}

// GetAllUserRegistrations 获取所有用户注册信息（管理员功能）
func GetAllUserRegistrations() []*UserRegistration {
	users := make([]*UserRegistration, 0, len(userRegistrations))
	for _, user := range userRegistrations {
		users = append(users, user)
	}
	return users
}

// DeleteUserRegistration 删除用户注册信息（管理员功能）
func DeleteUserRegistration(email string) error {
	if _, exists := userRegistrations[email]; !exists {
		return fmt.Errorf("用户不存在")
	}

	delete(userRegistrations, email)
	g.Log().Info(nil, "🗑️ 用户注册信息已删除:", email)
	return nil
}

// GenerateDID 生成DID标识符
func GenerateDID(email string) string {
	// 使用邮箱和时间戳生成唯一的DID
	timestamp := time.Now().Unix()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s_%d", email, timestamp)))
	// 统一规范为 did:example:<short-hash>
	return fmt.Sprintf("did:example:%x", hash[:8])
}

// SaveUserToDatabase 保存用户信息到MySQL数据库
func SaveUserToDatabase(user *UserRegistration) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 确保DID不为空，如果为空则生成一个唯一的DID
	did := user.DID
	if did == "" {
		did = GenerateDID(user.Email)
		g.Log().Warning(nil, "⚠️ 用户DID为空，生成新的DID:", did)
	}

	// 如果该邮箱已存在，则执行按邮箱的合并更新（UPSERT 语义）
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM `dids` WHERE email = ?", user.Email).Scan(&count).Error; err == nil && count > 0 {
		updates := map[string]interface{}{
			// 仅当上游提供非空值时才更新，避免把已有值覆盖为空
		}
		if user.FullName != "" {
			updates["full_name"] = user.FullName
		}
		if user.Role != "" {
			updates["role"] = user.Role
		}
		if user.Phone != "" {
			updates["phone"] = user.Phone
		}
		if user.Age != 0 {
			updates["age"] = user.Age
		}
		if did != "" {
			updates["did"] = did
		}
		if user.BuildingName != "" {
			updates["building_name"] = user.BuildingName
		}
		if user.BuildingAddr != "" {
			updates["building_addr"] = user.BuildingAddr
		}
		if user.BuildingType != "" {
			updates["building_type"] = user.BuildingType
		}
		if user.PropertyName != "" {
			updates["property_name"] = user.PropertyName
		}
		if user.Occupation != "" {
			updates["occupation"] = user.Occupation
		}
		if user.Institution != "" {
			updates["institution"] = user.Institution
		}
		// 注册时统一设置为 completed 状态
		updates["status"] = "completed"

		if len(updates) > 0 { // 至少包含一个业务字段
			if err := DB.Table("dids").Where("email = ?", user.Email).Updates(updates).Error; err != nil {
				return fmt.Errorf("更新用户记录失败: %v", err)
			}
			g.Log().Info(nil, "🔄 用户信息已按邮箱合并更新:", user.Email)
			return nil
		}
		// 如果没有可更新的字段，直接返回视为成功
		g.Log().Info(nil, "ℹ️ 无需更新，保持现有记录:", user.Email)
		return nil
	}

	// 不存在则创建新记录
	didRecord := &model.DID{
		Email:        user.Email,
		Password:     user.Password,
		FullName:     user.FullName,
		Role:         user.Role,
		Phone:        user.Phone,
		Age:          user.Age,
		DID:          did,
		BuildingName: user.BuildingName,
		BuildingAddr: user.BuildingAddr,
		BuildingType: user.BuildingType,
		PropertyName: user.PropertyName,
		Occupation:   user.Occupation,
		Institution:  user.Institution,
		Status:       user.Status,
		Language:     "zh", // 添加语言字段
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	g.Log().Info(nil, "🔍 准备创建用户记录:", user.Email, "DID:", did, "Status:", user.Status)
	
	if err := DB.Create(didRecord).Error; err != nil {
		g.Log().Error(nil, "❌ 保存用户记录失败:", err, "Email:", user.Email)
		return fmt.Errorf("保存用户记录失败: %v", err)
	}

	// 验证插入是否成功
	var verifyCount int64
	if err := DB.Raw("SELECT COUNT(*) FROM `dids` WHERE email = ?", user.Email).Scan(&verifyCount).Error; err != nil {
		g.Log().Error(nil, "❌ 验证插入结果失败:", err)
		return fmt.Errorf("验证插入结果失败: %v", err)
	}
	
	if verifyCount == 0 {
		g.Log().Error(nil, "❌ 插入后验证失败，记录不存在:", user.Email)
		return fmt.Errorf("插入后验证失败，记录不存在")
	}

	g.Log().Info(nil, "✅ 用户信息已保存到数据库:", user.Email, "DID:", did, "验证成功")
	return nil
}

// ResetPassword 重置用户密码（按邮箱）
func ResetPassword(email string, newPassword string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}
	if email == "" || newPassword == "" {
		return fmt.Errorf("邮箱与新密码不能为空")
	}
	// 先检查是否存在
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM `dids` WHERE email = ?", email).Scan(&count).Error; err != nil {
		return fmt.Errorf("查询用户失败: %v", err)
	}
	if count == 0 {
		return fmt.Errorf("用户不存在")
	}
	// 更新密码为新的哈希
	hashed := HashPassword(newPassword)
	if err := DB.Exec("UPDATE `dids` SET password = ? WHERE email = ?", hashed, email).Error; err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}
	// 同步刷新内存缓存
	if user, ok := userRegistrations[email]; ok {
		user.Password = hashed
		userRegistrations[email] = user
	}
	g.Log().Info(nil, "✅ 用户密码已重置:", email)
	return nil
}

// GetUploaderInfoFromDB 直接从数据库查询上传者信息
func GetUploaderInfoFromDB(uploader string) (*UserRegistration, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	var didRecord model.DID

	// 首先尝试按邮箱查询
	if err := DB.Where("email = ?", uploader).First(&didRecord).Error; err == nil {
		return &UserRegistration{
			Email:        didRecord.Email,
			Password:     didRecord.Password,
			FullName:     didRecord.FullName,
			Role:         didRecord.Role,
			Phone:        didRecord.Phone,
			Age:          int(didRecord.Age),
			DID:          didRecord.DID,
			BuildingName: didRecord.BuildingName,
			BuildingAddr: didRecord.BuildingAddr,
			BuildingType: didRecord.BuildingType,
			PropertyName: didRecord.PropertyName,
			Occupation:   didRecord.Occupation,
			Institution:  didRecord.Institution,
			Status:       didRecord.Status,
			CreatedAt:    parseCreatedAt(didRecord.CreatedAt), // 解析数据库中的创建时间
			UpdatedAt:    time.Now(),
		}, nil
	}

	// 如果按邮箱查不到，尝试按姓名查询
	if err := DB.Where("full_name = ?", uploader).First(&didRecord).Error; err == nil {
		return &UserRegistration{
			Email:        didRecord.Email,
			Password:     didRecord.Password,
			FullName:     didRecord.FullName,
			Role:         didRecord.Role,
			Phone:        didRecord.Phone,
			Age:          int(didRecord.Age),
			DID:          didRecord.DID,
			BuildingName: didRecord.BuildingName,
			BuildingAddr: didRecord.BuildingAddr,
			BuildingType: didRecord.BuildingType,
			PropertyName: didRecord.PropertyName,
			Occupation:   didRecord.Occupation,
			Institution:  didRecord.Institution,
			Status:       didRecord.Status,
			CreatedAt:    parseCreatedAt(didRecord.CreatedAt), // 解析数据库中的创建时间
			UpdatedAt:    time.Now(),
		}, nil
	}

	// 如果都查不到，返回 nil 和错误
	return nil, fmt.Errorf("未找到上传者信息: %s", uploader)
}

// parseCreatedAt 解析数据库中的CreatedAt字符串字段
func parseCreatedAt(createdAtStr string) time.Time {
	if createdAtStr == "" {
		return time.Now()
	}

	// 尝试解析 "2006-01-02 15:04:05" 格式
	if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
		return t
	}

	// 尝试解析 "2006-01-02T15:04:05Z" 格式
	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		return t
	}

	// 尝试解析 "2006-01-02 15:04:05 +0000 UTC" 格式
	if t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", createdAtStr); err == nil {
		return t
	}

	// 如果都解析失败，返回当前时间
	g.Log().Warning(nil, "⚠️ 无法解析CreatedAt时间:", createdAtStr, "使用当前时间")
	return time.Now()
}

// GetAllRegisteredEmails 获取所有注册成功的邮箱账号
func GetAllRegisteredEmails() (*model.EmailListResponse, error) {
	ctx := context.Background()

	// 从数据库获取所有注册成功的邮箱
	if DB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	// 只查询完全注册成功的用户邮箱（状态为completed）
	var results []model.DID
	err := DB.Where("status = ?", "completed").Order("created_at DESC").Find(&results).Error
	if err != nil {
		g.Log().Error(ctx, "❌ 查询邮箱列表失败:", err)
		return nil, fmt.Errorf("查询邮箱列表失败: %v", err)
	}

	// 提取邮箱地址
	var emails []string
	for _, record := range results {
		if record.Email != "" {
			emails = append(emails, record.Email)
		}
	}

	response := &model.EmailListResponse{
		Total:  len(emails),
		Emails: emails,
	}

	g.Log().Info(ctx, "✅ 成功获取邮箱列表，总数:", len(emails))
	return response, nil
}

// GetPendingUsersForAdmin 获取待审核用户（管理员审核用）
func GetPendingUsersForAdmin() (*model.EmailListResponse, error) {
	ctx := context.Background()

	// 从数据库获取待审核用户
	if DB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	// 查询待审核用户邮箱
	var results []model.DID
	err := DB.Where("status IN (?)", []string{"pending", "pending_review"}).Order("created_at DESC").Find(&results).Error
	if err != nil {
		g.Log().Error(ctx, "❌ 查询待审核用户列表失败:", err)
		return nil, fmt.Errorf("查询待审核用户列表失败: %v", err)
	}

	// 提取邮箱地址
	var emails []string
	for _, record := range results {
		email := record.Email
		if email != "" {
			emails = append(emails, email)
		}
	}

	response := &model.EmailListResponse{
		Total:  len(emails),
		Emails: emails,
	}

	g.Log().Info(ctx, "✅ 成功获取待审核用户列表，总数:", len(emails))
	return response, nil
}

// CleanupIncompleteRegistrations 清理未完成注册的账号
func CleanupIncompleteRegistrations() error {
	ctx := context.Background()
	
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 删除状态为 pending 且超过24小时未完成的注册
	cutoffTime := time.Now().Add(-24 * time.Hour)
	
	// 查询需要清理的账号
	var incompleteUsers []model.DID
	err := DB.Where("status = ? AND created_at < ?", "pending", cutoffTime).Find(&incompleteUsers).Error
	if err != nil {
		g.Log().Error(ctx, "❌ 查询未完成注册账号失败:", err)
		return fmt.Errorf("查询未完成注册账号失败: %v", err)
	}

	if len(incompleteUsers) == 0 {
		g.Log().Info(ctx, "ℹ️ 没有需要清理的未完成注册账号")
		return nil
	}

	// 删除未完成的注册记录
	for _, user := range incompleteUsers {
		// 从内存中删除
		delete(userRegistrations, user.Email)
		
		// 从数据库中删除
		if err := DB.Delete(&user).Error; err != nil {
			g.Log().Error(ctx, "❌ 删除未完成注册账号失败:", user.Email, err)
			continue
		}
		
		g.Log().Info(ctx, "🗑️ 已清理未完成注册账号:", user.Email)
	}

	g.Log().Info(ctx, "✅ 清理完成，共清理", len(incompleteUsers), "个未完成注册账号")
	return nil
}

// GetAllUsersForAdmin 获取所有用户（管理员用，显示待审核和已完成状态的用户）
func GetAllUsersForAdmin() (*model.EmailListResponse, error) {
	ctx := context.Background()

	// 从数据库获取所有用户邮箱（包括待审核和已完成的）
	if DB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	// 查询待审核和已完成状态的用户邮箱
	var results []model.DID
	err := DB.Where("status IN (?, ?)", "pending_review", "completed").Order("created_at DESC").Find(&results).Error
	if err != nil {
		g.Log().Error(ctx, "❌ 查询所有用户列表失败:", err)
		return nil, fmt.Errorf("查询所有用户列表失败: %v", err)
	}

	// 提取邮箱地址
	var emails []string
	for _, record := range results {
		if record.Email != "" {
			emails = append(emails, record.Email)
		}
	}

	response := &model.EmailListResponse{
		Total:  len(emails),
		Emails: emails,
	}

	g.Log().Info(ctx, "✅ 成功获取所有用户列表，总数:", len(emails), "（包括待审核和已完成）")
	return response, nil
}

// GetUserDetailByEmail 根据邮箱获取用户详细信息
func GetUserDetailByEmail(email string) (*model.UserDetailResponse, error) {
	ctx := context.Background()

	if email == "" {
		return nil, fmt.Errorf("邮箱地址不能为空")
	}

	// 从数据库获取用户详细信息
	if DB == nil {
		return nil, fmt.Errorf("数据库连接未初始化")
	}

	// 查询指定邮箱的用户信息
	var result model.DID
	err := DB.Where("email = ?", email).First(&result).Error
	if err != nil {
		g.Log().Error(ctx, "❌ 查询用户详细信息失败:", err)
		return nil, fmt.Errorf("查询用户详细信息失败: %v", err)
	}

	// 检查用户是否存在
	if result.ID == 0 {
		g.Log().Warning(ctx, "⚠️ 用户不存在，邮箱:", email)
		return nil, fmt.Errorf("用户不存在")
	}

	// 转换为响应格式
	userDetail := &model.UserDetailResponse{
		ID:           result.ID,
		Email:        result.Email,
		FullName:     result.FullName,
		Role:         result.Role,
		Phone:        result.Phone,
		Age:          result.Age,
		DID:          result.DID,
		BuildingName: result.BuildingName,
		BuildingAddr: result.BuildingAddr,
		BuildingType: result.BuildingType,
		PropertyName: result.PropertyName,
		Occupation:   result.Occupation,
		Institution:  result.Institution,
		Status:       result.Status,
		CreatedAt:    result.CreatedAt,
		UpdatedAt:    result.CreatedAt, // DID模型没有UpdatedAt字段，使用CreatedAt
	}

	g.Log().Info(ctx, "✅ 成功获取用户详细信息，邮箱:", email)
	return userDetail, nil
}
