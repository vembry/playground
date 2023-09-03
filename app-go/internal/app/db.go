package app

import (
	"app-go/internal/model"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newOrmDb is to initialize DB in ORM form using gorm
func NewOrmDb(cfg *EnvConfig) (*model.DB, func()) {
	// dsn := "host=localhost user=local password=local dbname=credits port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initiate connection to master db. err=%v", err)
	}

	slaveDb, err := gorm.Open(postgres.Open(cfg.SlaveDBConn), &gorm.Config{})
	if err != nil {
		log.Printf("failed to initiate connection to slave db. err=%v", err)

		// overwrite slave with master
		slaveDb = db
	}

	return &model.DB{
			Master: db,
			Slave:  slaveDb,
		}, func() {
			// close db connection
			// put this as deferer

			// close master db connection
			sqlDb, err := db.DB()
			if err != nil {
				log.Fatalf("found error on getting DB. err=%v", err)
			}

			err = sqlDb.Close()
			if err != nil {
				log.Fatalf("found error on closing DB connection. err=%v", err)
			}

			// close slave db connection
			sqlDb, err = slaveDb.DB()
			if err != nil {
				log.Fatalf("found error on getting slave DB. err=%v", err)
			}

			err = sqlDb.Close()
			if err != nil {
				log.Fatalf("found error on closing slave DB connection. err=%v", err)
			}
		}
}
