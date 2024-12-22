package rabbit

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbit struct {
	conn      *amqp.Connection
	consumers []iConsumer
}

type iConsumer interface {
	Channel() *amqp.Channel
	Name() string
	Handle(context.Context, amqp.Delivery) error
}

func New(rabbitUri string) *rabbit {
	// setup connection
	conn, err := amqp.Dial(rabbitUri)
	if err != nil {
		log.Fatalf("rabbit: failed to dial rabbit. err=%v", err)
	}

	return &rabbit{
		conn: conn,
	}
}

func (r *rabbit) GetConnection() *amqp.Connection {
	return r.conn
}

func (r *rabbit) RegisterWorkers(consumers ...iConsumer) {
	r.consumers = consumers
}

func (r *rabbit) Name() string {
	return "rabbit"
}

func (r *rabbit) Start() {
	for _, consumer := range r.consumers {
		log.Printf("rabbit: starting '%s' consumer", consumer.Name())

		r.startConsumer(consumer)
	}
}

// startConsumer is to start consumer by prepping needed preps
func (r *rabbit) startConsumer(consumer iConsumer) {
	ch := consumer.Channel()

	// declare queue in case it is missing
	// for now the config will be defined here
	// until further cases
	q, err := ch.QueueDeclare(
		consumer.Name(), // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Fatalf("rabbit: failed to declare a queue. consumer=%s. err=%v", consumer.Name(), err)
	}

	// is this proper?
	messageCh, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("rabbit: failed to setup message consumer. consumer=%s. err=%v", consumer.Name(), err)
	}

	// start consuming
	go r.consume(q, consumer.Handle, messageCh)
}

// consume polls messages from message-channel
func (r *rabbit) consume(queue amqp.Queue, handler func(context.Context, amqp.Delivery) error, messageCh <-chan amqp.Delivery) {
	for message := range messageCh {
		log.Printf("'%s' consuming %s", queue.Name, string(message.Body))
		r.consumeMessage(queue, message, handler)
	}
}

// consumeMessage consume the actual message
func (r *rabbit) consumeMessage(queue amqp.Queue, message amqp.Delivery, handler func(context.Context, amqp.Delivery) error) {
	ctx := context.Background()

	// consume incoming message
	err := handler(ctx, message)
	if err != nil {
		log.Printf("rabbit: failed to handle. consumer=%s. err=%v", queue.Name, err)

		err = message.Reject(true)
		if err != nil {
			log.Printf("rabbit: failed to reject. consumer=%s. err=%v", queue.Name, err)
		}
	} else {
		message.Ack(true)
	}
}

func (r *rabbit) Stop() {
	// close channels
	for _, consumer := range r.consumers {
		err := consumer.Channel().Close()
		if err != nil {
			log.Printf("rabbit: error trying to close channel for '%s'", consumer.Name())
		}
	}

	r.conn.Close() // close connection
}
