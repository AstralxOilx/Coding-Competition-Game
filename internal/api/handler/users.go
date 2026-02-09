package handler

import (
	"net/http"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/service" // ✅ Import service
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService // ✅ เปลี่ยนจาก repo เป็น service
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{userService: s}
}

func (h *UserHandler) FindAllUser(c *gin.Context) {
	// 1. เรียกใช้งาน Service แทน Repo
	users, err := h.userService.AllUsers()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"count":   len(users),
		"users":   users,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	profile, err := h.userService.Profile(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile data retrieved",
		"data":    profile,
	})
}

func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	// 1. ดึง userID (เหมือนเดิม)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. รับ JSON (เหมือนเดิม)
	var input struct {
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// 3. เรียกใช้ Service อัปเดตข้อมูล
	// หมายเหตุ: ตรงนี้ถ้า Service คืนค่าเป็น *dto.ProfileResponse
	// เราก็แค่หยิบเฉพาะฟิลด์ที่ต้องการมาตอบครับ
	updatedUser, err := h.userService.UpdateUserInfo(userID.(string), input.DisplayName, input.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. ✅ ส่งกลับแค่ ชื่อ และ รูป ตามที่พี่ต้องการ
	c.JSON(http.StatusOK, gin.H{
		"message":      "Update successful",
		"display_name": updatedUser.DisplayName,
		"avatar_url":   updatedUser.AvatarURL,
	})
}
