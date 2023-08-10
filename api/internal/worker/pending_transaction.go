package worker

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/segmentio/ksuid"
)

// pendingTransaction is the instance of pending-transaction worker
type pendingTransaction struct {
	worker      workerHandler
	transaction transactionProvider
}

// workerHandler is the expected spec of a worker. This is to give the individual worker
// the ability to enqueue task
type workerHandler interface {
	Enqueue(ctx context.Context, task *asynq.Task, taskOptions ...asynq.Option) (*asynq.TaskInfo, error)
}

// transactionProvider is the expected spec of transaction
type transactionProvider interface {
	ProcessPending(ctx context.Context, transactionId ksuid.KSUID) error
}

// NewPendingTransaction is to initialize pending transaction worker
func NewPendingTransaction(transaction transactionProvider) *pendingTransaction {
	return &pendingTransaction{
		transaction: transaction,
	}
}

// Type is to return the type of pending-transaction worker
func (tp *pendingTransaction) Type() string {
	return "pendingTransaction"
}

// Queue is to return the queue of pending-transaction worker
func (tp *pendingTransaction) Queue() string {
	return "pendingTransaction"
}

// Perform is the consumed task's handler
func (tp *pendingTransaction) Perform(ctx context.Context, task *asynq.Task) error {
	transactionId := ksuid.KSUID{}
	if err := transactionId.UnmarshalText(task.Payload()); err != nil {
		log.Printf("found error on converting task's payload. payload=%s. err=%v", string(task.Payload()), err)
		return nil
	}

	return tp.transaction.ProcessPending(ctx, transactionId)
}

// WithWorker is to inject worker to pending-transaction-worker
func (tp *pendingTransaction) WithWorker(worker workerHandler) {
	tp.worker = worker
}

// Enqueue is the pending-transaction's handler prior to enqueue task to main worker
func (tp *pendingTransaction) Enqueue(ctx context.Context, transactionId ksuid.KSUID) error {
	task := asynq.NewTask(tp.Type(), []byte(transactionId.String()))
	_, err := tp.worker.Enqueue(ctx, task, asynq.Queue(tp.Type()))
	return err
}
