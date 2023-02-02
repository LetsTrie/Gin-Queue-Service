package controllers

import (
	"fmt"
	"server/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTasks(c * gin.Context) {
    c.JSON(200, gin.H {
        "tasks": services.GetTasks(),
    })
}

func GetTask(c * gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    task := services.GetTask(id)
    if task == nil {
        c.JSON(404, gin.H {
            "message": "Task not found",
        })
        return
    }
    c.JSON(200, gin.H {
        "task": task,
    })
}

func CreateTask(c * gin.Context) {
    var task services.Task
    if err := c.ShouldBindJSON(&task);
    err != nil {
			  fmt.Println(err)
        c.JSON(400, gin.H {
            "message": "Bad request",
        })
        return
    }
    services.CreateTask(task)
    c.JSON(200, gin.H {
        "message": "Task created successfully",
    })
}