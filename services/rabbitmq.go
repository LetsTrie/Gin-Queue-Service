package service

import (
	"fmt"
	"log"
	"server/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ wraps a RabbitMQ RabbitMQ and channel.
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// handleNotification is a function that handles a notification message received from a queue.
func handleNotification(message []byte) {
	log.Printf("Received notification message: %s", message)
}

// NewRabbitMQ connects to a RabbitMQ server and returns a RabbitMQ.
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &RabbitMQ{conn, ch}, nil
}

// DeclareQueue creates a durable queue in RabbitMQ.
func (c *RabbitMQ) DeclareQueue(name string) error {
	_, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	return nil
}

// PublishMessage sends a message to a queue in RabbitMQ.
func (c *RabbitMQ) PublishMessage(queueName string, message []byte) error {
	err := c.ch.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}
	return nil
}

// ConsumeMessages(): starts consuming messages from a queue in RabbitMQ.
func (c *RabbitMQ) ConsumeMessages(queueName, consumerName string, messageHandler func([]byte), done <-chan struct{}) error {
	msgs, err := c.ch.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}
	log.Printf("registered consumer %v", msgs)

	for {
		select {
		case delivery, isOpen := <-msgs:
			if !isOpen {
				return nil
			}
			messageHandler(delivery.Body)
			delivery.Ack(true)
		case <-done:
			log.Println("Done consuming messages!")
			return nil
		}
	}
}

// Close(): closes the RabbitMQ and channel to RabbitMQ.
func (c *RabbitMQ) Close() error {
	err := c.ch.Close()
	if err != nil {
		return fmt.Errorf("failed to close the channel: %w", err)
	}
	err = c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close the RabbitMQ: %w", err)
	}
	return nil
}

func ManageQueue() {
	// Connect to RabbitMQ server
	rmq, err := NewRabbitMQ(config.Env.RabbitMqUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer rmq.Close()

	// Declare a queue for notifications
	notificationQueueName := "task"
	err = rmq.DeclareQueue(notificationQueueName)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Start consuming messages from the notification queue in a goroutine
	notificationChan := make(chan struct{})
	defer close(notificationChan)
	go rmq.ConsumeMessages("task", "notification-consumer", handleNotification, notificationChan)
}
