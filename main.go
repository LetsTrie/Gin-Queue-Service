package main

import (
	"server/controllers"
	"server/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	// Initialize Gin router
	r := gin.Default()
	
	// Define route handlers
	r.GET("/tasks", controllers.GetTasks)
	r.GET("/tasks/:id", controllers.GetTask)
	r.POST("/tasks", controllers.CreateTask)
	
	// Start server
	r.Run()
}