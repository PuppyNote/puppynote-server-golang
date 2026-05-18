package alert

import (
	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindSettingByUserID(userID int64) (*model.AlertSetting, error) {
	var s model.AlertSetting
	err := r.db.Where("user_id = ?", userID).First(&s).Error
	if err == gorm.ErrRecordNotFound {
		s = model.AlertSetting{
			UserID: userID,
			All:    model.AlertTypeOn,
			Walk:   model.AlertTypeOn,
			Friend: model.AlertTypeOn,
		}
		r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&s)
		return &s, nil
	}
	return &s, err
}

func (r *Repository) UpdateSetting(userID int64, req AlertSettingUpdateRequest) (*model.AlertSetting, error) {
	s, err := r.FindSettingByUserID(userID)
	if err != nil {
		return nil, err
	}
	s.All = req.All
	s.Walk = req.Walk
	s.Friend = req.Friend
	err = r.db.Save(s).Error
	return s, err
}

func (r *Repository) FindHistoriesByUserID(userID int64, page, size int) ([]model.AlertHistory, int64, error) {
	var histories []model.AlertHistory
	var total int64

	r.db.Model(&model.AlertHistory{}).Where("user_id = ?", userID).Count(&total)
	err := r.db.Where("user_id = ?", userID).
		Order("created_date DESC").
		Offset(page * size).Limit(size).
		Find(&histories).Error
	return histories, total, err
}

func (r *Repository) HasUnchecked(userID int64) bool {
	var count int64
	r.db.Model(&model.AlertHistory{}).
		Where("user_id = ? AND alert_history_status = 'UNCHECKED'", userID).Count(&count)
	return count > 0
}

func (r *Repository) UpdateHistoryStatus(id int64) error {
	return r.db.Model(&model.AlertHistory{}).Where("id = ?", id).
		Update("alert_history_status", model.AlertHistoryStatusChecked).Error
}

func (r *Repository) SaveHistory(h *model.AlertHistory) error {
	return r.db.Create(h).Error
}
