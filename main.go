package main

import (
	"log"
	"net/http"
	config "server/config"
	database "server/database"
	service "server/services"
	fcm "server/services/fcm"

	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	database.RedisConfig()
	fcm.FcmConfig()
}

func main() {

	// Connect to RabbitMQ server
	rmq, err := service.NewRabbitMQ(config.Env.RabbitMqUrl)
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
	// Start consuming messages from the notification queue in a goroutine
	notificationChan := make(chan struct{})
	defer close(notificationChan)
	go rmq.ConsumeMessages("task", "notification-consumer", handleNotification, notificationChan)

	// Set up HTTP server
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Queue Service!",
		})
	})
	router.Run()
}

// handleNotification is a function that handles a notification message received from a queue.
func handleNotification(message []byte) {
	log.Printf("Received notification message: %s", message)
}
