package events

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

type UnityEventReceiver func (event UnityEvent)
type UnityCommandReceiver func (event UnityCommand)

type UnityEventSubscriber interface {
	Subscriber() error
	AddHandler(handler UnityEventReceiver)
}

type UnityCommandSubscriber interface {
	Receiver() error
	AddHandler(handler UnityCommandReceiver)
}

type Subscriber struct {
	handlers []UnityEventReceiver
	queue amqp.Queue
	ch *amqp.Channel
	config Config
}

type CommandReceiver struct {
	handlers []UnityCommandReceiver
	queue amqp.Queue
	ch *amqp.Channel
	config Config
}

func NewSubscriber(cfg Config) (*Subscriber, error) {
	e := Subscriber{ config: cfg}
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

	events, err := e.ch.Consume(e.config.RabbitQueue, e.config.ConsumerName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
	}

	go func() {
		for ev := range events {
			// ... this consumer is responsible for sending pages per log
			var queueEvent UnityEvent
			json.Unmarshal(ev.Body, &queueEvent)
			if e.handlers != nil && len(e.handlers) > 0 {
				for h := range e.handlers {
					go e.handlers[h](queueEvent)
				}
			}
		}
	}()

	return &e, nil
}
func NewCommandReceiver(cfg Config) (*CommandReceiver, error) {
	e := CommandReceiver{ config: cfg}
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

	events, err := e.ch.Consume(e.config.RabbitQueue, e.config.ConsumerName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("basic.consume: %v", err)
	}

	go func() {
		for ev := range events {
			// ... this consumer is responsible for sending pages per log
			var queueCmd UnityCommand
			json.Unmarshal(ev.Body, &queueCmd)
			if e.handlers != nil && len(e.handlers) > 0 {
				for h := range e.handlers {
					go e.handlers[h](queueCmd)
				}
			}
		}
	}()

	return &e, nil
}
func (s *Subscriber) AddHandler(handler UnityEventReceiver) {
	if s.handlers == nil {
		s.handlers = make([]UnityEventReceiver, 0)
	}

	s.handlers = append(s.handlers, handler)
}
func (c *CommandReceiver) AddHandler(handler UnityCommandReceiver) {
	if c.handlers == nil {
		c.handlers = make([]UnityCommandReceiver, 0)
	}

	c.handlers = append(c.handlers, handler)
}
