package database

import (
	"context"
	"fmt"
	"log"
	"server/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func RedisConfig() {
	const maxRetries = 3
	const retryInterval = 3
	var retryCount int

	for {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     config.Env.RedisHost,
			Password: config.Env.RedisPassword,
		})

		_, err := RedisClient.Ping(ctx).Result()
		if err != nil && retryCount < maxRetries {
			retryCount++
			log.Printf("🚧 Redis connection attempt %d failed\n", retryCount)
			if retryCount == maxRetries {
				log.Fatalf("💔 Redis connection failed after %d attempts\n", maxRetries)
			}
			time.Sleep(retryInterval * time.Second)
			continue
		}
		fmt.Println("🎉 Redis connection established successfully!")
		break
	}
}
