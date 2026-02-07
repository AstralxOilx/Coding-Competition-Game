package handler

import (
	"net/http"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/config"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/model"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository/cache/redis"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userRepo repository.UserRepo
}

func NewAuthHandler(repo repository.UserRepo) *AuthHandler {
	return &AuthHandler{userRepo: repo}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req model.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Hash Password
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 2. เตรียม Model
	user := model.Users{
		ID:          util.GenerateID(14),
		DisplayName: req.DisplayName,
		UserName:    req.UserName,
		Email:       req.Email,
		Password:    hashedPassword,
	}

	// 3. บันทึกผ่าน Repository
	if err := h.userRepo.CreateUser(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or Username already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user_id": user.ID})
}

func (h *AuthHandler) Signin(c *gin.Context) {

	var req model.SigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	user, err := h.userRepo.FindByUserName(req.UserName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	token, err := util.GenerateAccessToken(user.ID, user.UserRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	refreshToken, err := util.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Authentication failed"})
		return
	}

	oldRefreshToken, err := redis.GetUserSession(c, user.ID)

	if err == nil && oldRefreshToken != "" {
		// 1. เตะฝั่ง WebSocket: สั่งปิด Connection ของเครื่องเก่าทันที
		util.WSManager.NotifyOldDevice(user.ID)

		// 2. เตะฝั่ง API: ลบ Session ใน Redis ทิ้ง
		// ทำให้คนเก่าที่ถือ Token เดิมอยู่ เมื่อเรียก API จะติด AuthMiddleware (401) ทันที
		_ = database.RDB.Del(database.Ctx, "session:"+user.ID).Err()

		// (ทางเลือก) ถ้าต้องการให้คนใหม่ต้องล็อกอินซ้ำอีกรอบเพื่อยืนยัน
		c.JSON(http.StatusConflict, gin.H{
			"error": "Account logged in on another device. Old sessions terminated. Please login again.",
		})
		return
	}

	err = redis.SetUserSession(c, user.ID, refreshToken, config.AppConfig.SessionDuration)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login Success",
		"token":         token,
		"refresh_token": refreshToken,
		"user": gin.H{
			"display_name": user.DisplayName,
			"role":         user.UserRole,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	// 1. ตรวจสอบความถูกต้องของ Refresh Token (ใช้ util.ValidateToken ตัวเดิม)
	token, err := util.ValidateRefreshToken(input.RefreshToken)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// 2. ดึง user_id จาก claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
		return
	}

	userID, ok := claims["user_id"].(string)

	// 3. (Optional) ไปดึง Role ปัจจุบันจาก Database
	user, _ := h.userRepo.FindById(userID)

	// 4. สร้าง Access Token ชุดใหม่
	newAccessToken, _ := util.GenerateAccessToken(user.ID, user.UserRole)
	newRefreshToken, _ := util.GenerateRefreshToken(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
