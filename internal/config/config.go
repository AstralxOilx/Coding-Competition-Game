package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	SessionDuration      time.Duration
}

var AppConfig *Config

func LoadConfig() {
	_ = godotenv.Load() // โหลดไฟล์ .env

	AppConfig = &Config{
		JWTSecret:            getEnv("JWT_SECRET", "default_secret"),
		AccessTokenDuration:  parseDuration(getEnv("ACCESS_TOKEN_DURATION", "15m")),
		RefreshTokenDuration: parseDuration(getEnv("REFRESH_TOKEN_DURATION", "168h")),
		SessionDuration:      parseDuration(getEnv("SESSION_DURATION", "168h")),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}
