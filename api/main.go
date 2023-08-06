package main

import (
	"embed"
	"log"
	"os"
	"os/signal"
	"syscall"

	"api/internal/app"
	"api/internal/app/handler"
	ledgerDomain "api/internal/domain/ledger"
	transactionDomain "api/internal/domain/transaction"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	//go:embed configs
	embedFs embed.FS
)

func main() {
	// setup config
	config := app.NewConfig(embedFs)

	// setup db
	db := newOrmDb(config)

	// close db connection on app closure
	defer func() {
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("found error on getting DB. err=%v", err)
		}

		err = sqlDb.Close()
		if err != nil {
			log.Fatalf("found error on closing DB connection. err=%v", err)
		}
	}()

	// setup ledger domain
	ledger := ledgerDomain.New(db)

	// setup transaction domain
	transaction := transactionDomain.New(db)

	// setup server's http-handler
	r := handler.NewHttpHandler(transaction, ledger)

	// start server
	serve(config, r)

	// awaits for interrupt signals
	watchForExitSignal()

	log.Printf("stopping http server...")

	// do shutdown handling
	// ...

	log.Printf("server stopped")
}

// serve is to start the server
func serve(cfg *app.EnvConfig, r *gin.Engine) {
	// start server
	log.Printf("starting http server...")
	go func() {
		if err := r.Run(cfg.HttpAddress); err != nil {
			log.Fatalf("gin stopped running. err=%v", err)
		}
	}()
}

// watchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func watchForExitSignal() os.Signal {
	ch := make(chan os.Signal, 4)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	return <-ch
}

// newOrmDb is to initialize DB in ORM form using gorm
func newOrmDb(cfg *app.EnvConfig) *gorm.DB {
	// dsn := "host=localhost user=local password=local dbname=credits port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initiate db. err=%v", err)
	}

	return db
}
