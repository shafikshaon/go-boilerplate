package handlers

import (
	"net/http"
	"strconv"

	"go-boilerplate/logger"
	"go-boilerplate/models/request"
	"go-boilerplate/models/response"
	"go-boilerplate/services/interfaces"
	"go-boilerplate/utilities"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService interfaces.UserService
}

func NewUserHandler(userService interfaces.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "CreateUser request received", nil)
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, "CreateUser: invalid request body", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := utilities.ValidateStruct(&req); err != nil {
		logger.Warn(ctx, "CreateUser: validation failed", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		logger.Error(ctx, "CreateUser failed", map[string]any{"email": req.Email, "error": err.Error()})
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	logger.Info(ctx, "User created successfully", map[string]any{"email": req.Email})
	c.JSON(http.StatusCreated, response.BaseResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")
	logger.Info(ctx, "GetUser request received", map[string]any{"id": idParam})
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logger.Warn(ctx, "GetUser: invalid user ID", map[string]any{"id": idParam, "error": err.Error()})
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		logger.Warn(ctx, "GetUser: not found", map[string]any{"id": id, "error": err.Error()})
		c.JSON(http.StatusNotFound, response.BaseResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	logger.Info(ctx, "GetUser: success", map[string]any{"id": id})
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	users, err := h.userService.GetUsers(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Failed to retrieve users",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
