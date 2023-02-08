package main

import (
	"server/config"
	db "server/database"

	"github.com/gin-gonic/gin"
)

func init() {
	config.Init()
	db.ConnectRedis()
}

// ~HOME/go/bin/CompileDaemon -command="go run main.go"
func main() {
	r := gin.Default()

	r.Run()
}