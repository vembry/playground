package main

import (
	"context"
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

	a := &resource{name: "a"}        // mock activity provider
	b := &resource{name: "b"}        // mock activity provider
	c := &resource{name: "c"}        // mock activity provider
	d := &failingresource{name: "d"} // mock activity provider

	// create workflow
	workflow1 := saga.NewWorkflow(
		"workflow-1",
		saga.NewActivity("activity-a", a.SomeFunctionToCommit, a.SomeFunctionToRollback),
		saga.NewActivity("activity-b", b.SomeFunctionToCommit, b.SomeFunctionToRollback),
		saga.NewActivity("activity-c", c.SomeFunctionToCommit, c.SomeFunctionToRollback),
		saga.NewActivity("activity-d", d.AnotherFunctionToCommit, d.AnotherFunctionToRollback),
	)

	// execute workflow
	workflow1.Commit(context.Background(), &someparameter{})
}

type resource struct {
	name string
}

func (r *resource) SomeFunctionToCommit(ctx context.Context, param *someparameter) error {
	param.task = append(param.task, fmt.Sprintf("task-%s", r.name))
	log.Printf("committing resource=%s", r.name)
	return nil
}

func (r *resource) SomeFunctionToRollback(ctx context.Context, param *someparameter) error {
	log.Printf("rolling back resource=%s", r.name)
	return nil
}

type failingresource struct {
	name string
}

func (r *failingresource) AnotherFunctionToCommit(ctx context.Context, param *someparameter) error {
	return errors.New("lol fail")
}

func (r *failingresource) AnotherFunctionToRollback(ctx context.Context, param *someparameter) error {
	log.Printf("rolling back failingresource=%s", r.name)
	return nil
}
