package helpers

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		
		poolSize, _ := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
		minIdleConns, _ := strconv.Atoi(os.Getenv("REDIS_MIN_IDLE_CONNS"))

		rdbInstance = redis.NewClient(&redis.Options{
				Addr:         os.Getenv("REDIS_CONNECTION"),
				Password:     os.Getenv("REDIS_PASSWORD"),
				DB:           0, 
				PoolSize:     poolSize,
				MinIdleConns: minIdleConns,
				PoolTimeout:  60000,
				// IdleTimeout:  60,
			})

		
		logrus.Info("[Redis] Connection OK")

	})
	return rdbInstance
}


func TestRedisConnection() {
	rdb := GetRedisConnection()
	err := rdb.Set(context.Background(), "healthcheck." + os.Getenv("SERVER_NAME") , "ok." + time.Now().String() , 0).Err()
    ExitOnError(err, "[Redis] Failed to Test Redis")	
}


func SetPessoa(apelido string, id string, termo string, pessoaInput string) error {
	rdb := GetRedisConnection()
	ctx := context.Background()

	_, err := rdb.Get(ctx, "pessoa|" + apelido).Result()
	if err == redis.Nil {
		rdb.Set(ctx, "pessoa|" + apelido, pessoaInput, 15*time.Minute)
		rdb.Set(ctx, "pessoa|" + id, pessoaInput, 15*time.Minute)
		rdb.Set(ctx, "pessoa|" + termo, pessoaInput, 15*time.Minute)
		return nil
	} else {
		return errors.New("pessoa ja cadastrada")
	}

}

func GetPessoaById(id string) (string, error) {
	rdb := GetRedisConnection()
	ctx := context.Background()

	jsonData, err := rdb.Get(ctx, "pessoa|" + id).Result()
	if err == redis.Nil {
		return "", errors.New("pessoa nao encontrada")
	} else {
		return jsonData, nil
	}

}


func GetPessoaByTermo(t string) ([]string, error) {
	rdb := GetRedisConnection()
	ctx := context.Background()

	var result []string
	var cursor uint64
	var err error
	pattern := "*" +   strings.ReplaceAll(strings.ToLower(t), " ", "") + "*"

	// logrus.Debug(">> Termo:", pattern)

	keys := make([]string, 0)
	count := 1
	for {
		var batch []string
		batch, cursor, err = rdb.Scan(ctx, cursor, pattern, 0).Result()

		LogOnError(err, fmt.Sprintf("Erro ao realizar SCAN: %v", err))

		keys = append(keys, batch...)

		if cursor == 0 {
			break
		}
		count++
		if count >= 10 {
			break
		}
	}


	for _, key := range keys {
		value, err := rdb.Get(ctx, key).Result()
		LogOnError(err, fmt.Sprintf("Erro ao obter valor para a chave %s: %v", key, err))

		result = append(result, value)
		// fmt.Printf("Chave: %s, Valor: %s\n", key, value)
	}

	return result, nil
}


