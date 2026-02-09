package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	Redis *redis.Client
}

func (r *UserCache) SetUserOnline(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user:online: %s", userID)
	return r.Redis.Set(ctx, key, "online", 60*time.Second).Err()
}

func (r *UserCache) GetUserOnline(ctx context.Context, userID string) (bool, error) {
	key := fmt.Sprintf("user:online: %s", userID)
	val, err := r.Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}

	return val == "online", err
}
