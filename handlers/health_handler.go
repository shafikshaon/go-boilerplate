package handlers

import (
	"net/http"
	"time"

	"go-boilerplate/logger"
	"go-boilerplate/models/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	ctx := c.Request.Context()
	logger.Info(ctx, "Health check", nil)
	c.JSON(http.StatusOK, response.BaseResponse{
		Success: true,
		Message: "Service is healthy",
		Data: gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
			"service":   "go-boilerplate",
		},
	})
}
