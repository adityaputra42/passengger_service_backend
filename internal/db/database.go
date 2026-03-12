package db

import (
	"fmt"
	"log"
	"passenger_service_backend/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: false,
		PrepareStmt:                              true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Connection pool tuning
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	log.Println("Database connected successfully!")
	return nil
}


func GetDBStats() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting DB stats: %v", err)
		return
	}
	s := sqlDB.Stats()
	log.Printf(
		"DB Stats — Open: %d | In Use: %d | Idle: %d | WaitCount: %d | WaitDuration: %v | MaxIdleClosed: %d | MaxLifetimeClosed: %d",
		s.OpenConnections, s.InUse, s.Idle,
		s.WaitCount, s.WaitDuration,
		s.MaxIdleClosed, s.MaxLifetimeClosed,
	)
}

func GetDB() *gorm.DB {
	return DB
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error closing DB: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing DB connection: %v", err)
	}
}
