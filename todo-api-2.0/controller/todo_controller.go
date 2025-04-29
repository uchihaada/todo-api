package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"todo-api-2.0/config"
	"todo-api-2.0/models"
)

func CreateTask(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.DB.Create(&todo).Error; err != nil {
		fmt.Println("Error creating task:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task created successfully", "task": todo})
}

func GetTasks(c *gin.Context) {
	var todos []models.Todo
	if err := config.DB.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	if len(todos) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No tasks found", "tasks": []models.Todo{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": todos})
}

func GetTasksById(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	if err := config.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"task": todo})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var existingTodo models.Todo
	if result := config.DB.Where("id = ?", id).First(&existingTodo); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
			return
		}
	}

	if err := config.DB.Model(&existingTodo).Update("title", todo.Title).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task title updated successfully", "task": existingTodo})
}

func MarksTaskAsCompleted(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo

	// Fetch the task by ID
	if err := config.DB.Where("id = ?", id).First(&todo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	// Check if the task is already completed
	if todo.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task already completed"})
		return
	}

	// Mark the task as completed
	if err := config.DB.Model(&todo).Update("completed", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Return the updated task
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as completed successfully", "task": todo})
}

func DeleteTask(c *gin.Context) {
	// Fetch all completed tasks
	var tasks []models.Todo
	if err := config.DB.Where("completed = ?", true).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	// Check if there are no completed tasks
	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No completed tasks found"})
		return
	}

	// Delete all completed tasks
	if err := config.DB.Where("completed = ?", true).Delete(&models.Todo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete completed tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All completed tasks deleted successfully"})
}
