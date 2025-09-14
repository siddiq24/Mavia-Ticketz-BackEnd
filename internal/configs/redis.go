package configs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client, error) {
	addr := os.Getenv("ADDR")
	port := os.Getenv("PORT")
	// user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")

	redisAddr := fmt.Sprintf("%s:%s", addr, port)

	// init redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		// Username: user,
		Password: password,
		// DB:       0,
	})

	// test connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", redisAddr, err)
	}

	log.Println("Connected to Redis at", redisAddr)
	return redisClient, nil
}
