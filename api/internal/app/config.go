package app

import (
	"embed"
	"log"
	"os"
	"strings"
)

type EnvConfig struct {
	HttpAddress string
	DBConn      string
	RedisUri    string
}

// NewConfig is to parse env
func NewConfig(embedFs embed.FS) *EnvConfig {
	// load file from embed
	envs, err := embedFs.ReadFile("configs/local.env")
	if err != nil {
		log.Fatalf("failed to load .env file from embed.Fs. err=%v", err)
	}

	// load envs to runtime(?) line by line
	lines := strings.Split(string(envs), "\n")
	for _, line := range lines {
		if line != "" {
			splits := strings.SplitN(line, "=", 2)

			// skip env if its defined already
			if os.Getenv(splits[0]) != "" {
				continue
			}
			if err := os.Setenv(splits[0], splits[1]); err != nil {
				log.Fatalf("failed to inject .env values. env=%s. value%s. err=%v", splits[0], splits[1], err)
			}
		}
	}

	return &EnvConfig{
		HttpAddress: os.Getenv("HTTP_ADDRESS"),
		DBConn:      os.Getenv("DB_CONN"),
		RedisUri:    os.Getenv("REDIS_URI"),
	}
}
