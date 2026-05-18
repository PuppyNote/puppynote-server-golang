package model

type SnsType string
type Role string

const (
	SnsTypeNormal SnsType = "NORMAL"
	SnsTypeKakao  SnsType = "KAKAO"
	SnsTypeGoogle SnsType = "GOOGLE"
	SnsTypeApple  SnsType = "APPLE"
)

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID         int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email      string  `gorm:"column:email" json:"email"`
	Password   string  `gorm:"column:password" json:"-"`
	NickName   string  `gorm:"column:nick_name" json:"nickName"`
	ProfileUrl string  `gorm:"column:profile_url" json:"profileUrl"`
	SnsType    SnsType `gorm:"column:sns_type;type:varchar(20)" json:"snsType"`
	Role       Role    `gorm:"column:role;type:varchar(20)" json:"role"`
	UseYn      string  `gorm:"column:use_yn;default:Y" json:"useYn"`
	BaseTimeEntity
}

func (User) TableName() string { return "users" }

type RefreshToken struct {
	ID           int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID     string `gorm:"column:device_id" json:"deviceId"`
	UserID       int64  `gorm:"column:user_id" json:"userId"`
	User         User   `gorm:"foreignKey:UserID" json:"-"`
	RefreshToken string `gorm:"column:refresh_token" json:"refreshToken"`
	BaseTimeEntity
}

func (RefreshToken) TableName() string { return "refresh_tokens" }

type Push struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceID  string `gorm:"column:device_id;uniqueIndex" json:"deviceId"`
	PushToken string `gorm:"column:push_token" json:"pushToken"`
	UserID    int64  `gorm:"column:user_id" json:"userId"`
	User      User   `gorm:"foreignKey:UserID" json:"-"`
	BaseTimeEntity
}

func (Push) TableName() string { return "pushes" }
