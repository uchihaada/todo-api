package handlers

import (
	"ada/models"
	"ada/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Storage *storage.FileStorage
}

func NewTaskHandler(storage *storage.FileStorage) *TaskHandler {
	return &TaskHandler{Storage: storage}
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	tasks, err := h.Storage.LoadTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tasks, _ := h.Storage.LoadTasks()
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	h.Storage.SaveTasks(tasks)
	c.JSON(http.StatusCreated, newTask)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tasks, _ := h.Storage.LoadTasks()
	for i, task := range tasks {
		if task.ID == id {
			updatedTask.ID = id
			tasks[i] = updatedTask
			h.Storage.SaveTasks(tasks)
			c.JSON(http.StatusOK, updatedTask)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	tasks, _ := h.Storage.LoadTasks()
	newTasks := []models.Task{}

	for _, task := range tasks {
		if task.ID != id {
			newTasks = append(newTasks, task)
		}
	}

	h.Storage.SaveTasks(newTasks)
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
