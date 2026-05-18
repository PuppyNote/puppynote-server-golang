package model

type PetTip struct {
	ID      int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Content string `gorm:"column:content;not null;size:500" json:"content"`
	BaseTimeEntity
}

func (PetTip) TableName() string { return "pet_tips" }
