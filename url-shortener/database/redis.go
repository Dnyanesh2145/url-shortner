package database

import (
	"context"
	"fiber-url-shortner/config"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisURL string = config.EnvDBURI("REDIS_URI")
var redisClient *redis.Client
var ctx = context.Background()

func RedisConnect() {
	var err error
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		panic("Error parsing redis uri: " + err.Error())
	}
	redisClient = redis.NewClient(opts)

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func RedisClose() {
	err := redisClient.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func SetHashkey(key string, field string, value string) error {
	err := redisClient.HSet(ctx, key, field, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetHashValue(key string, field string) (string, error) {
	value, err := redisClient.HGet(ctx, key, field).Result()
	if err != nil {
		return value, err
	}
	return value, nil
}

func RedisPing() (string, error) {
	return redisClient.Ping(ctx).Result()
}
