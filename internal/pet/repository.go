package pet

import (
	"puppynote/internal/model"
	pnerrors "puppynote/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SavePet(p *model.Pet) error {
	return r.db.Create(p).Error
}

func (r *Repository) FindPetByID(id int64) (*model.Pet, error) {
	var p model.Pet
	if err := r.db.First(&p, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.ErrPetNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *Repository) UpdatePet(p *model.Pet) error {
	return r.db.Save(p).Error
}

func (r *Repository) DeletePet(id int64) error {
	return r.db.Delete(&model.Pet{}, id).Error
}

func (r *Repository) FindFamilyMembersByUserID(userID int64) ([]model.FamilyMember, error) {
	var members []model.FamilyMember
	err := r.db.Preload("Pet").Preload("User").
		Where("user_id = ? AND status = 'DONE'", userID).
		Find(&members).Error
	return members, err
}

func (r *Repository) FindFamilyMember(userID, petID int64) (*model.FamilyMember, error) {
	var fm model.FamilyMember
	if err := r.db.Where("user_id = ? AND pet_id = ?", userID, petID).First(&fm).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.New("가족 구성원 정보를 찾을 수 없습니다.")
		}
		return nil, err
	}
	return &fm, nil
}

func (r *Repository) SaveFamilyMember(fm *model.FamilyMember) error {
	return r.db.Create(fm).Error
}

func (r *Repository) DeleteAllFamilyMembersByPetID(petID int64) error {
	return r.db.Where("pet_id = ?", petID).Delete(&model.FamilyMember{}).Error
}
