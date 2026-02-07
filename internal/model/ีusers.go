package model

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID          string     `gorm:"primaryKey;size:14" json:"id"`
	DisplayName string     `gorm:"not null" json:"display_name"`
	AvatarURL   string     `json:"avatar_url"`
	UserName    string     `gorm:"not null" json:"user_name"`
	Password    string     `gorm:"not null" json:"password"`
	Status      string     `gorm:"not null" json:"status"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	LastLogin   *time.Time `json:"last_login"` // เก็บสถิติการเข้าใช้งานล่าสุด
	gorm.Model
}
