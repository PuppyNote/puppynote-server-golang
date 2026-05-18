package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuppyNote/puppynote-server-golang/config"
	"github.com/PuppyNote/puppynote-server-golang/internal/batch"
	"github.com/PuppyNote/puppynote-server-golang/pkg/database"
	espkg "github.com/PuppyNote/puppynote-server-golang/pkg/elasticsearch"
	pnjwt "github.com/PuppyNote/puppynote-server-golang/pkg/jwt"
	"github.com/PuppyNote/puppynote-server-golang/pkg/middleware"
	"github.com/PuppyNote/puppynote-server-golang/pkg/redis"
	"github.com/PuppyNote/puppynote-server-golang/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	config.Load()
	pnjwt.Init(config.AppConfig.JWT.SecretKey)

	database.Connect()
	redis.Connect()
	espkg.Connect()

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.PrometheusMiddleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	base := r.Group(config.AppConfig.Server.ContextPath)
	{
		base.GET("/health-check", func(c *gin.Context) {
			response.OK(c, "서버가 정상 동작중입니다.")
		})
	}

	registerRoutes(base)

	c := batch.StartScheduler(database.DB)
	defer c.Stop()

	addr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("PuppyNote 서버 시작: http://localhost%s%s", addr, config.AppConfig.Server.ContextPath)
	if err := r.Run(addr); err != nil && err != http.ErrServerClosed {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}
