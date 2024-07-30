package rabbit

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbit struct{}

func New() *rabbit {
	conn, err := amqp.Dial("amqp://guest:guest@host.docker.internal:5672/")
	if err != nil {
		log.Fatalf("failed to dial rabbit. err=%v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. err=%v", err)
	}

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. err=%v", err)
	}

	body := "Hello World!"
	err = ch.PublishWithContext(context.Background(),
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		log.Fatalf("failed to publish a queue's task. err=%v", err)
	}

	log.Printf(" [x] Sent %s\n", body)

	return &rabbit{}
}
