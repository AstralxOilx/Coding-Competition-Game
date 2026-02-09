package router

import (
	socketHandler "github.com/AstralxOilx/Coding-Competition-Game/internal/api/handler/socket"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func InitWSRoutes(rg *gin.RouterGroup, h *socketHandler.WSHandler) {
	notifications := rg.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware())
	{
		// 🔌 สำหรับเชื่อมต่อ WebSocket (ใช้ ws://)
		notifications.GET("/broadcastfriend", h.HandleBroadcastFriendStatusWS)

		// 📊 สำหรับดึงรายชื่อคนออนไลน์เป็น JSON (ใช้ http://)
		// แนะนำให้เปลี่ยนชื่อ path ให้สื่อความหมายชัดเจนขึ้น
		notifications.GET("/online-list", h.HandleUserOnlineStats)
	}
}
