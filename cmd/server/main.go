package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/config"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
)

const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func SetupSecurity(r *gin.Engine) {
	// 1. ตั้งค่า CORS (Cross-Origin Resource Sharing)
	// ช่วยป้องกันไม่ให้เว็บไซต์อื่นที่ไม่ได้อนุญาตมาเรียก API ของเรา
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://yourdomain.com"}, // เพิ่ม localhost ไว้เทสฝั่ง Frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 2. เพิ่ม Security Headers พื้นฐาน
	// ป้องกันการโจมตีประเภท XSS, Clickjacking และการเดาประเภทไฟล์ (MIME Sniffing)
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Next()
	})
}

// 3. ฟังก์ชัน Helper สำหรับจัดการ Error แบบปลอดภัย
// ไม่ส่งรายละเอียดเชิงลึกของระบบออกไปให้ Client เพื่อป้องกันการรั่วไหลของข้อมูลโครงสร้างระบบ
func HandleError(c *gin.Context, err error) {
	if err != nil {
		log.Printf("[SERVER ERROR]: %v", err)                // บันทึกลง Log ฝั่ง Server เท่านั้น
		c.JSON(500, gin.H{"error": "Internal server error"}) // ส่งข้อความกลางๆ กลับไป
		c.Abort()
		return
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err) // ถ้าโหลดไม่ได้ให้หยุดทันที จะได้รู้ว่าไฟล์หาย
	}

	config.LoadConfig()

	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("%s[CORE] INITIALIZING SYSTEM...%s\n", ColorPurple, ColorReset)
	fmt.Println(strings.Repeat("=", 60))

	// 1. Environment Variables
	godotenv.Load()

	// 2. Database Connection
	database.InitDatabase()
	database.InitRedis()

	// ---------------------------------------------------------
	// 3.1 Migrate Main Database (ตารางหลัก)
	// ---------------------------------------------------------
	fmt.Printf("%s[ORM]%s Syncing Main Schema... ", ColorCyan, ColorReset)
	// รวมทุก Model ของ Main DB ไว้ในคำสั่งเดียว
	// errMain := database.DB.AutoMigrate(
	// 	&model.Users{},
	// 	&model.UserRanks{},
	// 	&model.Friendships{},
	// )
	// if errMain != nil {
	// 	fmt.Printf("%s[FAILED]%s\n", ColorRed, ColorReset)
	// 	log.Fatalf("Main Migration error: %v", errMain)
	// }
	fmt.Printf("%s[DONE]%s\n", ColorGreen, ColorReset)

	// ---------------------------------------------------------
	// 3.2 Migrate Log Database (ตาราง Log)
	// ---------------------------------------------------------
	// fmt.Printf("%s[ORM]%s Syncing Log Schema...  ", ColorCyan, ColorReset)
	// // รวมทุก Model ของ Log DB ไว้ในคำสั่งเดียว
	// errLog := database.LogDB.AutoMigrate(
	// 	&model.LoginLog{},
	// 	&model.MatchLogs{},
	// )
	// if errLog != nil {
	// 	fmt.Printf("%s[FAILED]%s\n", ColorRed, ColorReset)
	// 	log.Fatalf("Log Migration error: %v", errLog)
	// }
	fmt.Printf("%s[DONE]%s\n", ColorGreen, ColorReset)

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%s⚡ STATUS: %sSERVER READY %s| MODE: %sDEVELOPMENT%s\n", ColorWhite, ColorGreen, ColorWhite, ColorYellow, ColorReset)
	fmt.Println(strings.Repeat("-", 60))

	// 1. ตั้งค่าโหมดก่อนเริ่มระบบ
	if os.Getenv("APP_MODE") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// r := gin.Default()
	r := router.InitRouter()
	// 2. เรียกใช้ Security Setup
	SetupSecurity(r)
	r.Run(":8080")

}
