package saga

import (
	"context"
	"log"
)

// Activity contained a unit of activity which is part of a workflow
type AsyncActivity[T any] struct {
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
func NewAsyncActivity[T any](name string, commit execute[T], rollback execute[T]) *AsyncActivity[T] {
	return &AsyncActivity[T]{
		name:     name,
		commit:   commit,
		rollback: rollback,
	}
}

func (aa *AsyncActivity[T]) GetName() string {
	return aa.name
}

// Commit executes activity
func (aa *AsyncActivity[T]) Commit(ctx context.Context, param T) error {
	return aa.commit(ctx, param)
}

// Rollback rolling back what was committed by current activity
func (aa *AsyncActivity[T]) Rollback(ctx context.Context, param T) error {
	return aa.rollback(ctx, param)
}

// AsyncWorkflow a group of activities to be executed
type AsyncWorkflow[T any] struct {
	// name refer to workflow's name
	name string

	// activities refer to an ordered one or
	// more activity of a workflow
	activities []IActivity[T]

	taskTChs []chan T

	sagaProvider ISagaProvider
}

type ISagaProvider interface {
	Register(name string, taskCh chan string)
}

// NewAsyncWorkflow create a workflow
func NewAsyncWorkflow[T any](sagaProvider ISagaProvider, name string, activities ...IActivity[T]) *AsyncWorkflow[T] {
	return &AsyncWorkflow[T]{
		name:         name,
		activities:   activities,
		taskTChs:     make([]chan T, len(activities)),
		sagaProvider: sagaProvider,
	}
}

func (aw *AsyncWorkflow[T]) Commit(ctx context.Context, param T) {
	log.Printf("committing!")
	aw.taskTChs[0] <- param
}

func (aw *AsyncWorkflow[T]) Start() {
	// construct concurrent worker
	for i, activity := range aw.activities {
		// initiate channel
		aw.taskTChs[i] = make(chan T, 10)

		// start concurrent worker for activity at 'i'
		go func(step int, taskCh chan T, _activity IActivity[T]) {
			for task := range taskCh {
				// handle activity at 'i'
				_activity.Commit(context.Background(), task)

				// when there is next available activity
				// then push task to next channel
				if i+1 < len(aw.taskTChs) {
					aw.taskTChs[i+1] <- task
				}
			}
		}(i, aw.taskTChs[i], activity)
	}
}

func (aw *AsyncWorkflow[T]) Stop() {
	for i := range aw.taskTChs {
		close(aw.taskTChs[i])
	}
}
