package model

import (
	"time"

	"gorm.io/gorm"
)

type LoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:14;index" json:"user_id"` // เชื่อมกับ Users.ID
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	LoginAt   time.Time `json:"login_at"`
	Status    string    `json:"status"` // "SUCCESS" หรือ "FAILED"
}

type MatchLogs struct {
	ID       string `gorm:"primaryKey;size:20" json:"id"` // แนะนำใช้ ID ยาวขึ้น เช่น Match ID
	UserID   string `gorm:"size:14;index;not null" json:"user_id"`
	ModeName string `gorm:"type:varchar(20);index;not null" json:"mode_name"` // ใช้ค่าจาก const เช่น ModeClassic

	// --- ผลการแข่งขัน ---
	Result    string `gorm:"type:varchar(10);not null" json:"result"` // 'WIN', 'LOSS', 'DRAW'
	PointEarn int    `gorm:"default:0" json:"point_earn"`             // แต้มที่ได้รับหรือเสียไป (+/-)
	ExpEarn   int    `gorm:"default:0" json:"exp_earn"`               // Player EXP ที่ได้รับ

	// --- รายละเอียดการเล่น ---
	Duration   int    `json:"duration"`                   // เวลาที่ใช้ในการเล่น (วินาที)
	Score      int    `json:"score"`                      // คะแนนดิบที่ทำได้ในเกมนั้น
	OpponentID string `gorm:"size:14" json:"opponent_id"` // ID คู่แข่ง (ถ้ามี เช่นโหมด Duel)

	// --- Metadata ---
	PlayedAt time.Time `gorm:"autoCreateTime" json:"played_at"`
	gorm.Model
}
