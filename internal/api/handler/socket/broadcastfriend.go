package handler

import (
	"log"
	"net/http"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/service/socket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WSHandler struct {
	wsService socket.WSService
}

func NewWSHandler(s socket.WSService) *WSHandler {
	return &WSHandler{wsService: s}
}

func (h *WSHandler) HandleBroadcastFriendStatusWS(c *gin.Context) {
	h.baseWSHandler(c)
}

func (h *WSHandler) HandleUserOnlineStats(c *gin.Context) {
	// แก้กลับมาใช้ GetOnlineUserIDs เพราะเราต้องการดึงค่ามาตอบ JSON
	userIDs := h.wsService.GetOnlineUserIDs()
	c.JSON(http.StatusOK, gin.H{
		"total_online": len(userIDs),
		"online_users": userIDs,
	})
}

// 🔌 สำหรับเชื่อมต่อ WebSocket (Real-time อยู่ที่นี่)
func (h *WSHandler) baseWSHandler(c *gin.Context) {

	// 1. ดึง userID ของเราเองจาก Context
	userID, _ := c.Get("user_id")
	userIDStr := userID.(string)

	// 2. อัปเกรดเป็น WebSocket และ Register ตามปกติ
	conn, _ := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	h.wsService.RegisterClient(userIDStr, conn)

	// ✅ 3. แก้ไขตรงนี้: เรียกใช้ฟังก์ชันที่ดึงเฉพาะ "เพื่อนที่ออนไลน์" เท่านั้น
	// ไม่ใช่ดึง GetOnlineUserIDs() (ที่เห็นทุกคน)
	onlineFriends := h.wsService.GetOnlineFriendIDs(userIDStr)

	// 4. ส่งกลับไปหา Client
	// ถ้าไม่มีเพื่อนออนไลน์เลย onlineFriends จะเป็น [] (Array ว่าง)
	conn.WriteJSON(gin.H{
		"type":           "INITIAL_FRIEND_LIST",
		"online_friends": onlineFriends,
	})

	defer func() {
		log.Printf("Closing connection for user: %s", userIDStr)
		h.wsService.UnregisterClient(userIDStr, conn)
		conn.Close()
	}()

	// 🔄 Loop ค้างสายไว้เพื่อให้สถานะออนไลน์ยังคงอยู่
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("Connection closed by client or error: %v", err)
			break
		}
	}
}
