package service

import (
	"fmt"
	"log"
	"server/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ represents a connection to a RabbitMQ server.
type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewRabbitMQ connects to a RabbitMQ server and returns a RabbitMQ.
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	// Dial establishes a connection to the RabbitMQ server.
	conn, err := amqp.Dial(url)
	if err != nil {
		// If there was an error connecting to the server, return an error.
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a new channel for the connection.
	ch, err := conn.Channel()
	if err != nil {
		// If there was an error creating the channel, return an error.
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Return a pointer to a new RabbitMQ object containing the connection and channel.
	return &RabbitMQ{conn, ch}, nil
}

// DeclareQueue creates a durable queue in RabbitMQ.
func (c *RabbitMQ) DeclareQueue(name string) error {
	// Call QueueDeclare on the RabbitMQ channel to create a durable queue with the given name.
	// The queue is durable, which means it will survive a RabbitMQ server restart.
	// The queue is not auto-deleted when it is no longer in use.
	// The queue is not exclusive, which means other connections can access it.
	// There is no wait time for the queue to be declared.
	// Additional arguments can be passed in as a map[string]interface{}.
	_, err := c.ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		// If there was an error declaring the queue, return an error.
		return fmt.Errorf("failed to declare a queue: %w", err)
	}
	// Return nil if the queue was successfully declared.
	return nil
}

// ConsumeMessages starts consuming messages from a queue in RabbitMQ.
// It takes a queue name, a consumer name, a message handler function, and a done channel.
// The message handler function takes a []byte and processes the message.
// The done channel is used to signal that message consumption should stop.
// This function blocks until it receives a message from RabbitMQ or a signal on the done channel.
func (c *RabbitMQ) ConsumeMessages(queueName, consumerName string, messageHandler func([]byte), done <-chan struct{}) error {
	// Call Consume on the RabbitMQ channel to register a consumer for the given queue and consumer names.
	// The consumer is not auto-acked, which means that the message must be explicitly acknowledged with Ack.
	// There is no timeout for the consumer.
	// Additional arguments can be passed in as a map[string]interface{}.
	msgs, err := c.ch.Consume(queueName, consumerName, false, false, false, false, nil)
	if err != nil {
		// If there was an error registering the consumer, return an error.
		return fmt.Errorf("failed to register a consumer: %w", err)
	}
	// Log the registration of the consumer.
	log.Printf("registered consumer %v", msgs)

	// Start consuming messages until a message or signal is received.
	for {
		select {
		case delivery, isOpen := <-msgs:
			if !isOpen {
				// If the message channel has been closed, stop consuming messages and return nil.
				return nil
			}
			// Pass the message body to the message handler.
			messageHandler(delivery.Body)
			// Acknowledge the message to RabbitMQ.
			delivery.Ack(true)
		case <-done:
			// If a signal is received on the done channel, stop consuming messages and return nil.
			log.Println("Done consuming messages!")
			return nil
		}
	}
}

// Close closes the connection and channel to RabbitMQ.
// It returns an error if there was a problem closing either.
func (c *RabbitMQ) Close() error {
	// Close the RabbitMQ channel.
	err := c.ch.Close()
	if err != nil {
		return fmt.Errorf("failed to close the channel: %w", err)
	}
	// Close the RabbitMQ connection.
	err = c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close the RabbitMQ: %w", err)
	}
	// Return nil if the RabbitMQ connection and channel were closed successfully.
	return nil
}

// ManageQueue connects to RabbitMQ server and consumes messages from a specified queue using the handleNotification function.
// It returns an error if there was a problem with the RabbitMQ connection, failed to declare the queue, or failed to consume messages.
func ManageQueue(queueName string, forever <-chan struct{}) error {
	// Connect to RabbitMQ server
	rmq, err := NewRabbitMQ(config.Env.RabbitMqUrl)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer rmq.Close()

	// Declare a queue for notifications
	err = rmq.DeclareQueue(queueName)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Start consuming messages from the queue
	err = rmq.ConsumeMessages(queueName, "notification-consumer", handleNotification, forever)
	if err != nil {
		log.Fatalf("Failed to consume messages: %s", err)
	}

	// Return nil if the function executes successfully
	return nil
}

// handleNotification is a function that handles a notification message received from a queue.
func handleNotification(message []byte) {
	log.Printf("Received notification message: %s", message)
}
