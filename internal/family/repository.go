package family

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

func (r *Repository) FindFamilyMembersByPetID(petID int64) ([]model.FamilyMember, error) {
	var members []model.FamilyMember
	err := r.db.Preload("User").
		Where("pet_id = ? AND status = 'DONE'", petID).
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

func (r *Repository) UpdateFamilyMemberStatus(userID, petID int64, status model.FamilyMemberStatus) error {
	return r.db.Model(&model.FamilyMember{}).
		Where("user_id = ? AND pet_id = ?", userID, petID).
		Update("status", status).Error
}

func (r *Repository) DeleteFamilyMember(userID, petID int64) error {
	return r.db.Where("user_id = ? AND pet_id = ?", userID, petID).
		Delete(&model.FamilyMember{}).Error
}

func (r *Repository) ExistsFamilyMember(userID, petID int64) bool {
	var count int64
	r.db.Model(&model.FamilyMember{}).Where("user_id = ? AND pet_id = ?", userID, petID).Count(&count)
	return count > 0
}

func (r *Repository) SearchUsersByEmail(email string) ([]model.User, error) {
	var users []model.User
	err := r.db.Where("email LIKE ? AND use_yn = 'Y'", "%"+email+"%").Find(&users).Error
	return users, err
}
