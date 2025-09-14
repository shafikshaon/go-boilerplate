package handlers

import (
	"net/http"

	"go-boilerplate/logger"
	"go-boilerplate/models/request"
	"go-boilerplate/models/response"
	"go-boilerplate/services/interfaces"
	"go-boilerplate/utilities"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService interfaces.UserService
}

func NewAuthHandler(userService interfaces.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Register request received", nil)
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, "Register: invalid request body", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := utilities.ValidateStruct(&req); err != nil {
		logger.Warn(ctx, "Register: validation failed", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		logger.Error(ctx, "Register failed", map[string]any{"email": req.Email, "error": err.Error()})
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Failed to register user",
			Error:   err.Error(),
		})
		return
	}

	logger.Info(ctx, "Register successful", map[string]any{"email": req.Email})
	c.JSON(http.StatusCreated, response.BaseResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Login request received", nil)
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn(ctx, "Login: invalid request body", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := utilities.ValidateStruct(&req); err != nil {
		logger.Warn(ctx, "Login: validation failed", map[string]any{"error": err.Error()})
		c.JSON(http.StatusBadRequest, response.BaseResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	loginResponse, err := h.userService.Login(&req)
	if err != nil {
		logger.Warn(ctx, "Login failed", map[string]any{"email": req.Email, "error": err.Error()})
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
		return
	}

	logger.Info(ctx, "Login successful", map[string]any{"email": req.Email})
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Login successful",
		Data:    loginResponse,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Logout request received", nil)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		logger.Warn(ctx, "Logout: unauthenticated", nil)
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		logger.Error(ctx, "Logout: invalid user id type", nil)
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	if err := h.userService.Logout(userID); err != nil {
		logger.Error(ctx, "Logout failed", map[string]any{"user_id": userID, "error": err.Error()})
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Failed to logout",
			Error:   err.Error(),
		})
		return
	}

	logger.Info(ctx, "Logout successful", map[string]any{"user_id": userID})
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Logout successful",
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Me request received", nil)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		logger.Warn(ctx, "Me: unauthenticated", nil)
		c.JSON(http.StatusUnauthorized, response.BaseResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		logger.Error(ctx, "Me: invalid user id type", nil)
		c.JSON(http.StatusInternalServerError, response.BaseResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		logger.Warn(ctx, "Me: user not found", map[string]any{"user_id": userID, "error": err.Error()})
		c.JSON(http.StatusNotFound, response.BaseResponse{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	logger.Info(ctx, "Me: success", map[string]any{"user_id": userID})
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		Data:    user,
	})
}
