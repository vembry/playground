package locker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
)

// locker contains implementation related to locking capabilities
type locker struct {
	redis *redsync.Redsync
}

func New(client *redis.Client) *locker {
	return &locker{
		redis: redsync.New(goredis.NewPool(client)),
	}
}

func (l *locker) Lock(ctx context.Context, key string) (func(context.Context), error) {
	// construct lock
	mutex := l.redis.NewMutex(key,
		redsync.WithTries(1),
		redsync.WithExpiry(time.Minute),
	)
	if mutex == nil {
		return nil, fmt.Errorf("constructing mutex return nil")
	}

	// locking
	err := mutex.LockContext(ctx)
	if err != nil {
		// TODO:
		// need to skim thru the errors and adapt them
		// to this locker module dedicated errors
		return func(ctx context.Context) {}, err
	}

	return func(ctx context.Context) {
		mutex.UnlockContext(ctx)
	}, nil
}
