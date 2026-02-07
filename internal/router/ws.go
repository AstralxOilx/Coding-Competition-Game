package router

import (
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler"
	wsHandler "github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler/ws"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitWSRoutes(rg *gin.RouterGroup, h *handler.UserHandler) {
	// สร้างกลุ่ม /notifications ต่อท้าย /api/v1
	// URL จะกลายเป็น /api/v1/notifications
	notifications := rg.Group("/notifications")

	// ใส่ Middleware เพื่อให้ปลอดภัย (ดึง user_id จาก token)
	notifications.Use(middleware.AuthMiddleware())

	{
		notifications.GET("/session", wsHandler.HandleSessionWS)
	}
}
