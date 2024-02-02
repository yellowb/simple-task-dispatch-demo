package global

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var redisCli *redis.Client
var redisInitErr error

func NewRedisClient(redisCfg RedisCfg) {
	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisCfg.Addr, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.Db,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		redisInitErr = err
		return
	}
	redisCli = cli
}

func GetRedisCli() *redis.Client {
	if redisInitErr != nil {
		return nil
	}
	return redisCli
}
