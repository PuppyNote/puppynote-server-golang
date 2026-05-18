package main

import (
	"puppynote/internal/alert"
	"puppynote/internal/auth"
	"puppynote/internal/community"
	"puppynote/internal/family"
	"puppynote/internal/food"
	"puppynote/internal/home"
	"puppynote/internal/pet"
	"puppynote/internal/petitem"
	"puppynote/internal/pettip"
	"puppynote/internal/storage"
	"puppynote/internal/user"
	"puppynote/internal/walk"
	"puppynote/internal/walkalarm"
	"puppynote/internal/weather"
	"puppynote/pkg/database"
	redispkg "puppynote/pkg/redis"
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

	// Pet
	petRepo := pet.NewRepository(database.DB)
	petSvc := pet.NewService(petRepo)
	pet.NewHandler(petSvc).RegisterRoutes(base)

	// FamilyMember
	familyRepo := family.NewRepository(database.DB)
	familySvc := family.NewService(familyRepo)
	family.NewHandler(familySvc).RegisterRoutes(base)

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
