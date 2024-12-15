package app

import (
	"context"
	"embed"
	"log"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

// EnvConfig is the instance to compile all env vars
type EnvConfig struct {
	HttpAddress string `env:"HTTP_ADDRESS"`
	DBConn      string `env:"DB_CONN"`
	RedisUri    string `env:"REDIS_URI"`
	KafkaBroker string `env:"KAFKA_BROKER"`
	RabbitUri   string `env:"RABBIT_URI"`
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

	// read all required env values
	var cfg EnvConfig
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatalf("failed to read environment variables. err=%v", err)
	}

	return &cfg
}
