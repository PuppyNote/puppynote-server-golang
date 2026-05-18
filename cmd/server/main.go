package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuppyNote/puppynote-server-golang/config"
	"github.com/PuppyNote/puppynote-server-golang/pkg/database"
	"github.com/PuppyNote/puppynote-server-golang/pkg/middleware"
	"github.com/PuppyNote/puppynote-server-golang/pkg/redis"
	"github.com/PuppyNote/puppynote-server-golang/pkg/response"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()

	database.Connect()
	redis.Connect()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	base := r.Group(config.AppConfig.Server.ContextPath)
	{
		base.GET("/health-check", func(c *gin.Context) {
			response.OK(c, "서버가 정상 동작중입니다.")
		})
	}

	registerRoutes(base)

	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("PuppyNote 서버 시작: http://localhost%s%s", addr, config.AppConfig.Server.ContextPath)
	if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
