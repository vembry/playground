package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sandbox/saga"
	"sandbox/saga/provider"
	"syscall"
	"time"
)

type someparameter struct {
	task []string
}

func main() {
	// construct mock activities
	a := &resource{name: "a"}
	b := &resource{name: "b"}
	c := &resource{name: "c"}
	d := &failingresource{name: "d"}

	// create sync workflow
	// ====================

	// construct workflow
	workflow1 := saga.NewWorkflow(
		"workflow-1",
		saga.NewActivity("activity-a", a.SomeFunctionToCommit, a.SomeFunctionToRollback),
		saga.NewActivity("activity-b", b.SomeFunctionToCommit, b.SomeFunctionToRollback),
		saga.NewActivity("activity-c", c.SomeFunctionToCommit, c.SomeFunctionToRollback),
		saga.NewActivity("activity-d", d.AnotherFunctionToCommit, d.AnotherFunctionToRollback),
	)
	workflow2 := saga.NewWorkflow(
		"workflow-2",
		saga.NewActivity("activity-a", a.SomeFunctionToCommit, a.SomeFunctionToRollback),
		saga.NewActivity("activity-b", b.SomeFunctionToCommit, b.SomeFunctionToRollback),
		saga.NewActivity("activity-c", c.SomeFunctionToCommit, c.SomeFunctionToRollback),
	)

	// execute sync workflow
	workflow1.Commit(context.Background(), &someparameter{})
	workflow2.Commit(context.Background(), &someparameter{})

	// construct async workflow
	// ========================

	// construct async client
	// using hibiken/asynq redis-worker for easy setup
	asynqSaga := provider.NewAsynqProvider()
	defer asynqSaga.Close()

	// construct workflow
	asycnworkflow1 := saga.NewAsyncWorkflow(
		asynqSaga,
		"asynchronous-workflow-1",
		saga.NewActivity("activity-a", a.SomeFunctionToCommit, a.SomeFunctionToRollback),
		saga.NewActivity("activity-b", b.SomeFunctionToCommit, b.SomeFunctionToRollback),
		saga.NewActivity("activity-c", c.SomeFunctionToCommit, c.SomeFunctionToRollback),
	)

	asycnworkflow1.Start()
	defer asycnworkflow1.Stop()

	// asynqSaga.Start()

	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		go func() {
			asycnworkflow1.Commit(context.Background(), &someparameter{})
		}()
		time.Sleep(100 * time.Millisecond)
	}

	WatchForExitSignal()
}

// WatchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func WatchForExitSignal() os.Signal {
	log.Printf("awaiting sigterm...")
	ch := make(chan os.Signal, 4)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	return <-ch
}
