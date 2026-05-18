package auth

import (
	"puppynote/internal/model"
	pnemail "puppynote/pkg/email"
	pnerrors "puppynote/pkg/errors"
	pnjwt "puppynote/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) NormalLogin(req LoginRequest) (*LoginResponse, error) {
	u, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if u.SnsType != model.SnsTypeNormal {
		return nil, pnerrors.New("SNS 로그인 계정입니다.")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, pnerrors.ErrInvalidPassword
	}

	return s.issueTokens(u, req.DeviceID, req.PushKey)
}

func (s *Service) OAuthLogin(req LoginOauthRequest) (*LoginResponse, error) {
	email, err := getOAuthEmail(req.SnsType, req.Token)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		newUser := &model.User{
			Email:   email,
			SnsType: req.SnsType,
			Role:    model.RoleUser,
			UseYn:   "Y",
		}
		if err = s.repo.SaveUser(newUser); err != nil {
			return nil, err
		}
		u = newUser
	}

	if u.SnsType != req.SnsType {
		return nil, pnerrors.New("다른 방식으로 가입된 계정입니다.")
	}

	return s.issueTokens(u, req.DeviceID, req.PushKey)
}

func (s *Service) Refresh(req TokenRefreshRequest) (*TokenRefreshResponse, error) {
	stored, err := s.repo.FindRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	info, err := pnjwt.Parse(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	pair, err := pnjwt.GeneratePair(*info)
	if err != nil {
		return nil, err
	}

	if err = s.repo.UpdateRefreshToken(stored.ID, pair.RefreshToken); err != nil {
		return nil, err
	}

	return &TokenRefreshResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}, nil
}

func (s *Service) SendPasswordResetEmail(req PasswordEmailSendRequest) error {
	u, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if u.SnsType != model.SnsTypeNormal {
		return pnerrors.New("SNS 로그인 계정은 비밀번호를 변경할 수 없습니다.")
	}
	_, err = pnemail.SendVerificationCode(req.Email)
	return err
}

func (s *Service) ResetPassword(req PasswordResetRequest) error {
	u, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if u.SnsType != model.SnsTypeNormal {
		return pnerrors.New("SNS 로그인 계정은 비밀번호를 변경할 수 없습니다.")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.UpdatePassword(u.ID, string(hashed))
}

func (s *Service) issueTokens(u *model.User, deviceID, pushKey string) (*LoginResponse, error) {
	pair, err := pnjwt.GeneratePair(pnjwt.LoginUserInfo{UserID: u.ID, Email: u.Email})
	if err != nil {
		return nil, err
	}

	if err = s.repo.UpsertRefreshToken(u.ID, deviceID, pair.RefreshToken); err != nil {
		return nil, err
	}
	if pushKey != "" {
		_ = s.repo.UpsertPushToken(u.ID, deviceID, pushKey)
	}

	return &LoginResponse{
		Email:        u.Email,
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
	}, nil
}
