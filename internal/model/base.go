package model

import "time"

type BaseTimeEntity struct {
	CreatedDate time.Time `gorm:"autoCreateTime;column:created_date" json:"createdDate"`
	UpdatedDate time.Time `gorm:"autoUpdateTime;column:updated_date" json:"updatedDate"`
}
