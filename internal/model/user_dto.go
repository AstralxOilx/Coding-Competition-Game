package model

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
