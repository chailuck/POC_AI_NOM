// internal/database/database.go
package database

import (
	"fmt"

	"github.com/your-username/tmf632-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schemas
	err = autoMigrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Individual{},
		&models.ContactMedium{},
		&models.ExternalReference{},
		&models.IndividualIdentification{},
		&models.PartyCharacteristic{},
	)
}
