package model

type AlertType            string
type AlertHistoryStatus   string
type AlertDestinationType string

const (
	AlertTypeOn  AlertType = "ON"
	AlertTypeOff AlertType = "OFF"

	AlertHistoryStatusChecked   AlertHistoryStatus = "CHECKED"
	AlertHistoryStatusUnchecked AlertHistoryStatus = "UNCHECKED"

	AlertDestDailyReport AlertDestinationType = "DAILY_REPORT"
	AlertDestFriend      AlertDestinationType = "FRIEND"
	AlertDestFriendCode  AlertDestinationType = "FRIEND_CODE"
	AlertDestWalk        AlertDestinationType = "WALK"
	AlertDestPetItem     AlertDestinationType = "PET_ITEM"
	AlertDestFamilyInvite AlertDestinationType = "FAMILY_INVITE"
)

type AlertSetting struct {
	ID     int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID int64     `gorm:"column:user_id;uniqueIndex" json:"userId"`
	User   User      `gorm:"foreignKey:UserID" json:"-"`
	All    AlertType `gorm:"column:all_alerts;type:varchar(10)" json:"all"`
	Walk   AlertType `gorm:"column:walk;type:varchar(10)" json:"walk"`
	Friend AlertType `gorm:"column:friend;type:varchar(10)" json:"friend"`
	BaseTimeEntity
}

func (AlertSetting) TableName() string { return "alert_settings" }

type AlertHistory struct {
	ID                   int64                `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID               int64                `gorm:"column:user_id;not null" json:"userId"`
	User                 User                 `gorm:"foreignKey:UserID" json:"-"`
	AlertDescription     string               `gorm:"column:alert_description" json:"alertDescription"`
	AlertHistoryStatus   AlertHistoryStatus   `gorm:"column:alert_history_status;type:varchar(20)" json:"alertHistoryStatus"`
	AlertDestinationType AlertDestinationType `gorm:"column:alert_destination_type;type:varchar(50)" json:"alertDestinationType"`
	AlertDestinationInfo string               `gorm:"column:alert_destination_info" json:"alertDestinationInfo"`
	BaseTimeEntity
}

func (AlertHistory) TableName() string { return "alert_histories" }
