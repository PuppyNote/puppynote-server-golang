package main

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/auth"
	"github.com/PuppyNote/puppynote-server-golang/internal/pet"
	"github.com/PuppyNote/puppynote-server-golang/internal/petitem"
	"github.com/PuppyNote/puppynote-server-golang/internal/user"
	"github.com/PuppyNote/puppynote-server-golang/internal/walk"
	"github.com/PuppyNote/puppynote-server-golang/internal/walkalarm"
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

	// Pet & FamilyMember
	petRepo := pet.NewRepository(database.DB)
	petSvc := pet.NewService(petRepo)
	pet.NewHandler(petSvc).RegisterRoutes(base)

	// Walk
	walkRepo := walk.NewRepository(database.DB)
	walkSvc := walk.NewService(walkRepo)
	walk.NewHandler(walkSvc).RegisterRoutes(base)

	// WalkAlarm
	walkAlarmRepo := walkalarm.NewRepository(database.DB)
	walkAlarmSvc := walkalarm.NewService(walkAlarmRepo)
	walkalarm.NewHandler(walkAlarmSvc).RegisterRoutes(base)

	// PetItem & PetItemPurchase
	petItemRepo := petitem.NewRepository(database.DB)
	petItemSvc := petitem.NewService(petItemRepo)
	petitem.NewHandler(petItemSvc).RegisterRoutes(base)
}
