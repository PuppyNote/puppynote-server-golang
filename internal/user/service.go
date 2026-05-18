package user

import (
	"puppynote/internal/model"
	pnemail "puppynote/pkg/email"
	pnerrors "puppynote/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendVerificationEmail(email string) error {
	if s.repo.ExistsByEmail(email) {
		return pnerrors.ErrDuplicateEmail
	}
	_, err := pnemail.SendVerificationCode(email)
	return err
}

func (s *Service) SignUp(req SignUpRequest) (*SignUpResponse, error) {
	if s.repo.ExistsByEmail(req.Email) {
		return nil, pnerrors.ErrDuplicateEmail
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &model.User{
		Email:    req.Email,
		Password: string(hashed),
		NickName: req.NickName,
		SnsType:  model.SnsTypeNormal,
		Role:     model.RoleUser,
		UseYn:    "Y",
	}
	if err = s.repo.Save(u); err != nil {
		return nil, err
	}

	return &SignUpResponse{Email: u.Email, NickName: u.NickName}, nil
}

func (s *Service) GetProfile(userID int64) (*UserProfileResponse, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return &UserProfileResponse{
		UserID:     u.ID,
		Email:      u.Email,
		NickName:   u.NickName,
		ProfileUrl: u.ProfileUrl,
	}, nil
}

func (s *Service) UpdateProfile(userID int64, req UserProfileUpdateRequest) error {
	if req.NickName != "" {
		if err := s.repo.UpdateNickName(userID, req.NickName); err != nil {
			return err
		}
	}
	if req.ProfileUrl != "" {
		if err := s.repo.UpdateProfileUrl(userID, req.ProfileUrl); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) Withdraw(userID int64) error {
	return s.repo.Withdraw(userID)
}
