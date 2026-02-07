package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึงค่าจาก Context (ที่ถูก Set มาจาก AuthMiddleware)
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required: role not found"})
			c.Abort()
			return
		}

		// 2. แปลงประเภทข้อมูลให้เป็น int อย่างปลอดภัย (Safe Type Conversion)
		// JSON numbers ใน JWT มักจะกลายเป็น float64 เมื่อถูกแกะออกมา
		var userRole int
		switch v := roleValue.(type) {
		case int:
			userRole = v
		case float64:
			userRole = int(v)
		case int64:
			userRole = int(v)
		default:
			// หากเป็น Type อื่นที่ไม่คาดคิด ให้แจ้ง Error และ Print ดู Type จริง
			fmt.Printf("[RoleMiddleware] Unexpected type: %T for value: %v\n", roleValue, roleValue)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: invalid role data format"})
			c.Abort()
			return
		}

		// 3. ตรวจสอบว่าสิทธิ์ (Role) ตรงกับที่อนุญาตหรือไม่
		isAllowed := false
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				isAllowed = true
				break
			}
		}

		// 4. ถ้าไม่มีสิทธิ์ ให้ส่ง 403 Forbidden
		if !isAllowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: insufficient permissions"})
			c.Abort()
			return
		}

		// ผ่านฉลุย! ให้ไปทำ Handler ถัดไป
		c.Next()
	}
}
