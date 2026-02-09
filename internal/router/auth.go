package router

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(rg *gin.RouterGroup, h *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		// 🔓 Public Routes (ไม่ต้องใส่ Token ก็เรียกได้)
		auth.POST("/signup", h.Signup)
		auth.POST("/signin", h.Signin)

		// 🔄 Refresh Token (ปกติส่ง Refresh Token มาใน Body จึงไม่ต้องใช้ AuthMiddleware)
		auth.POST("/refresh", h.Refresh)

		// 🔒 Protected Routes (ต้องมี Access Token ที่ยังไม่หมดอายุ)
		protected := auth.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// ตัวอย่าง: logout หรือ profile
			// protected.POST("/logout", h.Logout)
		}
	}
}
