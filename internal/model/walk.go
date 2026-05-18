package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type AlarmStatus string
type AlarmDay string

const (
	AlarmStatusYes AlarmStatus = "YES"
	AlarmStatusNo  AlarmStatus = "NO"

	AlarmDayMon AlarmDay = "MON"
	AlarmDayTue AlarmDay = "TUE"
	AlarmDayWed AlarmDay = "WED"
	AlarmDayThu AlarmDay = "THU"
	AlarmDayFri AlarmDay = "FRI"
	AlarmDaySat AlarmDay = "SAT"
	AlarmDaySun AlarmDay = "SUN"
)

type Walk struct {
	ID         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PetID      int64          `gorm:"column:pet_id;not null" json:"petId"`
	Pet        Pet            `gorm:"foreignKey:PetID" json:"-"`
	StartTime  time.Time      `gorm:"column:start_time;not null;index:idx_walks_start_time" json:"startTime"`
	EndTime    time.Time      `gorm:"column:end_time;not null" json:"endTime"`
	Latitude   *float64       `gorm:"column:latitude;type:decimal(11,7)" json:"latitude"`
	Longitude  *float64       `gorm:"column:longitude;type:decimal(11,7)" json:"longitude"`
	Location   string         `gorm:"column:location;size:100" json:"location"`
	Memo       string         `gorm:"column:memo;size:500" json:"memo"`
	Photos     []WalkPhoto    `gorm:"foreignKey:WalkID" json:"photos,omitempty"`
	BaseTimeEntity
}

func (Walk) TableName() string { return "walks" }

type PetWalkAlarm struct {
	ID          int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	PetID       int64       `gorm:"column:pet_id;not null" json:"petId"`
	Pet         Pet         `gorm:"foreignKey:PetID" json:"-"`
	AlarmStatus AlarmStatus `gorm:"column:alarm_status;type:varchar(3);not null" json:"alarmStatus"`
	AlarmTime   LocalTime   `gorm:"column:alarm_time;not null" json:"alarmTime"`
	AlarmDays   []PetAlarmDay `gorm:"foreignKey:PetAlarmID" json:"alarmDays,omitempty"`
	BaseTimeEntity
}

func (PetWalkAlarm) TableName() string { return "pet_walk_alarms" }

type PetAlarmDay struct {
	PetAlarmID int64    `gorm:"column:pet_alarm_id;not null" json:"petAlarmId"`
	AlarmDay   AlarmDay `gorm:"column:alarm_day;type:varchar(3);not null" json:"alarmDay"`
}

func (PetAlarmDay) TableName() string { return "pet_alarm_days" }

// LocalTime은 time.Time을 HH:MM:SS 형태로 DB에 저장하기 위한 커스텀 타입
type LocalTime struct {
	time.Time
}

func (lt LocalTime) Value() (driver.Value, error) {
	return lt.Format("15:04:05"), nil
}

func (lt *LocalTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		lt.Time = v
		return nil
	case []byte:
		t, err := time.Parse("15:04:05", string(v))
		if err != nil {
			return err
		}
		lt.Time = t
		return nil
	case string:
		t, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		lt.Time = t
		return nil
	}
	return fmt.Errorf("unsupported type: %T", value)
}

func (lt LocalTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + lt.Format("15:04") + `"`), nil
}
