package middleware

import (
	"net/http"
	"strings"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/database" // เพิ่มอันนี้
	"github.com/AstralxOilx/Coding-Competition-Game/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. ดึง Token จาก Header (โค้ดเดิมของคุณ)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 2. ตรวจสอบ JWT (โค้ดเดิมของคุณ)
		token, err := util.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 3. ดึง Claims มาเพื่อเอา User ID
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// 1. ดึง userID แบบปลอดภัย
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
			c.Abort()
			return
		}

		// 2. ❌ ห้ามใช้: tokenInHand := claims["refresh_token"].(string)
		// ✅ ให้ใช้ "Comma ok" เพื่อเช็คว่ามีไหม ถ้าไม่มีก็ให้เป็นค่าว่างไป ไม่ให้ Panic
		tokenInHand, _ := claims["refresh_token"].(string)

		// 3. ดึง Session จาก Redis
		storedRefreshToken, err := database.RDB.Get(database.Ctx, "session:"+userID).Result()
		if err != nil {
			// ถ้าใน Redis ไม่มีค่า แปลว่า Session นี้ถูกลบไปแล้ว (อาจโดนเครื่องอื่นเตะ)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired or logged out"})
			c.Abort()
			return
		}

		// 4. ตรวจสอบว่าตรงกันไหม (Optional: ถ้า Access Token คุณไม่ได้เก็บ refresh_token ไว้
		// ให้ข้ามบรรทัดนี้ไปเลย หรือเช็คแค่ว่าใน Redis มีค่าก็พอ)
		if tokenInHand != "" && storedRefreshToken != tokenInHand {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Logged in from another device"})
			c.Abort()
			return
		}
		// 5. ตรวจสอบว่า Refresh Token ที่ User มี (ถ้าคุณเก็บใส่ JWT ไว้) ตรงกับ Redis ไหม
		// หรือในกรณีที่ง่ายที่สุด: ถ้ามีค่าใน Redis แปลว่า "ยังอยู่ในระบบ"
		// แต่ถ้าต้องการ Single Device จริงๆ ต้องเช็คว่า "Token ชุดนี้" คือชุดล่าสุดหรือไม่

		// --------------------------------------------------

		// เก็บข้อมูลลง Context ให้ Handler ใช้งานต่อ
		c.Set("user_id", userID)
		c.Set("role", claims["role"])

		c.Next()
	}
}
