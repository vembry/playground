package rabbit

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/ksuid"
)

type transfer struct {
	ch                *amqp.Channel
	transferProcessor ITransferProcessor
}

func NewTransfer() *transfer {
	return &transfer{}
}

func (t *transfer) Name() string {
	return "transfer"
}
func (t *transfer) Produce(ctx context.Context, transferId ksuid.KSUID) error {
	body, _ := transferId.MarshalText()

	// enqueue task
	err := t.ch.PublishWithContext(ctx,
		"",       // exchange
		t.Name(), // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish a queue's task. err=%v", err)
	}
	return nil
}

type ITransferProcessor interface {
	ProcessTransfer(ctx context.Context, withdrawId ksuid.KSUID) error
}

func (t *transfer) InjectDeps(transferProcessor ITransferProcessor) {
	t.transferProcessor = transferProcessor
}

func (t *transfer) Handle(ctx context.Context, task amqp.Delivery) error {
	id, _ := ksuid.Parse(string(task.Body))
	return t.transferProcessor.ProcessTransfer(ctx, id)
}
