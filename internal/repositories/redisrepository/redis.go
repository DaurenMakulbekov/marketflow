package redisrepository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

type redisRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisRepository(rdb *redis.Client, ctx context.Context) *redisRepository {
	var newRedisRepository = redisRepository{
		rdb: rdb,
		ctx: ctx,
	}

	return &newRedisRepository
}

func (redisRepo *redisRepository) Write(exchange string) {
	_, err := redisRepo.rdb.LPush(redisRepo.ctx, "exchange", exchange).Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to add key-value pair")
	}
}

func (redisRepo *redisRepository) Read() string {
	result, err := redisRepo.rdb.RPop(redisRepo.ctx, "exchange").Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Key not found in Redis cache:", err)
	}

	return result
}

func (redisRepo *redisRepository) LLen() int {
	length, err := redisRepo.rdb.LLen(redisRepo.ctx, "exchange").Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, "LLen error:", err)
	}

	return int(length)
}
