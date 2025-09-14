package middleware

import (
	"net/http"
	"strings"

	"go-boilerplate/models/response"
	"go-boilerplate/services/interfaces"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService interfaces.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.BaseResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, response.BaseResponse{
				Success: false,
				Message: "Bearer token is required",
			})
			c.Abort()
			return
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.BaseResponse{
				Success: false,
				Message: "Invalid token",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		userID, err := authService.GetUserIDFromToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.BaseResponse{
				Success: false,
				Message: "Invalid token",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
