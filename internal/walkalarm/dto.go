package walkalarm

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/model"
)

type WalkAlarmCreateRequest struct {
	PetID       int64               `json:"petId" binding:"required"`
	AlarmStatus model.AlarmStatus   `json:"alarmStatus" binding:"required"`
	AlarmDays   []model.AlarmDay    `json:"alarmDays" binding:"required,min=1"`
	AlarmTime   string              `json:"alarmTime" binding:"required"` // HH:mm
}

type WalkAlarmUpdateRequest struct {
	AlarmID     int64               `json:"alarmId" binding:"required"`
	AlarmStatus model.AlarmStatus   `json:"alarmStatus" binding:"required"`
	AlarmDays   []model.AlarmDay    `json:"alarmDays" binding:"required,min=1"`
	AlarmTime   string              `json:"alarmTime" binding:"required"`
}

type WalkAlarmStatusUpdateRequest struct {
	AlarmID     int64             `json:"alarmId" binding:"required"`
	AlarmStatus model.AlarmStatus `json:"alarmStatus" binding:"required"`
}

type WalkAlarmResponse struct {
	AlarmID     int64            `json:"alarmId"`
	AlarmStatus model.AlarmStatus `json:"alarmStatus"`
	AlarmDays   []model.AlarmDay `json:"alarmDays"`
	AlarmTime   string           `json:"alarmTime"`
}
