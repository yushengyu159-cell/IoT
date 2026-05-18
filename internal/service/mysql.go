package service

import (
	"fabric-sdk/internal/model"
	"fmt"
	"os"
	"strings"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

const (
	MaxIdleConns    = 10
	MaxOpenConns    = 100
	ConnMaxLifetime = 3600
	ConnMaxIdleTime = 600
)

func GetDB() *gorm.DB {
	return DB
}

func InitMySQL() error {
	dsn := getDSN()
	isDev := os.Getenv("ENV") == "development" || os.Getenv("DEBUG") == "true"
	
	logLevel := logger.Silent
	if isDev {
		logLevel = logger.Info
	}
	
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return fmt.Errorf("MySQL连接失败: %v", err)
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	
	sqlDB.SetMaxIdleConns(MaxIdleConns)
	sqlDB.SetMaxOpenConns(MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(ConnMaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(ConnMaxIdleTime) * time.Second)
	
	err = safeDropAndMigrate(db)
	if err != nil {
		return fmt.Errorf("MySQL自动迁移失败: %v", err)
	}
	
	DB = db
	return nil
}

func getDSN() string {
	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		return dsn
	}
	
	// 从环境变量构建DSN
	host := os.Getenv("MYSQL_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("MYSQL_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("MYSQL_USER")
	if user == "" {
		user = "root"
	}
	password := os.Getenv("MYSQL_PASSWORD")
	if password == "" {
		password = "Test@123456"
	}
	database := os.Getenv("MYSQL_DATABASE")
	if database == "" {
		database = "esg"
	}
	
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, database)
}

func PingDatabase() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

func safeDropAndMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&model.ESGFile{}, &model.DID{}, &model.FailedRegistration{}, &model.BuildingAsset{})
	if err == nil {
		return nil
	}
	
	errStr := err.Error()
	if isForeignKeyError(errStr) {
		safeDropForeignKeys(db)
		db.Exec("DROP TABLE IF EXISTS esg_files")
		db.Exec("DROP TABLE IF EXISTS dids")
		db.Exec("DROP TABLE IF EXISTS failed_registrations")
		db.Exec("DROP TABLE IF EXISTS building_assets")

		err = db.AutoMigrate(&model.ESGFile{}, &model.DID{}, &model.FailedRegistration{}, &model.BuildingAsset{})
		if err != nil {
			return fmt.Errorf("重新迁移失败: %v", err)
		}
		return nil
	}
	return err
}

func isForeignKeyError(errStr string) bool {
	keywords := []string{"Error 1091", "DROP FOREIGN KEY", "DROP INDEX", "check that column/key exists", "Can't DROP", "doesn't exist"}
	for _, keyword := range keywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
}

func safeDropForeignKeys(db *gorm.DB) {
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
	
	for _, constraint := range constraints {
		if constraint.ConstraintName != "" {
			db.Exec("ALTER TABLE esg_files DROP FOREIGN KEY ?", constraint.ConstraintName)
		}
	}
	
	db.Raw(`
		SELECT CONSTRAINT_NAME 
		FROM information_schema.KEY_COLUMN_USAGE 
		WHERE TABLE_SCHEMA = DATABASE() 
		AND TABLE_NAME = 'dids' 
		AND REFERENCED_TABLE_NAME IS NOT NULL
	`).Scan(&constraints)
	
	for _, constraint := range constraints {
		if constraint.ConstraintName != "" {
			db.Exec("ALTER TABLE dids DROP FOREIGN KEY ?", constraint.ConstraintName)
		}
	}
}
