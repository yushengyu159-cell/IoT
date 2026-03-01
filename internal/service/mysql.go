package service

import (
	"fabric-sdk/internal/model"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return DB
}

func InitMySQL() error {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:Test@123456@tcp(127.0.0.1:3306)/esg?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("MySQL连接失败: %v", err)
	}
	db = db.Debug() // 启用GORM Debug日志
	
	// 安全地删除表并重新创建
	err = safeDropAndMigrate(db)
	if err != nil {
		return fmt.Errorf("MySQL自动迁移失败: %v", err)
	}
	
	DB = db
	return nil
}

// safeDropAndMigrate 安全地删除表并重新迁移
func safeDropAndMigrate(db *gorm.DB) error {
	// 先尝试正常迁移
	err := db.AutoMigrate(&model.ESGFile{}, &model.DID{}, &model.FailedRegistration{})
	if err == nil {
		return nil
	}
	
	// 如果迁移失败，检查是否是外键约束错误
	errStr := err.Error()
	if isForeignKeyError(errStr) {
		// 安全地删除外键约束
		safeDropForeignKeys(db)
		// 删除表并重新创建
		db.Exec("DROP TABLE IF EXISTS esg_files")
		db.Exec("DROP TABLE IF EXISTS dids")
		db.Exec("DROP TABLE IF EXISTS failed_registrations")

		// 重新迁移
		err = db.AutoMigrate(&model.ESGFile{}, &model.DID{}, &model.FailedRegistration{})
		if err != nil {
			return fmt.Errorf("重新迁移失败: %v", err)
		}
		return nil
	}
	
	return err
}

// isForeignKeyError 检查是否是外键相关的错误
func isForeignKeyError(errStr string) bool {
	keywords := []string{
		"Error 1091",
		"DROP FOREIGN KEY",
		"DROP INDEX", 
		"check that column/key exists",
		"Can't DROP",
		"doesn't exist",
	}
	
	for _, keyword := range keywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
}

// safeDropForeignKeys 安全地删除外键约束
func safeDropForeignKeys(db *gorm.DB) {
	// 获取esg_files表的外键约束
	var constraints []struct {
		ConstraintName string `gorm:"column:CONSTRAINT_NAME"`
	}
	
	db.Raw(`
		SELECT CONSTRAINT_NAME 
		FROM information_schema.KEY_COLUMN_USAGE 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'esg_files' 
		AND REFERENCED_TABLE_NAME IS NOT NULL
	`).Scan(&constraints)
	
	// 安全地删除每个外键约束
	for _, constraint := range constraints {
		if constraint.ConstraintName != "" {
			db.Exec(fmt.Sprintf("ALTER TABLE esg_files DROP FOREIGN KEY IF EXISTS %s", constraint.ConstraintName))
		}
	}
	
	// 同样处理dids表
	db.Raw(`
		SELECT CONSTRAINT_NAME 
		FROM information_schema.KEY_COLUMN_USAGE 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'dids' 
		AND REFERENCED_TABLE_NAME IS NOT NULL
	`).Scan(&constraints)
	
	for _, constraint := range constraints {
		if constraint.ConstraintName != "" {
			db.Exec(fmt.Sprintf("ALTER TABLE dids DROP FOREIGN KEY IF EXISTS %s", constraint.ConstraintName))
		}
	}
}

// contains 是 strings.Contains 的简写
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
