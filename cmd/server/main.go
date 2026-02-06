package main

import (
	"log"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/model"
	"github.com/joho/godotenv"
)

func main() {
	// 1. โหลดไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: ไม่พบไฟล์ .env จะใช้ค่าจาก System Environment")
	}

	// 2. เริ่มการเชื่อมต่อ DB
	database.InitDatabase()

	// 3. ทำการ Auto Migrate (สร้าง/อัปเดตตาราง)
	err := database.DB.AutoMigrate(
		&model.Users{},
		&model.LoginLog{}, // ตารางบันทึก log ที่คุยกันก่อนหน้า
	)
	if err != nil {
		log.Fatalf("Migration ล้มเหลว: %v", err)
	}

	log.Println("🚀 ระบบพร้อมทำงาน...")

	// ส่วนของ Server (เช่น Gin หรือ Echo) จะเริ่มตรงนี้...
}
