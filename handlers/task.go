package handlers

import (
	"ada/models"
	"ada/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Taskhandler struct {
	Storage *storage.FileStorage
}

func NewTaskHandler(storage *storage.FileStorage) *Taskhandler {
	return &Taskhandler{Storage: storage}
}

func (h *Taskhandler) getTask(c *gin.Context) {
	tasks, err := h.Storage.LoadTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load tasks"})
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Taskhandler) CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request Body"})
		return
	}

	tasks, _ := h.Storage.LoadTasks()
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)
	h.Storage.SaveTasks(tasks)
	c.JSON(http.StatusCreated, newTask)
}

func (h *Taskhandler) updateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task id"})
		return
	}
	var updatdTask models.Task
	if err := c.ShouldBindJSON(&updatdTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	tasks, _ := h.Storage.LoadTasks()
	for i, task := range tasks {
		if task.ID == id {
			tasks[i] = updatdTask
			h.Storage.SaveTasks(tasks)
			c.JSON(http.StatusOK, updatdTask)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func (h *Taskhandler) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	tasks, _ := h.Storage.LoadTasks()
	newTask := []models.Task{}
	for _, task := range tasks {
		if task.ID != id {
			newTask = append(newTask, task)
		}
	}
	h.Storage.SaveTasks(newTask)
	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
