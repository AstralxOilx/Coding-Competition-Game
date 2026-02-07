package router

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(rg *gin.RouterGroup, h *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/signup", h.Signup)
		auth.POST("/signin", h.Signin)
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/refresh", h.Refresh)
		}

	}
}
