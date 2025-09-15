package cmd

import (
	"context"

	"go-boilerplate/config"
	"go-boilerplate/database"
	"go-boilerplate/handlers"
	"go-boilerplate/logger"
	"go-boilerplate/repository"
	"go-boilerplate/routes"
	"go-boilerplate/services"

	"github.com/gin-gonic/gin"
)

// initializeAppManual wires dependencies manually. It is used by the wire-generated code.
func initializeAppManual() (*gin.Engine, error) {
	ctx := context.Background()
	// Load configuration
	cfg := config.Load()

	// Database connections
	db, err := database.NewConnection(cfg)
	if err != nil {
		return nil, err
	}

	rdb, err := database.NewRedisConnection(cfg)
	if err != nil {
		return nil, err
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)

	// Services
	authService := services.NewAuthService(cfg)
	redisService := services.NewRedisService(redisRepo)
	userService := services.NewUserService(userRepo, authService, redisService)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService)
	healthHandler := handlers.NewHealthHandler()

	// Setup routes
	router := routes.SetupRoutes(userHandler, authHandler, healthHandler, authService)
	// Attach tracing middleware
	router.Use(logger.GinMiddleware())

	logger.Info(ctx, "Application initialized successfully", nil)
	return router, nil
}
