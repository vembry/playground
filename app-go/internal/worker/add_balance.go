package worker

import (
	"app/internal/model"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/hibiken/asynq"
)

// addBalance is the instance of topup worker
type addBalance struct {
	worker  workerHandler
	balance balanceProvider
}

// balanceProvider is the expected spec of balance
type balanceProvider interface {
	Add(ctx context.Context, in *model.AddBalanceParam) error
}

// NewTopup is to initialize topup worker
func NewAddBalance(balance balanceProvider) *addBalance {
	return &addBalance{
		balance: balance,
	}
}

// Type is to return the type of pending-transaction worker
func (ab *addBalance) Type() string {
	return "addBalance"
}

// Queue is to return the queue of pending-transaction worker
func (ab *addBalance) Queue() string {
	return "addBalance"
}

// Perform is the consumed task's handler
func (ab *addBalance) Perform(ctx context.Context, task *asynq.Task) error {
	var in model.AddBalanceParam
	err := json.Unmarshal(task.Payload(), &in)
	if err != nil {
		log.Printf("found error on processing task. payload=%s. err=%v", string(task.Payload()), err)
	}

	// process task
	err = ab.balance.Add(ctx, &in)
	if err != nil {
		if errors.Is(err, model.ErrBalanceLocked) {
			// if getting lock-balance error, then requeue
			err = ab.Enqueue(ctx, &in)
			if err != nil {
				log.Printf("found error on enqueuing add-balance task. payload=%s. err=%v", string(task.Payload()), err)
			}
		} else {
			return err
		}
	}
	return nil
}

// WithWorker is to inject worker instance to topup-worker
func (ab *addBalance) WithWorker(worker workerHandler) {
	ab.worker = worker
}

// Enqueue is the add-balance's handler prior to enqueue task to main worker
func (ab *addBalance) Enqueue(ctx context.Context, in *model.AddBalanceParam) error {
	raw, _ := json.Marshal(in)
	task := asynq.NewTask(ab.Type(), raw)
	_, err := ab.worker.Enqueue(ctx, task, asynq.Queue(ab.Type()))
	return err
}
