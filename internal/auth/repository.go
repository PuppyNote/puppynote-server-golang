package auth

import (
	"puppynote/internal/model"
	pnerrors "puppynote/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByEmail(email string) (*model.User, error) {
	var u model.User
	if err := r.db.Where("email = ? AND use_yn = 'Y'", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) SaveUser(u *model.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) UpsertRefreshToken(userID int64, deviceID, token string) error {
	rt := model.RefreshToken{
		UserID:       userID,
		DeviceID:     deviceID,
		RefreshToken: token,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"refresh_token"}),
	}).Create(&rt).Error
}

func (r *Repository) FindRefreshToken(token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	if err := r.db.Where("refresh_token = ?", token).First(&rt).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.NewWithCode(401, "유효하지 않은 RefreshToken입니다.")
		}
		return nil, err
	}
	return &rt, nil
}

func (r *Repository) UpdateRefreshToken(id int64, token string) error {
	return r.db.Model(&model.RefreshToken{}).Where("id = ?", id).Update("refresh_token", token).Error
}

func (r *Repository) UpsertPushToken(userID int64, deviceID, pushToken string) error {
	p := model.Push{
		UserID:    userID,
		DeviceID:  deviceID,
		PushToken: pushToken,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"push_token", "user_id"}),
	}).Create(&p).Error
}

func (r *Repository) UpdatePassword(userID int64, hashedPwd string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("password", hashedPwd).Error
}
