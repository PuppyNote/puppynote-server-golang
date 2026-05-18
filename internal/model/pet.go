package model

import "time"

type FamilyMemberRole   string
type FamilyMemberStatus string

const (
	FamilyRoleOwner  FamilyMemberRole = "OWNER"
	FamilyRoleFamily FamilyMemberRole = "FAMILY"

	FamilyStatusPending FamilyMemberStatus = "PENDING"
	FamilyStatusDone    FamilyMemberStatus = "DONE"
)

type Pet struct {
	ID                 int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name               string     `gorm:"column:name;not null;size:50" json:"name"`
	BirthDate          *time.Time `gorm:"column:birth_date" json:"birthDate"`
	ProfileImage       string     `gorm:"column:profile_image;size:255" json:"profileImage"`
	RegistrationNumber string     `gorm:"column:registration_number;size:50" json:"registrationNumber"`
	BaseTimeEntity
}

func (Pet) TableName() string { return "pets" }

type FamilyMember struct {
	UserID int64              `gorm:"primaryKey;column:user_id" json:"userId"`
	PetID  int64              `gorm:"primaryKey;column:pet_id" json:"petId"`
	User   User               `gorm:"foreignKey:UserID" json:"-"`
	Pet    Pet                `gorm:"foreignKey:PetID" json:"-"`
	Role   FamilyMemberRole   `gorm:"column:role;type:varchar(50)" json:"role"`
	Status FamilyMemberStatus `gorm:"column:status;type:varchar(20);not null" json:"status"`
}

func (FamilyMember) TableName() string { return "family_members" }

type WalkPhoto struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	WalkID   int64  `gorm:"column:walk_id;not null" json:"walkId"`
	Walk     Walk   `gorm:"foreignKey:WalkID" json:"-"`
	ImageKey string `gorm:"column:image_key;not null;size:255" json:"imageKey"`
	BaseTimeEntity
}

func (WalkPhoto) TableName() string { return "walk_photos" }
