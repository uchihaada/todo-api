package main

import (
	"ada/handlers"
	"ada/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	fileStorage := storage.NewFileStorage("task.json")
	defer fileStorage.Close()

	router := gin.Default()

	taskHandler := handlers.NewTaskHandler(fileStorage)

	router.GET("/tasks", taskHandler.GetTask)
	router.POST("/tasks", taskHandler.CreateTask)
	router.PUT("/tasks/:id", taskHandler.UpdateTask)
	router.DELETE("/tasks/:id", taskHandler.DeleteTask)

	router.Run(":8000")
}
