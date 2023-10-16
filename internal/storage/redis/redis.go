package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go_test/domains"
	"time"
)

type RDBClient struct {
	App     *domains.PersonProcessingApp
	RedisDB *redis.Client
}

func New(app *domains.PersonProcessingApp) *RDBClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     app.Cfg.RedisUrl,
		Password: app.Cfg.RedisPassword,
		DB:       app.Cfg.RedisDB,
	})
	return &RDBClient{
		App:     app,
		RedisDB: rdb,
	}
}

func (rClient *RDBClient) SetToCache(key string, person []byte) {
	rClient.RedisDB.Set(context.Background(), key, person, 10*time.Second)
}

func (rClient *RDBClient) GetFromCache(key string) ([]byte, error) {
	return rClient.RedisDB.Get(context.Background(), key).Bytes()
}
