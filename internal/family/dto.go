package family

import "puppynote/internal/model"

type FamilyMemberResponse struct {
	UserID     int64                    `json:"userId"`
	NickName   string                   `json:"nickName"`
	ProfileUrl string                   `json:"profileUrl"`
	Role       model.FamilyMemberRole   `json:"role"`
	Status     model.FamilyMemberStatus `json:"status"`
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
