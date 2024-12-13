package app

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newOrmDb is to initialize DB in ORM form using gorm
func NewOrmDb(cfg *EnvConfig) (*gorm.DB, func()) {
	// dsn := "host=localhost user=local password=local dbname=playground_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initiate db. err=%v", err)
	}

	return db, func() {
		// close db connection
		// put this as deferer
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("found error on getting DB. err=%v", err)
		}

		err = sqlDb.Close()
		if err != nil {
			log.Fatalf("found error on closing DB connection. err=%v", err)
		}
	}
}
