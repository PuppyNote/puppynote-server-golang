package main

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/alert"
	"github.com/PuppyNote/puppynote-server-golang/internal/auth"
	"github.com/PuppyNote/puppynote-server-golang/internal/community"
	"github.com/PuppyNote/puppynote-server-golang/internal/food"
	"github.com/PuppyNote/puppynote-server-golang/internal/home"
	"github.com/PuppyNote/puppynote-server-golang/internal/pet"
	"github.com/PuppyNote/puppynote-server-golang/internal/petitem"
	"github.com/PuppyNote/puppynote-server-golang/internal/pettip"
	"github.com/PuppyNote/puppynote-server-golang/internal/storage"
	"github.com/PuppyNote/puppynote-server-golang/internal/user"
	"github.com/PuppyNote/puppynote-server-golang/internal/walk"
	"github.com/PuppyNote/puppynote-server-golang/internal/walkalarm"
	"github.com/PuppyNote/puppynote-server-golang/internal/weather"
	"github.com/PuppyNote/puppynote-server-golang/pkg/database"
	redispkg "github.com/PuppyNote/puppynote-server-golang/pkg/redis"
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

	// Community
	communityRepo := community.NewRepository(database.DB)
	communitySvc := community.NewService(communityRepo, redispkg.Client)
	community.NewHandler(communitySvc).RegisterRoutes(base)

	// Storage (S3)
	storage.NewHandler().RegisterRoutes(base)

	// Alert Setting & History
	alertRepo := alert.NewRepository(database.DB)
	alert.NewHandler(alertRepo).RegisterRoutes(base)

	// Weather
	weather.NewHandler().RegisterRoutes(base)

	// Food AI
	food.NewHandler(database.DB).RegisterRoutes(base)

	// Home
	home.NewHandler(database.DB).RegisterRoutes(base)

	// PetTip
	pettip.NewHandler(database.DB).RegisterRoutes(base)
}
