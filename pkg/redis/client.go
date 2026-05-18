package redis

import (
	"context"
	"fmt"
	"log"

	"puppynote/config"
	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func Connect() {
	cfg := config.AppConfig.Redis
	Client = goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})

	if err := Client.Ping(context.Background()).Err(); err != nil {
		log.Printf("Redis 연결 실패 (선택적): %v", err)
		return
	}
	log.Println("Redis 연결 성공")
}
