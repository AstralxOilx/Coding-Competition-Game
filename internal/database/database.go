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

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		host, user, pass, name, port, ssl,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		fmt.Printf("%s[FAILED]%s\n", ColorRed, ColorReset)
		fmt.Printf("%sError:%s %v\n", ColorRed, ColorReset, err)
		panic(fmt.Sprintf("Failed to connect to %s", label))
	}

	// Connection Pool Settings
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Printf("%s[SUCCESS]%s\n", ColorGreen, ColorReset)
	fmt.Printf("%s »%s Name: %s%s%s\n", ColorYellow, ColorReset, ColorCyan, name, ColorReset)

	return db
}
