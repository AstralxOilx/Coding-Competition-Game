package redis

import (
	"context"
	"time"

	"github.com/AstralxOilx/Coding-Competition-Game/internal/database"
)

type RedisCache struct {
	// เก็บ dependencies ถ้าจำเป็น
}

func SetUserSession(ctx context.Context, userID string, token string, duration time.Duration) error {
	key := "session:" + userID
	return database.RDB.Set(ctx, key, token, duration).Err()
}

func GetUserSession(ctx context.Context, userID string) (string, error) {
	key := "session:" + userID
	return database.RDB.Get(ctx, key).Result()
}

func DeleteSession(ctx context.Context, userID string) error {
	key := "session:" + userID
	return database.RDB.Del(ctx, key).Err()
}
