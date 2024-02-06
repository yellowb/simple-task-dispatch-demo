package global

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"log"
)

var redisCli *redis.Client
var redisInitErr error

func RedisInit(redisCfg RedisCfg) {
	// 创建Redis连接
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisCfg.Addr, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.Db,
	})
	// 测试连接是否成功
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logrus.Errorf("Failed to connect to Redis:%v", err)
		redisInitErr = err
	}
	log.Println("Connected to Redis:", pong)
	redisCli = redisClient
}

func GetRedisCli() *redis.Client {
	if redisInitErr != nil {
		return nil
	}
	return redisCli
}
