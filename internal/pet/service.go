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
