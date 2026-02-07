package router

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Setup Repositories
	userRepo := repository.NewUserRepository(database.DB)

	// Setup Handlers
	authHandler := handler.NewAuthHandler(userRepo)
	userHandler := handler.NewUserHandler(userRepo)

	v1 := r.Group("/api/v1")
	{
		// แยกกลุ่มการตั้งค่า Route ออกไปเป็นฟังก์ชันย่อย
		InitAuthRoutes(v1, authHandler)
		InitUserRoutes(v1, userHandler)
		InitWSRoutes(v1, userHandler)
	}

	return r
}
