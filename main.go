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
	// we are creating a done channel and deferring its closure.
	// We are then passing this channel to the ManageQueue() function,
	// which will use it to gracefully stop consuming messages when the channel is closed.

	forever := make(chan struct{})
	defer close(forever)

	err := service.ManageQueue("task", forever)

	if err != nil {
		log.Fatalf("Failed to manage queue: %s", err)
	}

	// Set up HTTP server
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Queue Service!",
		})
	})
	router.Run()
}
