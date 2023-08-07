package app

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newOrmDb is to initialize DB in ORM form using gorm
func NewOrmDb(cfg *EnvConfig) *gorm.DB {
	// dsn := "host=localhost user=local password=local dbname=credits port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initiate db. err=%v", err)
	}

	return db
}
