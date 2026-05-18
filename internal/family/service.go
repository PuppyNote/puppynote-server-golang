package family

import (
	"puppynote/internal/model"
	pnerrors "puppynote/pkg/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetFamilyMembers(currentUserID, petID int64) ([]FamilyMemberResponse, error) {
	members, err := s.repo.FindFamilyMembersByPetID(petID)
	if err != nil {
		return nil, err
	}

	responses := make([]FamilyMemberResponse, 0, len(members))
	for _, fm := range members {
		if fm.UserID == currentUserID {
			continue
		}
		responses = append(responses, FamilyMemberResponse{
			UserID:     fm.User.ID,
			NickName:   fm.User.NickName,
			ProfileUrl: fm.User.ProfileUrl,
			Role:       fm.Role,
			Status:     fm.Status,
		})
	}
	return responses, nil
}

func (s *Service) SearchUsersByEmail(currentUserID int64, email string) ([]UserSearchResponse, error) {
	users, err := s.repo.SearchUsersByEmail(email)
	if err != nil {
		return nil, err
	}

	responses := make([]UserSearchResponse, 0, len(users))
	for _, u := range users {
		if u.ID == currentUserID {
			continue
		}
		responses = append(responses, UserSearchResponse{
			UserID:     u.ID,
			Email:      u.Email,
			NickName:   u.NickName,
			ProfileUrl: u.ProfileUrl,
		})
	}
	return responses, nil
}

func (s *Service) InviteFamilyMember(inviterID int64, req FamilyMemberInviteRequest) error {
	fm, err := s.repo.FindFamilyMember(inviterID, req.PetID)
	if err != nil {
		return err
	}
	if fm.Role != model.FamilyRoleOwner {
		return pnerrors.New("가족 초대 권한이 없습니다.")
	}

	if s.repo.ExistsFamilyMember(req.InviteeUserID, req.PetID) {
		return pnerrors.New("이미 등록되었거나 초대 대기 중인 구성원입니다.")
	}

	return s.repo.SaveFamilyMember(&model.FamilyMember{
		UserID: req.InviteeUserID,
		PetID:  req.PetID,
		Role:   model.FamilyRoleFamily,
		Status: model.FamilyStatusPending,
	})
}

func (s *Service) RegisterFamilyMember(req FamilyMemberRegisterRequest) error {
	fm, err := s.repo.FindFamilyMember(req.UserID, req.PetID)
	if err != nil {
		return err
	}
	if fm.Status == model.FamilyStatusDone {
		return pnerrors.New("이미 등록된 가족 구성원입니다.")
	}
	return s.repo.UpdateFamilyMemberStatus(req.UserID, req.PetID, model.FamilyStatusDone)
}

func (s *Service) DeleteFamilyMember(currentUserID, targetUserID, petID int64) error {
	currentFM, err := s.repo.FindFamilyMember(currentUserID, petID)
	if err != nil {
		return err
	}

	if currentFM.Role == model.FamilyRoleOwner {
		targetFM, err := s.repo.FindFamilyMember(targetUserID, petID)
		if err != nil {
			return err
		}
		if targetFM.Role == model.FamilyRoleOwner {
			return pnerrors.New("주인은 삭제할 수 없습니다.")
		}
		return s.repo.DeleteFamilyMember(targetUserID, petID)
	}

	if currentUserID != targetUserID {
		return pnerrors.New("삭제 권한이 없습니다.")
	}
	return s.repo.DeleteFamilyMember(currentUserID, petID)
}
