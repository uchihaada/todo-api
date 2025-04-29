package routes

import (
	"github.com/gin-gonic/gin"
	"todo-api-2.0/controller"
)

func SetUpRoutes(r *gin.Engine) {
	r.GET("/api/todo", controller.GetTasks)
	r.GET("/api/todo/:id", controller.GetTasksById)
	r.POST("/api/todo", controller.CreateTask)
	r.PUT("/api/todo/:id", controller.UpdateTask)
	r.DELETE("/api/todo/:id", controller.DeleteTask)
	r.PATCH("/api/todo/:id/markascomplete", controller.MarksTaskAsCompleted) // Use PATCH here
}
