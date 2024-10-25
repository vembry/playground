package saga

import (
	"log"
)

type saga struct {
}

func New() *saga {
	return &saga{}
}

type Workflow[T any] struct {
	name       string
	activities []Activity[T]
}

type Activity[T any] interface {
	Commit(T) error
	Rollback(T)
}

// CreateWorkflow create a workflow group
func CreateWorkflow[T any](name string, activities ...Activity[T]) *Workflow[T] {
	return &Workflow[T]{
		name:       name,
		activities: activities,
	}
}

func (w *Workflow[T]) Commit(param T) {
	fallbacks := []func(T){}

	isFallingBack := false
	for _, activity := range w.activities {
		err := activity.Commit(param)
		if err != nil {
			isFallingBack = true
			break
		}

		fallbacks = append([]func(T){activity.Rollback}, fallbacks...)
	}

	if isFallingBack {
		for _, fallback := range fallbacks {
			fallback(param)
		}
	}

	log.Printf("closing workflow='%s'", w.name)
	log.Println(param)
}
