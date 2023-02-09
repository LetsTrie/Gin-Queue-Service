package main

import (
	"net/http"
	"server/config"
	db "server/database"
	service "server/services"
	fcm "server/services/fcm"

	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	db.RedisConfig()
	fcm.FcmConfig()
}

var userId string = "63ba96fe96202297db352f7c"

func main() {
	router := gin.Default()

	emailQueue := service.NewRedisQueue("emailQueue")
	notificationQueue := service.NewRedisQueue("notificationQueue")

	go func() {
		emailQueue.Set("item1")
		notificationQueue.Set("item2")
	}()

	router.GET("/", func(c *gin.Context) {
		message := &fcm.MulticastMessage{
			Notification: &fcm.Notification{
				Title: "🔥 Rock On 🎸",
				Body:  "Get ready to rock with our latest update 🤘",
			},
		}

		fcm.SendPushNotificationToUser(userId, message)
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Gin Queue Service!",
		})
	})

	router.Run()
}

// ~HOME/go/bin/CompileDaemon -command="go run main.go"
