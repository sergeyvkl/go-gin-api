package controllers

import (
	"go-gin-api/dto"
	"go-gin-api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// GetAllUsers возвращает всех пользователей с пагинацией
func (uc *UserController) GetAllUsers(c *gin.Context) {
	// Получаем параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	users, total, err := uc.service.GetAllUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to fetch users: " + err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Users retrieved successfully",
		Data: gin.H{
			"users":     users,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserByID возвращает пользователя по ID
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	
	user, err := uc.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Database error",
			Message: "Failed to fetch user: " + err.Error(),
			Status:  http.StatusInternalServerError,
		})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "User not found",
			Message: "User with the specified ID does not exist",
			Status:  http.StatusNotFound,
		})
		return
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser создает нового пользователя
func (uc *UserController) CreateUser(c *gin.Context) {
	var createReq dto.CreateUserRequest
	
	if err := c.ShouldBindJSON(&createReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation error",
			Message: "Invalid request data: " + err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}
	
	user, err := uc.service.CreateUser(&createReq)
	if err != nil {
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Error:   "Creation failed",
			Message: err.Error(),
			Status:  http.StatusConflict,
		})
		return
	}
	
	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser обновляет данные пользователя
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	
	var updateReq dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation error",
			Message: "Invalid request data: " + err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}
	
	user, err := uc.service.UpdateUser(id, &updateReq)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		} else if err.Error() == "email already in use" {
			status = http.StatusConflict
		}
		
		c.JSON(status, dto.ErrorResponse{
			Error:   "Update failed",
			Message: err.Error(),
			Status:  status,
		})
		return
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser удаляет пользователя
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	
	if err := uc.service.DeleteUser(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		
		c.JSON(status, dto.ErrorResponse{
			Error:   "Deletion failed",
			Message: err.Error(),
			Status:  status,
		})
		return
	}
	
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "User deleted successfully",
	})
}

// HealthCheck проверяет состояние сервиса
func (uc *UserController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"service": "go-gin-api",
	})
}