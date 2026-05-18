package pet

import (
	"time"

	"puppynote/internal/model"
)

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
	PetID         int64                  `json:"petId"`
	PetName       string                 `json:"petName"`
	PetProfileUrl string                 `json:"petProfileUrl"`
	RoleType      model.FamilyMemberRole `json:"roleType"`
}
