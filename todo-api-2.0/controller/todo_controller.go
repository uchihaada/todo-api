package controller

import (
	"fmt"
	"net/http"
	"sync"

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
	if err := config.DB.Create(&todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task created successfully", "task": todo})
}

func GetTasks(c *gin.Context) {
	var todos []models.Todo
	if err := config.DB.Find(&todos); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": todos})
}

func GetTasksById(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := config.DB.Where("id = ?", id).First(&todo); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
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

	if result := config.DB.Where("id = ?", id).First(&todo); result.Error != nil {

		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
			return
		}
	}

	if todo.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task already completed"})
		return
	}

	if err := config.DB.Model(&todo).Updates(models.Todo{Title: todo.Title, Completed: todo.Completed}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": todo})
}

func MarksTaskAsCompleted(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if result := config.DB.Where("id = ?", id).Find(&todo); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
			return
		}
	}
	if todo.Completed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task already completed"})
		return
	}
	if err := config.DB.Model(&todo).Updates(models.Todo{Title: todo.Title, Completed: true}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task marked as completed successfully", "task": todo})
}

func DeleteTask(c *gin.Context) {

	var tasks []models.Todo

	if err := config.DB.Where("Completed = ?", true).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No completed tasks found"})
		return
	}

	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(t models.Todo) {
			defer wg.Done()
			if err := config.DB.Delete(&t).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete task with id %d", t.ID)})
				return
			}
		}(task)
	}
	wg.Wait()

	c.JSON(http.StatusOK, gin.H{"message": "All completed tasks deleted successfully"})
}
