package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"product-challenge/internal/models"
	"product-challenge/pkg/config"
)

func NewCockroachDB(cfg *config.Config) (*gorm.DB, error) {
	// data source name config. CockroachDB can use postgres diver
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Create sequence if not exists
	err = db.Exec(`CREATE SEQUENCE IF NOT EXISTS order_number_seq START 1`).Error
	if err != nil {
		return nil, err
	}
	
	// Auto Migrate the schema
	err = db.AutoMigrate(
		&models.Products{},
		&models.User{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return db, nil
}
