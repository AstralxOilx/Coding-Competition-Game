package util

import (
	"fmt"
	"os"
	"time"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte(os.Getenv("AUTHKEY")) // ในงานจริงควรดึงจาก .env
var refresh_key = []byte(os.Getenv("REFRESH_KEY"))

func GenerateAccessToken(userID string, role int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(config.AppConfig.AccessTokenDuration).Unix(), // หมดอายุใน 2 ชม.
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID, // เก็บเป็น string ลงใน Token
		"exp":     time.Now().Add(config.AppConfig.RefreshTokenDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refresh_key)
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่า Algorithm การเข้ารหัสถูกต้องไหม (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}

func ValidateRefreshToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบ Algorithm (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// สำคัญ: ต้องใช้ refresh_key ในการตรวจ (คนละตัวกับ secretKey)
		return refresh_key, nil
	})
}
