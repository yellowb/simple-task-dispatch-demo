package initialize

import (
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"time"
)

var redSync *redsync.Redsync

func RedSyncLockInit() {
	// 创建redsync的客户端连接池
	pool := goredis.NewPool(redisCli)
	// 创建redsync实例
	redSync = redsync.New(pool)
}

func GetRedSync() *redsync.Redsync {
	if redSync == nil {
		RedSyncLockInit()
	}
	return redSync
}

func CreateLock(lockName string, timeout time.Duration, retries int) *redsync.Mutex {
	if redSync == nil {
		RedSyncLockInit()
	}
	//创建基于key的互斥锁
	mutex := redSync.NewMutex(lockName, redsync.WithExpiry(timeout), redsync.WithTries(retries))
	return mutex
}
