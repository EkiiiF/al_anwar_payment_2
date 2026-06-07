package config

import (
	"fmt"
	"log"

	"github.com/EkiiiF/al_anwar_payment_2.git/internal/model/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection pool: %v", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	log.Println("Running database migrations...")
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	db.Exec("DROP INDEX IF EXISTS idx_student_primary ON addresses")

	log.Println("Database connected and migrated successfully")
	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.Role{},
		&domain.StudentStatusType{},
		&domain.Category{},
		&domain.User{},
		&domain.Student{},
		&domain.Guardian{},
		&domain.Address{},
		&domain.Invoice{},
		&domain.Payment{},
		&domain.Notification{},
		&domain.ActivityLog{},
		&domain.TokenBlacklist{},
		&domain.Setting{},
	)
}
