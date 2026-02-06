package model

import "time"

type LoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:14;index" json:"user_id"` // เชื่อมกับ Users.ID
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	LoginAt   time.Time `json:"login_at"`
	Status    string    `json:"status"` // "SUCCESS" หรือ "FAILED"
}
