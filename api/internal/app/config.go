package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	HttpAddress   string
	DBConn        string
	SomethingElse string
}

// NewConfig is to parse env
func NewConfig() *EnvConfig {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatalf("failed to load .env. err=%v", err)
	}

	return &EnvConfig{
		HttpAddress:   os.Getenv("HTTP_ADDRESS"),
		DBConn:        os.Getenv("DB_CONN"),
		SomethingElse: os.Getenv("SOMETHING_ELSE"),
	}
}
