package routes

import (
	"go-boilerplate/handlers"
	"go-boilerplate/middleware"
	"go-boilerplate/services/interfaces"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	healthHandler *handlers.HealthHandler,
	authService interfaces.AuthService,
) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", healthHandler.Check)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)

			// Protected auth routes
			authProtected := auth.Use(middleware.AuthMiddleware(authService))
			{
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.GET("/me", authHandler.Me)
			}
		}

		// User routes
		users := v1.Group("/users")
		{
			users.POST("/", userHandler.CreateUser)
			users.GET("/", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUser)

			// Protected routes
			protected := users.Use(middleware.AuthMiddleware(authService))
			{
				protected.PUT("/:id", userHandler.UpdateUser)
				protected.DELETE("/:id", userHandler.DeleteUser)
			}
		}
	}

	return router
}
