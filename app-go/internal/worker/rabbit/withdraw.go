package rabbit

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/ksuid"
)

type withdraw struct {
	ch                *amqp.Channel
	withdrawProcessor IWithdrawProcessor
}

func NewWithdraw(conn *amqp.Connection) *withdraw {
	// setup channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. err=%v", err)
	}

	return &withdraw{
		ch: ch,
	}
}

func (t *withdraw) Name() string {
	return "withdraw"
}
func (t *withdraw) Channel() *amqp.Channel {
	return t.ch
}

func (t *withdraw) Produce(ctx context.Context, withdrawId ksuid.KSUID) error {
	body, _ := withdrawId.MarshalText()

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

type IWithdrawProcessor interface {
	ProcessWithdraw(ctx context.Context, withdrawId ksuid.KSUID) error
}

func (t *withdraw) InjectDeps(withdrawProcessor IWithdrawProcessor) {
	t.withdrawProcessor = withdrawProcessor
}

func (t *withdraw) Handle(ctx context.Context, task amqp.Delivery) error {
	id, _ := ksuid.Parse(string(task.Body))
	return t.withdrawProcessor.ProcessWithdraw(ctx, id)
}
