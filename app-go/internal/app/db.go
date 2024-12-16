package app

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// newOrmDb is to initialize DB in ORM form using gorm
func NewOrmDb(cfg *EnvConfig) (*gorm.DB, func()) {
	// dsn := "host=localhost user=local password=local dbname=playground_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error on initiating db. err=%v", err)
	}

	// add otel plugin
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Fatalf("error on installing tracing into gorm. err=%v", err)
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
