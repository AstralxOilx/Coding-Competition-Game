package router

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitUserRoutes(rg *gin.RouterGroup, h *handler.UserHandler) {
	users := rg.Group("/users")

	// ใช้ Middleware กับทั้งกลุ่มนี้เลย
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("/", middleware.RoleMiddleware(1), h.FindAllUser)

	}
}
