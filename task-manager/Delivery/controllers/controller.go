package controllers

import (
	"net/http"
	"task-manager/Domain"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	taskUsecase domain.TaskUsecase
	userUsecase domain.UserUsecase
}

func NewController(taskUsecase domain.TaskUsecase, userUsecase domain.UserUsecase) *Controller {
	return &Controller{
		taskUsecase: taskUsecase,
		userUsecase: userUsecase,
	}
}

// Task Controllers
func (ctrl *Controller) GetTasks(c *gin.Context) {
	tasks, err := ctrl.taskUsecase.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (ctrl *Controller) GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := ctrl.taskUsecase.GetTaskByID(id)
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

func (ctrl *Controller) CreateTask(c *gin.Context) {
	var newTask domain.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := ctrl.taskUsecase.CreateTask(newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (ctrl *Controller) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var updatedTask domain.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := ctrl.taskUsecase.UpdateTask(id, updatedTask)
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

func (ctrl *Controller) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := ctrl.taskUsecase.DeleteTask(id)
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

// User Controllers
func (ctrl *Controller) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdUser, err := ctrl.userUsecase.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdUser)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginResp, err := ctrl.userUsecase.LoginUser(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loginResp)
}

func (ctrl *Controller) Promote(c *gin.Context) {
	id := c.Param("id")
	updatedUser, err := ctrl.userUsecase.PromoteUser(id)
	if err != nil {
		if err.Error() == "invalid user ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (ctrl *Controller) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := ctrl.userUsecase.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, user)
}