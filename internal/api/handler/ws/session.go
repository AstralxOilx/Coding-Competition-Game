package handler

import (
	"log"
	"net/http"
	"sync"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// ตั้งค่า Upgrader เพื่ออัปเกรด HTTP เป็น WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// ตรวจสอบ Origin (ในโปรดักชั่นควรเช็คโดเมนเพื่อป้องกัน CSRF)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleSessionWS(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userIDStr := userID.(string)

	// 2. อัปเกรดการเชื่อมต่อ
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return // Upgrade failure ไม่ส่ง 500 แต่จะตัดการเชื่อมต่อ
	}

	// 3. ป้องกัน Panic จาก WSManager (ตรวจสอบว่า Map ถูกสร้างหรือยัง)
	util.WSManager.Mu.Lock()
	if util.WSManager.Clients == nil {
		util.WSManager.Clients = make(map[string]*util.Client)
	}

	// 4. เตะคนเก่า (NotifyOldDevice)
	// ย้ายมาทำในนี้เพื่อความชัวร์ หรือเรียก util.WSManager.NotifyOldDevice(userIDStr)
	if oldClient, ok := util.WSManager.Clients[userIDStr]; ok {
		_ = oldClient.Conn.WriteJSON(gin.H{"type": "FORCE_LOGOUT"})
		oldClient.Conn.Close()
	}

	// 5. บันทึกเครื่องใหม่
	util.WSManager.Clients[userIDStr] = &util.Client{
		Conn: conn,
		Mu:   sync.Mutex{},
	}
	util.WSManager.Mu.Unlock()

	// 6. รักษาการเชื่อมต่อ
	defer func() {
		util.WSManager.Mu.Lock()
		delete(util.WSManager.Clients, userIDStr)
		util.WSManager.Mu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
