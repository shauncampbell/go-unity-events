package events

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type UnityEventPublisher interface {
	Publish(event UnityEvent) error
}

type Events struct {
	queue amqp.Queue
	ch *amqp.Channel
	config Config
	UnityEventPublisher
}

type Config struct {
	RabbitHost string
	RabbitPort int
	RabbitUser string
	RabbitPass string
	RabbitQueue string
}

func New(cfg Config) (*Events, error) {
	e := Events{ config: cfg}

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitUser, cfg.RabbitPass, cfg.RabbitHost, cfg.RabbitPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	e.ch, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	e.queue, err = e.ch.QueueDeclare(
		cfg.RabbitQueue, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &e, nil
}

func (e *Events) Publish(event UnityEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = e.ch.Publish(
		"",     // exchange
		e.config.RabbitQueue, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "application/json",
			Body: body,
		})

	return err
}