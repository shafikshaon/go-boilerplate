package main

import (
	"context"

	"go-boilerplate/cmd"
	"go-boilerplate/config"
	"go-boilerplate/logger"
)

func main() {
	// Initialize logger from env/Viper
	_ = logger.Init()
	ctx := context.Background()

	// Load configuration using Viper (.env + environment variables)
	cfg := config.Load()
	port := cfg.Port

	// Initialize application
	app, err := cmd.InitializeApp()
	logger.Must(ctx, err, "Failed to initialize app", nil)

	logger.Info(ctx, "Server starting", map[string]any{"port": port})
	if err := app.Run(":" + port); err != nil {
		logger.Must(ctx, err, "Server failed to start", nil)
	}
}
