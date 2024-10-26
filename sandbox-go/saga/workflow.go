package saga

import (
	"context"
	"log"
)

type execute[T any] func(context.Context, T) error

// IActivity defines specification of workflow's activity
type IActivity[T any] interface {
	GetName() string
	Commit(context.Context, T) error
	Rollback(context.Context, T) error
}

// Workflow a group of activities to be executed
type Workflow[T any] struct {
	// name refer to workflow's name
	name string

	// activities refer to an ordered one or
	// more activity of a workflow
	activities []IActivity[T]
}

// NewWorkflow create a workflow
func NewWorkflow[T any](name string, activities ...IActivity[T]) *Workflow[T] {
	return &Workflow[T]{
		name:       name,
		activities: activities,
	}
}

// Commit executes the workflow
func (w *Workflow[T]) Commit(ctx context.Context, param T) {
	log.Printf("starting workflow='%s'", w.name)

	fallbacks := []execute[T]{}
	isFallingBack := false

	// executes activity
	for _, activity := range w.activities {
		err := activity.Commit(ctx, param)

		// when error found,
		// then workflow will fallback
		// and break activities loop
		if err != nil {
			isFallingBack = true
			break
		}

		fallbacks = append(fallbacks, activity.Rollback)
	}

	if isFallingBack {
		// when falling back, we need to
		// start from the latest activity committed

		for i := len(fallbacks) - 1; i >= 0; i-- {
			// get latest fallback entry
			fallback := fallbacks[i]

			// execute fallback
			fallback(ctx, param)
		}
	}

	log.Printf("closing workflow='%s'", w.name)
	log.Println(param)
}
