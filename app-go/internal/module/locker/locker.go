package locker

import "context"

// locker contains implementation related to locking capabilities
type locker struct {
}

func New() *locker {
	return &locker{}
}

func (l *locker) Lock(ctx context.Context, key string) error {
	return nil
}

func (l *locker) Unlock(ctx context.Context, key string) {
}
