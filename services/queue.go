package service

import (
	"context"

	database "server/database"
)

var ctx = context.Background()

type RedisQueue struct {
	Key string
}

func (q *RedisQueue) Set(item string) error {
	return database.RedisClient.LPush(ctx, q.Key, item).Err()
}

func (q *RedisQueue) Get() (string, error) {
	return database.RedisClient.RPop(ctx, q.Key).Result()
}

func NewRedisQueue(key string) *RedisQueue {
	return &RedisQueue{
		Key: key,
	}
}
