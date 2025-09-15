//go:build wireinject

package cmd

import (
	"github.com/google/wire"

	"go-boilerplate/config"
	"go-boilerplate/database"
	"go-boilerplate/handlers"
	"go-boilerplate/logger"
	"go-boilerplate/repository"
	repoInterfaces "go-boilerplate/repository/interfaces"
	"go-boilerplate/routes"
	"go-boilerplate/services"
	serviceInterfaces "go-boilerplate/services/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Provider functions used by Wire
func provideConfig() *config.Config                          { return config.Load() }
func provideDB(cfg *config.Config) (*gorm.DB, error)         { return database.NewConnection(cfg) }
func provideRedis(cfg *config.Config) (*redis.Client, error) { return database.NewRedisConnection(cfg) }

// Repositories
func provideUserRepository(db *gorm.DB) repoInterfaces.UserRepository {
	return repository.NewUserRepository(db)
}
func provideRedisRepository(rdb *redis.Client) repoInterfaces.RedisRepository {
	return repository.NewRedisRepository(rdb)
}

// Services
func provideAuthService(cfg *config.Config) serviceInterfaces.AuthService {
	return services.NewAuthService(cfg)
}
func provideRedisService(redisRepo repoInterfaces.RedisRepository) serviceInterfaces.RedisService {
	return services.NewRedisService(redisRepo)
}
func provideUserService(userRepo repoInterfaces.UserRepository, auth serviceInterfaces.AuthService, redis serviceInterfaces.RedisService) serviceInterfaces.UserService {
	return services.NewUserService(userRepo, auth, redis)
}

// Handlers
func provideUserHandler(svc serviceInterfaces.UserService) *handlers.UserHandler {
	return handlers.NewUserHandler(svc)
}
func provideAuthHandler(svc serviceInterfaces.UserService) *handlers.AuthHandler {
	return handlers.NewAuthHandler(svc)
}
func provideHealthHandler() *handlers.HealthHandler { return handlers.NewHealthHandler() }

// Router
func provideRouter(uh *handlers.UserHandler, ah *handlers.AuthHandler, hh *handlers.HealthHandler, auth serviceInterfaces.AuthService) *gin.Engine {
	r := routes.SetupRoutes(uh, ah, hh, auth)
	r.Use(logger.GinMiddleware())
	return r
}

// InitializeApp is the Wire injector. The actual implementation is generated into wire_gen.go.
func InitializeApp() (*gin.Engine, error) {
	wire.Build(
		provideConfig,
		provideDB,
		provideRedis,
		provideUserRepository,
		provideRedisRepository,
		provideAuthService,
		provideRedisService,
		provideUserService,
		provideUserHandler,
		provideAuthHandler,
		provideHealthHandler,
		provideRouter,
	)
	return nil, nil
}
