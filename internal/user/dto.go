package user

type EmailSendRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	NickName string `json:"nickName" binding:"required,max=20"`
}

type SignUpResponse struct {
	Email    string `json:"email"`
	NickName string `json:"nickName"`
}

type UserProfileResponse struct {
	UserID     int64  `json:"userId"`
	Email      string `json:"email"`
	NickName   string `json:"nickName"`
	ProfileUrl string `json:"profileUrl"`
}

type UserProfileUpdateRequest struct {
	NickName   string `json:"nickName" binding:"omitempty,max=20"`
	ProfileUrl string `json:"profileUrl"`
}
