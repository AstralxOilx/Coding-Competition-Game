package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global variables for database connections
var (
	DB    *gorm.DB // Main Database
	LogDB *gorm.DB // Log Database
)

const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorCyan   = "\033[36m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func InitDatabase() {
	// Configure GORM Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	// 1. Connect to Main Database
	DB = connectDB(
		os.Getenv("DB_MAIN_HOST"),
		os.Getenv("DB_MAIN_USER"),
		os.Getenv("DB_MAIN_PASSWORD"),
		os.Getenv("DB_MAIN_NAME"),
		os.Getenv("DB_MAIN_PORT"),
		os.Getenv("DB_MAIN_SSLMODE"),
		"MAIN DATABASE",
		newLogger,
	)

	// 2. Connect to Log Database
	LogDB = connectDB(
		os.Getenv("DB_LOG_HOST"),
		os.Getenv("DB_LOG_USER"),
		os.Getenv("DB_LOG_PASSWORD"),
		os.Getenv("DB_LOG_NAME"),
		os.Getenv("DB_LOG_PORT"),
		os.Getenv("DB_LOG_SSLMODE"),
		"LOG DATABASE",
		newLogger,
	)

	fmt.Println(ColorBlue + "---------------------------------------" + ColorReset)
}

// Helper function to handle connection logic for multiple DBs
func connectDB(host, user, pass, name, port, ssl, label string, gormLogger logger.Interface) *gorm.DB {
	fmt.Printf("%s[DB]%s Connecting to %s... ", ColorCyan, ColorReset, label)
	// --- เพิ่มบรรทัดนี้เพื่อเช็คค่าที่แอปอ่านได้จริง ---
	fmt.Printf("\n[DEBUG] Connecting to %s -> Host: %s, Port: %s, DB: %s, User: %s\n", label, host, port, name, user)
	// -------------------------------------------
	// ป้องกันค่าว่าง
	if ssl == "" {
		ssl = "disable"
	}

	// เปลี่ยนมาใช้รูปแบบ URL (Stronger & More Reliable for SASL/SCRAM)
	// รูปแบบ: postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=Asia/Bangkok",
		user, pass, host, port, name, ssl,
	)

	// ใช้ postgres.Open(dsn) เหมือนเดิม แต่ dsn เป็น URL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		fmt.Printf("%s[FAILED]%s\n", ColorRed, ColorReset)
		fmt.Printf("%sError:%s %v\n", ColorRed, ColorReset, err)
		// แสดง DSN ออกมาเช็ค (เฉพาะตอน Debug) เพื่อดูว่าค่าที่อ่านจาก .env ถูกไหม
		// fmt.Printf("DEBUG DSN: host=%s port=%s user=%s\n", host, port, user)
		panic(fmt.Sprintf("Failed to connect to %s", label))
	}

	// Connection Pool (เหมือนเดิม)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Printf("%s[SUCCESS]%s\n", ColorGreen, ColorReset)
	fmt.Printf("%s »%s Name: %s%s%s\n", ColorYellow, ColorReset, ColorCyan, name, ColorReset)

	return db
}
