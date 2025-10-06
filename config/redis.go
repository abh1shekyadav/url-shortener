package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Unable to connect to redis: %v", err)
	}
	log.Println("Connected to redis")
}
