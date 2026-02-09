package socket

import (
	"log"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/model"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/util"
	"github.com/gorilla/websocket"
)

// 1. รวม Interface ไว้ที่เดียว (ใช้ชื่อ WSService ให้สื่อความหมายครอบคลุม)
type WSService interface {
	RegisterClient(userID string, conn *websocket.Conn)
	UnregisterClient(userID string, conn *websocket.Conn)
	IsUserOnline(userID string) bool
	// GetOnlineUserCount() int
	GetOnlineUserIDs() []string
	BroadcastFriendStatus(userID string, action string)
	GetOnlineFriendIDs(userID string) []string
}

type wsService struct {
	userRepo repository.UserRepo // สมมติว่านี่คือ Repository ที่จัดการเรื่อง User/Friends
}

func NewWSService(repo repository.UserRepo) WSService {
	return &wsService{
		userRepo: repo,
	}
}

// 2. ลงทะเบียนและจัดการ Force Logout
func (s *wsService) RegisterClient(userID string, conn *websocket.Conn) {
	util.WSManager.Mu.Lock()
	util.WSManager.Clients[userID] = &util.Client{Conn: conn}
	util.WSManager.Mu.Unlock()

	// 📣 ตะโกนบอกทุกคนทันที
	s.BroadcastFriendStatus(userID, "joined")
}

func (s *wsService) UnregisterClient(userID string, conn *websocket.Conn) {
	util.WSManager.Mu.Lock()
	delete(util.WSManager.Clients, userID)
	util.WSManager.Mu.Unlock()

	// 📣 ตะโกนบอกทุกคนว่าคนนี้ออกไปแล้ว
	s.BroadcastFriendStatus(userID, "left")
}

func (s *wsService) BroadcastFriendStatus(userID string, action string) {
	// 1. ดึงรายชื่อเพื่อนทั้งหมดจาก DB
	friendIDs, err := s.userRepo.FindFriendIDs(userID)
	if err != nil {
		log.Printf("Error fetching friends: %v", err)
		return
	}

	// ✅ 1.1 ดึงข้อมูล ชื่อ และ รูป ของคนที่เป็นต้นเหตุ (userID)
	userInfo, err := s.userRepo.FindUserInfo(userID)
	if err != nil {
		log.Printf("Error fetching user info for broadcast: %v", err)
		// ถ้าดึงข้อมูลไม่ได้ อาจจะใช้ค่า default หรือ return ออกไปเลย
		return
	}

	// 2. ทำ Map เพื่อ Search เพื่อนออนไลน์ได้เร็ว (O(1))
	isFriend := make(map[string]bool)
	for _, id := range friendIDs {
		isFriend[id] = true
	}

	util.WSManager.Mu.Lock()
	defer util.WSManager.Mu.Unlock()

	// 🚀 ปรับ Payload ให้มี DisplayName และ AvatarURL
	msg := map[string]interface{}{
		"type":         "FRIEND_STATUS_UPDATE",
		"friend_id":    userID,
		"display_name": userInfo.DisplayName, // ✅ เพิ่มชื่อ
		"avatar_url":   userInfo.AvatarURL,   // ✅ เพิ่มรูป
		"action":       action,
		"status":       model.UserStatus[model.StatusOnline],
	}

	if action == "left" {
		msg["status"] = model.UserStatus[model.StatusOffline]
	}

	// 3. วนลูปส่งให้เฉพาะเพื่อนที่ออนไลน์อยู่
	for id, client := range util.WSManager.Clients {
		if isFriend[id] {
			// แนะนำ: ใช้ Lock ที่ตัว client เองด้วยถ้ามี (ป้องกันการเขียนซ้อน)
			err := client.Conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Could not send message to friend %s: %v", id, err)
			}
		}
	}
}

// 4. เช็คสถานะออนไลน์ (สำหรับ UserService)
func (s *wsService) IsUserOnline(userID string) bool {
	util.WSManager.Mu.Lock()
	defer util.WSManager.Mu.Unlock()

	_, online := util.WSManager.Clients[userID]
	return online
}

func (s *wsService) GetOnlineUserIDs() []string {
	util.WSManager.Mu.Lock()
	defer util.WSManager.Mu.Unlock()

	ids := make([]string, 0, len(util.WSManager.Clients))
	for id := range util.WSManager.Clients {
		ids = append(ids, id)
	}
	return ids
}

func (s *wsService) GetOnlineFriendIDs(userID string) []string {
	// 1. ไปถาม Database ว่า "ฉันมีเพื่อนเป็นใครบ้าง" (Status = 1)
	allFriendIDs, err := s.userRepo.FindFriendIDs(userID)
	if err != nil {
		return []string{} // ถ้า Error หรือไม่มีเพื่อนเลย ให้ส่ง Array ว่าง
	}

	util.WSManager.Mu.Lock()
	defer util.WSManager.Mu.Unlock()

	var onlineFriends []string

	// 2. เช็คทีละคนว่า เพื่อนคนนั้น "ออนไลน์อยู่จริง" ในท่อ WebSocket หรือไม่
	for _, fID := range allFriendIDs {
		if _, online := util.WSManager.Clients[fID]; online {
			// ถ้าเพื่อนคนนี้มีชื่ออยู่ในท่อออนไลน์ ให้เก็บลงลิสต์
			onlineFriends = append(onlineFriends, fID)
		}
	}

	return onlineFriends
}
