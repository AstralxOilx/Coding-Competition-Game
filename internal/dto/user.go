package dto

import "time"

// SignupRequest ใช้สำหรับรับข้อมูลตอนสมัครสมาชิก
type SignupRequest struct {
	UserName    string `json:"user_name" binding:"required,min=4,max=20"`
	DisplayName string `json:"display_name" binding:"required,max=14"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

type SigninRequest struct {
	UserName string `json:"user_name" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=6"`
}

type UserInfo struct {
	DisplayName string `json:"display_name"`
	Role        int    `json:"role"`
}

type SigninResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	User         UserInfo `json:"user"`
}

type ProfileRequest struct {
	DisplayName string `json:"display_name" binding:"omitempty,max=14"`
	AvatarURL   string `json:"avatar_url" binding:"omitempty,url"`
}

// ✅ อันนี้ดีแล้ว: ใช้สำหรับส่งข้อมูลกลับไปแสดงผลเท่านั้น (Read-only)
// ห้ามส่งตัวแปรที่มี Type นี้เข้าไปใน s.repo.UpdateUserInfo เด็ดขาด!
type ProfileResponse struct {
	ID          string             `json:"id"`
	DisplayName string             `json:"display_name"`
	AvatarURL   string             `json:"avatar_url"`
	PlayerLevel int                `json:"player_level"`
	PlayerExp   int                `json:"player_exp"`
	Status      int                `json:"status"`
	LastLogin   *time.Time         `json:"last_login"`
	Ranks       []UserRankResponse `json:"ranks"` // ตัวปัญหาถ้าเอาไปใส่ใน GORM Updates
}

type UserRankResponse struct {
	ModeName   string  `json:"mode_name"`
	Rank       int     `json:"rank"`      // 0-13
	RankTier   int     `json:"rank_tier"` // 3, 2, 1
	RankPoint  int     `json:"rank_point"`
	WinRate    float64 `json:"win_rate"` // คำนวณให้ Frontend เลย
	TotalGames int     `json:"total_games"`
}
