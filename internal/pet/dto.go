package pet

import (
	"time"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
)

// Pet DTOs
type PetCreateRequest struct {
	Name               string     `json:"name" binding:"required,max=50"`
	BirthDate          *time.Time `json:"birthDate"`
	ProfileImage       string     `json:"profileImage"`
	RegistrationNumber string     `json:"registrationNumber"`
}

type PetCreateResponse struct {
	PetID   int64  `json:"petId"`
	PetName string `json:"petName"`
}

type PetUpdateRequest struct {
	Name               string     `json:"name" binding:"required"`
	BirthDate          *time.Time `json:"birthDate"`
	ProfileImage       string     `json:"profileImage"`
	RegistrationNumber string     `json:"registrationNumber"`
}

type PetResponse struct {
	PetID         int64                    `json:"petId"`
	PetName       string                   `json:"petName"`
	PetProfileUrl string                   `json:"petProfileUrl"`
	RoleType      model.FamilyMemberRole   `json:"roleType"`
}

// FamilyMember DTOs
type FamilyMemberResponse struct {
	UserID     int64                      `json:"userId"`
	NickName   string                     `json:"nickName"`
	ProfileUrl string                     `json:"profileUrl"`
	Role       model.FamilyMemberRole     `json:"role"`
	Status     model.FamilyMemberStatus   `json:"status"`
}

type UserSearchResponse struct {
	UserID     int64  `json:"userId"`
	Email      string `json:"email"`
	NickName   string `json:"nickName"`
	ProfileUrl string `json:"profileUrl"`
}

type FamilyMemberInviteRequest struct {
	InviteeUserID int64 `json:"inviteeUserId" binding:"required"`
	PetID         int64 `json:"petId" binding:"required"`
}

type FamilyMemberRegisterRequest struct {
	UserID int64 `json:"userId" binding:"required"`
	PetID  int64 `json:"petId" binding:"required"`
}
