package walk

import (
	"time"

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

func (r *Repository) Save(w *model.Walk) error {
	return r.db.Create(w).Error
}

func (r *Repository) FindByID(id int64) (*model.Walk, error) {
	var w model.Walk
	if err := r.db.Preload("Photos").First(&w, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.New("산책 기록을 찾을 수 없습니다.")
		}
		return nil, err
	}
	return &w, nil
}

func (r *Repository) FindByPetIDAndDate(petID int64, date time.Time) ([]model.Walk, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24 * time.Hour)

	var walks []model.Walk
	err := r.db.Preload("Photos").
		Where("pet_id = ? AND start_time >= ? AND start_time < ?", petID, start, end).
		Order("end_time DESC").
		Find(&walks).Error
	return walks, err
}

func (r *Repository) FindByPetIDAndMonth(petID int64, year, month int) ([]model.Walk, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)

	var walks []model.Walk
	err := r.db.Select("id, pet_id, start_time, end_time").
		Where("pet_id = ? AND start_time >= ? AND start_time < ?", petID, start, end).
		Find(&walks).Error
	return walks, err
}

func (r *Repository) Delete(id int64) error {
	return r.db.Select("Photos").Delete(&model.Walk{}, id).Error
}

func (r *Repository) DeleteAllByPetID(petID int64) error {
	return r.db.Where("pet_id = ?", petID).Delete(&model.Walk{}).Error
}
