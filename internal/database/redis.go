package database

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	// addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	password := os.Getenv("REDIS_PASSWORD")

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv(password), // ถ้าไม่มีจะเป็นค่าว่าง ""
		DB:       0,                   // Default DB
	})

	fmt.Printf("\033[34m[REDIS]\033[0m Connecting to Redis... ")
	if _, err := RDB.Ping(Ctx).Result(); err != nil {
		fmt.Printf("\033[31m[FAILED]\033[0m\n")
		// ในช่วงพัฒนาอาจจะไม่ต้อง Fatalf เพื่อให้แอปยังรันต่อได้ถ้าไม่มี Redis
		fmt.Printf("Warning: Redis not connected: %v\n", err)
	} else {
		fmt.Printf("\033[32m[SUCCESS]\033[0m\n")
	}

}
