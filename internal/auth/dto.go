package auth

import "github.com/PuppyNote/puppynote-server-golang/internal/model"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	DeviceID string `json:"deviceId" binding:"required"`
	PushKey  string `json:"pushKey"`
}

type LoginOauthRequest struct {
	Token    string        `json:"token" binding:"required"`
	SnsType  model.SnsType `json:"snsType" binding:"required"`
	DeviceID string        `json:"deviceId" binding:"required"`
	PushKey  string        `json:"pushKey"`
}

type LoginResponse struct {
	Email        string `json:"email"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenRefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type TokenRefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type PasswordEmailSendRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetRequest struct {
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}
