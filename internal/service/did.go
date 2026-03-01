package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fabric-sdk/internal/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateDID 新增DID登记，返回did身份ID
func CreateDID(name, phone, email, role, password string, age int) (string, error) {
	if DB == nil {
		return "", fmt.Errorf("数据库未初始化")
	}
	hash := sha256.Sum256([]byte(password))
	didStr := "did:" + uuid.NewString()
	beijing := time.Now().In(time.FixedZone("CST", 8*3600))
	did := model.DID{
		Name:      name,
		Phone:     phone,
		Email:     email,
		Role:      role,
		Age:       age,
		Password:  hex.EncodeToString(hash[:]),
		CreatedAt: beijing.Format("2006-01-02 15:04:05"),
		DID:       didStr,
	}
	if err := DB.Create(&did).Error; err != nil {
		return "", err
	}
	return did.DID, nil
}

// VerifyDID 通过did和密码校验
func VerifyDID(didStr, password string) (*model.DID, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	var did model.DID
	if err := DB.Where("did = ?", didStr).First(&did).Error; err != nil {
		return nil, fmt.Errorf("DID不存在")
	}
	hash := sha256.Sum256([]byte(password))
	if did.Password != hex.EncodeToString(hash[:]) {
		return nil, fmt.Errorf("密码错误")
	}
	return &did, nil
}
