package main

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/auth"
	"github.com/PuppyNote/puppynote-server-golang/internal/user"
	"github.com/PuppyNote/puppynote-server-golang/pkg/database"
	"github.com/gin-gonic/gin"
)

func registerRoutes(base *gin.RouterGroup) {
	// Auth
	authRepo := auth.NewRepository(database.DB)
	authSvc := auth.NewService(authRepo)
	auth.NewHandler(authSvc).RegisterRoutes(base)

	// User
	userRepo := user.NewRepository(database.DB)
	userSvc := user.NewService(userRepo)
	user.NewHandler(userSvc).RegisterRoutes(base)
}
