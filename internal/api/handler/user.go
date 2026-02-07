package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/repository"
)

type UserHandler struct {
	userRepo repository.UserRepo
}

func NewUserHandler(repo repository.UserRepo) *UserHandler {
	return &UserHandler{
		userRepo: repo,
	}
}

func (h *UserHandler) FindAllUser(c *gin.Context) {
	// 1. เรียกใช้งาน Repository
	users, err := h.userRepo.FindAllUser()

	// 2. ตรวจสอบ Error จาก Database
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user data",
		})
		return
	}

	// 3. ตรวจสอบว่ามีข้อมูลไหม (Optional)
	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No users found",
			"users":   []interface{}{}, // ส่ง Array ว่างกลับไป
		})
		return
	}

	// 4. ส่งข้อมูลกลับเมื่อสำเร็จ
	c.JSON(http.StatusOK, gin.H{
		"message": "Users retrieved successfully",
		"count":   len(users),
		"users":   users,
	})
}
