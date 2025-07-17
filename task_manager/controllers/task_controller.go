package controllers

import (
	"net/http"

	"task_manager/data"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(c *gin.Context) {
	tasks, err := data.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := data.CreateTask(newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var updatedTask models.Task

	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	task, err := data.UpdateTask(id, updatedTask)

	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"message": "Task not Found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := data.DeleteTask(id)
	if err != nil {
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
