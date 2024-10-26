package saga

import "context"

// Activity contained a unit of activity which is part of a workflow
type Activity[T any] struct {
	// name refer to the activity name
	name string

	// commit refer to the execution we want
	// to achieve in the current activity
	commit execute[T]

	// rollback refer to rolling back process
	// which already executed in the current acitivity
	rollback execute[T]
}

// NewActivity creates a new unit of activity
func NewActivity[T any](name string, commit execute[T], rollback execute[T]) *Activity[T] {
	return &Activity[T]{
		name:     name,
		commit:   commit,
		rollback: rollback,
	}
}

func (a *Activity[T]) GetName() string {
	return a.name
}

// Commit executes activity
func (a *Activity[T]) Commit(ctx context.Context, param T) error {
	return a.commit(ctx, param)
}

// Rollback rolling back what was committed by current activity
func (a *Activity[T]) Rollback(ctx context.Context, param T) error {
	return a.rollback(ctx, param)
}
