package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"time"
)

type resource struct {
	name string
}

func (r *resource) SomeFunctionToCommit(ctx context.Context, param *someparameter) error {
	param.task = append(param.task, fmt.Sprintf("task-%s", r.name))
	log.Printf("committing resource=%s", r.name)
	duration := time.Duration(rand.IntN(100)) * time.Millisecond
	time.Sleep(duration)
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
