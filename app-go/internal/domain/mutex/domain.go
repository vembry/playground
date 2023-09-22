package mutex

import (
	"app/internal/model"
	"context"
	"fmt"

	redsync "github.com/go-redsync/redsync/v4"
	redsyncRedis "github.com/go-redsync/redsync/v4/redis/goredis/v9"

	redis "github.com/redis/go-redis/v9"
)

// mutex is the instance of mutex
type mutex struct {
	client *redsync.Redsync
}

// New is to setup mutex instance
func New(client redis.UniversalClient) *mutex {
	redisPool := redsyncRedis.NewPool(client)
	return &mutex{
		client: redsync.New(redisPool),
	}
}

// Acquire acquires a new mutex.
func (s *mutex) Acquire(ctx context.Context, name string) (*model.Mutex, error) {
	// setup mutex
	// with-retries = 1 will force red-sync to only attempt to lock once
	mutex := s.client.NewMutex(name, redsync.WithTries(1))

	// get lock
	if err := mutex.LockContext(ctx); err != nil {
		// when theres error, we'll always assume that
		// the lock already been taken
		return &model.Mutex{}, err
	}

	return &model.Mutex{
		Name:   name,
		Locker: mutex,
	}, nil
}

// Delete deletes a mutex.
func (s *mutex) Delete(ctx context.Context, mutex *model.Mutex) error {
	if ok, err := mutex.Locker.UnlockContext(ctx); !ok || err != nil {
		if !ok {
			err = fmt.Errorf("error obtaining redsync quorum while attempting to delete mutex. name=%s. mutexName=%s", mutex.Name, mutex.Locker.Name())
		}
		return err
	}

	return nil
}
