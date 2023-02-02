package main

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.GET("/tasks", controllers.GetTasks)
	r.GET("/tasks/:id", controllers.GetTask)
	r.POST("/tasks", controllers.CreateTask)
	
	r.Run(":8080")
}