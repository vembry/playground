package app

import (
	"log"

	"github.com/redis/go-redis/v9"
)

type cache struct {
	client *redis.Client
}

// NewCache is to initiate cache instance
func NewCache(cfg *EnvConfig) *cache {
	opt, err := redis.ParseURL(cfg.RedisUri)
	if err != nil {
		log.Fatalf("found error on parsing redis url. err=%v", err)
	}

	return &cache{
		client: redis.NewClient(opt),
	}
}

// GetClient is to retrieve redis client
func (ch *cache) GetClient() redis.UniversalClient {
	return ch.client
}
