package router

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/service"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 1. Repositories
	userRepo := repository.NewUserRepository(database.DB)

	// 2. Services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)

	// 3. Handlers
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)

	v1 := r.Group("/api/v1")
	{
		InitAuthRoutes(v1, authHandler)
		InitUserRoutes(v1, userHandler)
	}

	return r
}
