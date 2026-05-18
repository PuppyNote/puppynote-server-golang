package pet

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

// Pet Services
func (s *Service) GetMyPets(userID int64) ([]PetResponse, error) {
	members, err := s.repo.FindFamilyMembersByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]PetResponse, 0, len(members))
	for _, fm := range members {
		responses = append(responses, PetResponse{
			PetID:         fm.Pet.ID,
			PetName:       fm.Pet.Name,
			PetProfileUrl: fm.Pet.ProfileImage,
			RoleType:      fm.Role,
		})
	}
	return responses, nil
}

func (s *Service) CreatePet(userID int64, req PetCreateRequest) (*PetCreateResponse, error) {
	p := &model.Pet{
		Name:               req.Name,
		BirthDate:          req.BirthDate,
		ProfileImage:       req.ProfileImage,
		RegistrationNumber: req.RegistrationNumber,
	}
	if err := s.repo.SavePet(p); err != nil {
		return nil, err
	}

	owner := &model.FamilyMember{
		UserID: userID,
		PetID:  p.ID,
		Role:   model.FamilyRoleOwner,
		Status: model.FamilyStatusDone,
	}
	if err := s.repo.SaveFamilyMember(owner); err != nil {
		return nil, err
	}

	return &PetCreateResponse{PetID: p.ID, PetName: p.Name}, nil
}

func (s *Service) UpdatePet(userID, petID int64, req PetUpdateRequest) error {
	p, err := s.repo.FindPetByID(petID)
	if err != nil {
		return err
	}

	p.Name = req.Name
	p.BirthDate = req.BirthDate
	p.ProfileImage = req.ProfileImage
	p.RegistrationNumber = req.RegistrationNumber

	return s.repo.UpdatePet(p)
}

func (s *Service) DeletePet(userID, petID int64) error {
	fm, err := s.repo.FindFamilyMember(userID, petID)
	if err != nil {
		return err
	}
	if fm.Role != model.FamilyRoleOwner {
		return pnerrors.New("반려동물 삭제 권한이 없습니다.")
	}

	if err = s.repo.DeleteAllFamilyMembersByPetID(petID); err != nil {
		return err
	}
	return s.repo.DeletePet(petID)
}

// FamilyMember Services
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

	newMember := &model.FamilyMember{
		UserID: req.InviteeUserID,
		PetID:  req.PetID,
		Role:   model.FamilyRoleFamily,
		Status: model.FamilyStatusPending,
	}
	return s.repo.SaveFamilyMember(newMember)
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

	// FAMILY 역할인 경우 자기 자신만 탈퇴 가능
	if currentUserID != targetUserID {
		return pnerrors.New("삭제 권한이 없습니다.")
	}
	return s.repo.DeleteFamilyMember(currentUserID, petID)
}
