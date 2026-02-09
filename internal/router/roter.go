package router

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler"
	socketHandler "github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler/socket"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/service"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/service/socket"
	socketService "github.com/AstralxOilx/Coding-Competition-Game/internal/service/socket"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 1. Repo
	userRepo := repository.NewUserRepository(database.DB)
	wsSvc := socket.NewWSService(userRepo)
	// 2. Services
	userService := service.NewUserService(userRepo, wsSvc) // ✅ สร้าง UserService
	authService := service.NewAuthService(userRepo)
	sessionService := socketService.NewWSService(userRepo)

	// 3. Handlers
	userHandler := handler.NewUserHandler(userService) // ✅ ส่ง Service ให้ UserHandler
	authHandler := handler.NewAuthHandler(authService)
	wsHandler := socketHandler.NewWSHandler(sessionService)

	v1 := r.Group("/api/v1")
	{
		InitAuthRoutes(v1, authHandler)
		InitWSRoutes(v1, wsHandler)
		InitUserRoutes(v1, userHandler)
	}

	return r
}
