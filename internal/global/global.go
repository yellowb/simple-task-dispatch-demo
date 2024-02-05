package global

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"github.com/yellowb/simple-task-dispatch-demo/internal/global/initialize"
)

// 全局变量
var (
	RedisCli *redis.Client
	RedSync  *redsync.Redsync
)

func GlobalInit(redisCfg RedisCfg) {
	initialize.RedisInit(redisCfg)
	RedisCli = initialize.GetRedisCli()
	initialize.RedSyncLockInit()
	RedSync = initialize.GetRedSync()
}

func GlobalClose() {
	RedisCli.Close()
}
