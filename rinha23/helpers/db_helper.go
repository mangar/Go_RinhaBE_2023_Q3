package helpers

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	once sync.Once
	pool *pgxpool.Pool
)

func GetDBConnection() *pgxpool.Pool {
	once.Do(func() {

		dbConfig, err := pgxpool.ParseConfig(os.Getenv("DB_CONNECTION"))
		ExitOnError(err, "[DB] Failed to create a config")
		
		dbConfig.MaxConns = 25
		dbConfig.MinConns = 2
		dbConfig.MaxConnLifetime = time.Hour
		dbConfig.MaxConnIdleTime = time.Minute * 30
		dbConfig.HealthCheckPeriod = time.Minute
		dbConfig.ConnConfig.ConnectTimeout = time.Second * 5
	
		pool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		ExitOnError(err, "[DB] Unable to create connection pool")
		// defer pool.Close()
	
		err = pool.Ping(context.Background())
		ExitOnError(err, "[DB] Unable to ping database")

	})
	return pool
}

func TestDBConnection() {
	GetDBConnection()
	logrus.Info("[DB] Connection OK")
}

