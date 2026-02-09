package model

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID          string `gorm:"primaryKey;size:14" json:"id"`
	DisplayName string `gorm:"type:varchar(14);not null" json:"display_name"`
	UserName    string `gorm:"type:varchar(20);not null;unique" json:"user_name"`
	Password    string `gorm:"not null" json:"-"`
	Email       string `gorm:"not null;unique" json:"email"`

	PlayerLevel int `gorm:"default:1" json:"player_level"`
	PlayerExp   int `gorm:"default:0" json:"player_exp"`

	Ranks []UserRanks `gorm:"foreignKey:UserID" json:"ranks"`

	// เก็บเป็นตัวเลข (0=DEV, 1=GM, 2=USER, 3=GUEST)
	UserRole int `gorm:"default:2" json:"user_role"`
	// เก็บเป็นตัวเลข (0=ONLINE, 1=OFFLINE, 2=UNAVAILABLE)
	Status int `gorm:"default:1" json:"status"`

	IsActive  bool       `gorm:"default:true" json:"is_active"`
	AvatarURL string     `json:"avatar_url"`
	LastLogin *time.Time `json:"last_login"`
	gorm.Model
}

type UserRanks struct {
	gorm.Model
	UserID   string `gorm:"size:14;index" json:"user_id"`
	ModeName string `gorm:"type:varchar(20);index" json:"mode_name"`

	// --- ระบบความเก่ง (เปลี่ยนเป็น int ทั้งหมด) ---
	Rank      int `gorm:"default:0" json:"rank"`      // เก็บ index 0-13
	RankTier  int `gorm:"default:3" json:"rank_tier"` // เก็บเลข 3, 2, 1
	RankPoint int `gorm:"default:0" json:"rank_point"`

	TotalGames int `gorm:"default:0" json:"total_games"`
	Win        int `gorm:"default:0" json:"win"`
	Loss       int `gorm:"default:0" json:"loss"`
	Draw       int `gorm:"default:0" json:"draw"`
}

type Friendships struct {
	gorm.Model
	// ผู้ส่งคำขอ
	UserID string `gorm:"size:14;index;not null" json:"user_id"`
	// ผู้รับคำขอ
	FriendID string `gorm:"size:14;index;not null" json:"friend_id"`

	// สถานะความสัมพันธ์: 0=PENDING, 1=ACCEPTED, 2=BLOCKED
	Status int `gorm:"default:0" json:"status"`

	// ทำ BelongsTo เพื่อให้ดึงข้อมูล User ออกมาดูได้ง่าย
	User   Users `gorm:"foreignKey:UserID" json:"-"`
	Friend Users `gorm:"foreignKey:FriendID" json:"-"`
}
