package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient is a global variable for redis client.
var RedisClient *redis.Client

// ctx is a global variable for context.Background().
var ctx = context.Background()

// ConnectRedis creates a new Redis client and connects to a Redis server.
// It retries connecting to the Redis server if the connection failed, with a limit of maxRetries.
// It sleeps for retryInterval seconds between each retry.
// On successful connection, it returns a success message with the pong received from the Redis server.
func ConnectRedis() {
	const maxRetries = 3
	const retryInterval = 3
	var retryCount int

	for {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
		})

		// Try to ping the Redis server to check the connection.
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
