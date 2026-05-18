package user

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("email = ? AND use_yn = 'Y'", email).Count(&count)
	return count > 0
}

func (r *Repository) Save(u *model.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) FindByID(id int64) (*model.User, error) {
	var u model.User
	if err := r.db.Where("id = ? AND use_yn = 'Y'", id).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("email = ? AND use_yn = 'Y'", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) UpdateNickName(id int64, nickName string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("nick_name", nickName).Error
}

func (r *Repository) UpdateProfileUrl(id int64, profileUrl string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("profile_url", profileUrl).Error
}

func (r *Repository) Withdraw(id int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("use_yn", "N").Error
}
