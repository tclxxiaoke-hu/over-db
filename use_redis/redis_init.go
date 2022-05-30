package use_redis

import (
	"github.com/go-redis/redis"
)

var redisDb *redis.Client

//InitClient 初始化连接
func InitClient() (err error) {
	redisDb = redis.NewClient(&redis.Options{
		Addr:     "192.168.3.8:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = redisDb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
