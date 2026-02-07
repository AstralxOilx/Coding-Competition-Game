package util

import (
	"log"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

func GenerateID(length int) string {
	// กำหนดเฉพาะตัวอักษรและตัวเลขที่ต้องการ (ไม่มีอักขระพิเศษ)
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// ใช้ฟังก์ชัน Generate เพื่อระบุ alphabet และความยาว
	id, err := gonanoid.Generate(alphabet, length)
	if err != nil {
		log.Printf("[Util] Failed to generate ID: %v", err)
		return ""
	}
	return id
}
