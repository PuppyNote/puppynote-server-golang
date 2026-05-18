package main

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/user"
	"github.com/PuppyNote/puppynote-server-golang/pkg/database"
	"github.com/gin-gonic/gin"
)

func registerRoutes(base *gin.RouterGroup) {
	userRepo := user.NewRepository(database.DB)
	userSvc := user.NewService(userRepo)
	user.NewHandler(userSvc).RegisterRoutes(base)
}
