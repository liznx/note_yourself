package db

import (
	"log"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init() {
	var err error
	//单机版(如果能使用阿里云redis的话，业务量很大也可以用单机方式)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	_, err = RedisClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
