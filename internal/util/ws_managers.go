package util

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Mu   sync.Mutex // ป้องกัน Write พร้อมกัน
}

type Manager struct {
	Clients map[string]*Client
	Mu      sync.Mutex
}

var WSManager = Manager{
	Clients: make(map[string]*Client),
}

func (m *Manager) NotifyOldDevice(userID string) {
	m.Mu.Lock()

	if oldClient, ok := m.Clients[userID]; ok {
		// ✅ ส่งข้อความภาษาอังกฤษแจ้งเตือนเครื่องเก่า
		_ = oldClient.Conn.WriteJSON(gin.H{
			"type":    "FORCE_LOGOUT",
			"message": "Your account has been logged in from another device. This session will be terminated.",
		})

		// ให้เวลา Network ในการส่ง JSON ออกไปเล็กน้อยก่อนตัดการเชื่อมต่อ
		time.Sleep(100 * time.Millisecond)

		// ✅ ปิดการเชื่อมต่อ
		oldClient.Conn.Close()

		// ✅ ลบข้อมูลออกจากหน่วยความจำ
		delete(m.Clients, userID)
	}
	m.Mu.Unlock()
}
