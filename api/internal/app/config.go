package app

import (
	"embed"
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
func NewConfig(embedFs embed.FS) *EnvConfig {
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
