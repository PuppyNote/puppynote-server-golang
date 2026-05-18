package model

type SafetyLevel string

const (
	SafetyLevelGood   SafetyLevel = "GOOD"
	SafetyLevelNotion SafetyLevel = "NOTION"
	SafetyLevelBad    SafetyLevel = "BAD"
)

type FoodChatHistory struct {
	ID          int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Question    string      `gorm:"column:question;type:text;not null" json:"question"`
	Answer      string      `gorm:"column:answer;type:text;not null" json:"answer"`
	SafetyLevel SafetyLevel `gorm:"column:safety_level;type:varchar(10)" json:"safetyLevel"`
	BaseTimeEntity
}

func (FoodChatHistory) TableName() string { return "food_chat_histories" }
