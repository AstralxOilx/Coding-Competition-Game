package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
	"github.com/AstralxOilx/Coding-Competition-Game/internal/model"
	"github.com/joho/godotenv"
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

func main() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("%s[CORE] INITIALIZING SYSTEM...%s\n", ColorPurple, ColorReset)
	fmt.Println(strings.Repeat("=", 60))

	// 1. Environment Variables
	fmt.Printf("%s[ENV]%s Loading configuration... ", ColorBlue, ColorReset)
	if err := godotenv.Load(); err != nil {
		fmt.Printf("%s[MISSING]%s Using system environment variables\n", ColorYellow, ColorReset)
	} else {
		fmt.Printf("%s[SUCCESS]%s\n", ColorGreen, ColorReset)
	}

	// 2. Database Connection (Main & Log)
	database.InitDatabase()

	// 3. Database Migration
	// 3.1 Migrate Main Database
	fmt.Printf("%s[ORM]%s Syncing Main Schema... ", ColorCyan, ColorReset)
	if err := database.DB.AutoMigrate(&model.Users{}); err != nil {
		fmt.Printf("%s[ERROR]%s\n", ColorRed, ColorReset)
		log.Fatalf("Main Migration failed: %v", err)
	}
	fmt.Printf("%s[DONE]%s\n", ColorGreen, ColorReset)

	// 3.2 Migrate Log Database
	fmt.Printf("%s[ORM]%s Syncing Log Schema...  ", ColorCyan, ColorReset)
	if err := database.LogDB.AutoMigrate(&model.LoginLog{}); err != nil {
		fmt.Printf("%s[ERROR]%s\n", ColorRed, ColorReset)
		log.Fatalf("Log Migration failed: %v", err)
	}
	fmt.Printf("%s[DONE]%s\n", ColorGreen, ColorReset)

	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%s⚡ STATUS: %sSERVER READY %s| MODE: %sDEVELOPMENT%s\n", ColorWhite, ColorGreen, ColorWhite, ColorYellow, ColorReset)
	fmt.Println(strings.Repeat("-", 60))

	// Start your Server (Gin/Echo) here...
}
