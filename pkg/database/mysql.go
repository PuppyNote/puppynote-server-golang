package database

import (
	"fmt"
	"log"

	"puppynote/config"
	"puppynote/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	cfg := config.AppConfig.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("DB 인스턴스 획득 실패: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err = autoMigrate(); err != nil {
		log.Fatalf("AutoMigrate 실패: %v", err)
	}

	log.Println("DB 연결 성공")
}

func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.RefreshToken{},
		&model.Push{},
		&model.Pet{},
		&model.FamilyMember{},
		&model.Walk{},
		&model.WalkPhoto{},
		&model.PetWalkAlarm{},
		&model.PetAlarmDay{},
		&model.PetItem{},
		&model.PetItemPurchase{},
		&model.Post{},
		&model.PostHashtag{},
		&model.PostImage{},
		&model.PostLike{},
		&model.FoodChatHistory{},
		&model.AlertSetting{},
		&model.AlertHistory{},
		&model.PetTip{},
	)
}
