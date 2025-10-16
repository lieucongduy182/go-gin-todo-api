package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lieucongduy182/go-gin-todo-api/database"
	"github.com/lieucongduy182/go-gin-todo-api/models"
)

func GetTasks(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var tasks []models.Task
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"count": len(tasks),
	})
}

func GetTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	taskID := c.Param("id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": task})
}

func CreateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var input models.CreateTaskInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := models.Task{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Completed:   false,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": task, "message": "Created Task successfully"})
}

func UpdateTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	taskID := c.Param("id")

	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Default().Printf("Error %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Update fields
	if input.Title != "" {
		task.Title = input.Title
	}

	if input.Description != "" {
		task.Description = input.Description
	}

	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	if err := database.DB.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    task,
		"message": "Updated Task Successfully",
	})
}

func DeleteTask(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	taskID := c.Param("id")

	var task models.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted Task Successfully"})
}
