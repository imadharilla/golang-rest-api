package redisdb

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})



func SetValue(key , value string) error {
	// cache for 1 hour
	err := client.Set(ctx, key, value, 60*time.Minute).Err()
    return err
}

func GetValue( key string) (string,error) {
	val, err := client.Get(ctx, key).Result()
    return val, err
}