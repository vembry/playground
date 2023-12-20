package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

type locker struct {
	c *redsync.Redsync
}

func NewLocker(client *redis.Client) *locker {
	return &locker{
		c: redsync.New(goredis.NewPool(client)),
	}
}

func (l *locker) AcquireLock(ctx context.Context, key string) (func(_ctx context.Context), error) {
	mutex := l.c.NewMutex(key,
		redsync.WithTries(1),
		redsync.WithExpiry(time.Hour),
	)
	if mutex == nil {
		return nil, fmt.Errorf("constructing mutex return nil")
	}

	err := mutex.LockContext(ctx)
	if err != nil {
		log.Printf("got error on locking. err=%v", err)
		return nil, err
	}

	return func(_ctx context.Context) {
		mutex.UnlockContext(_ctx)
	}, nil
}
