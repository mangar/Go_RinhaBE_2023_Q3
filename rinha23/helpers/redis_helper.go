package helpers

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	onceRedis sync.Once
	rdbInstance *redis.Client
)


func GetRedisConnection() (*redis.Client) {
	onceRedis.Do(func() {
		opts, err := redis.ParseURL(os.Getenv("REDIS_CONNECTION"))
		ExitOnError(err, "[Redis] Failed to connect to Redis")	
		logrus.Info("[Redis] Connection OK")
		
		rdbInstance = redis.NewClient(opts)

	})
	return rdbInstance
}


func TestRedisConnection() {
	rdb := GetRedisConnection()
	err := rdb.Set(context.Background(), "healthcheck." + os.Getenv("SERVER_NAME") , "ok." + time.Now().String() , 0).Err()
    ExitOnError(err, "[Redis] Failed to Test Redis")	
}