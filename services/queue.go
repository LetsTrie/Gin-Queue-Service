package service

import (
	"context"

	db "server/database"
)

var ctx = context.Background()

type RedisQueue struct {
	Key string
}

func (q *RedisQueue) Set(item string) error {
	return db.RedisClient.LPush(ctx, q.Key, item).Err()
}

func (q *RedisQueue) Get() (string, error) {
	return db.RedisClient.RPop(ctx, q.Key).Result()
}

func NewRedisQueue(key string) *RedisQueue {
	return &RedisQueue{
		Key: key,
	}
}
