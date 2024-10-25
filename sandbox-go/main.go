package main

import (
	"errors"
	"fmt"
	"log"
	"sandbox/saga"
)

type someparameter struct {
	task []string
}

func main() {
	// s := saga.New()

	a := &resource{name: "a"}
	b := &resource{name: "b"}
	c := &resource{name: "c"}
	d := &failingresource{name: "d"}

	workflow1 := saga.CreateWorkflow("workflow-1", a, b, c, d)

	workflow1.Commit(&someparameter{})
}

type resource struct {
	name string
}

func (r *resource) Commit(param *someparameter) error {
	param.task = append(param.task, fmt.Sprintf("task-%s", r.name))
	log.Printf("committing resource=%s", r.name)
	return nil
}

func (r *resource) Rollback(param *someparameter) {
	log.Printf("rolling back resource=%s", r.name)
}

type failingresource struct {
	name string
}

func (r *failingresource) Commit(param *someparameter) error {
	return errors.New("lol fail")
}

func (r *failingresource) Rollback(param *someparameter) {
	log.Printf("rolling back failingresource=%s", r.name)
}
