package main

import (
	"server/config"
	db "server/database"
	service "server/services"

	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	db.ConnectRedis()
}

func main() {
	r := gin.Default()
	emailQueue := service.NewRedisQueue("emailQueue")
	notificationQueue := service.NewRedisQueue("notificationQueue")

	go func() {
		emailQueue.Set("item1")
		notificationQueue.Set("item2")
	}()

	r.Run()
}

// ~HOME/go/bin/CompileDaemon -command="go run main.go"
