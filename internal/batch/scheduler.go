package batch

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func StartScheduler(db *gorm.DB) *cron.Cron {
	c := cron.New()
	c.AddJob("*/5 * * * *", NewWalkAlarmJob(db))
	c.AddJob("0 23 * * *", NewPetItemJob(db))
	c.Start()
	return c
}
