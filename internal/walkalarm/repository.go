package walkalarm

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

func (r *Repository) Save(alarm *model.PetWalkAlarm) error {
	return r.db.Create(alarm).Error
}

func (r *Repository) FindByID(id int64) (*model.PetWalkAlarm, error) {
	var alarm model.PetWalkAlarm
	if err := r.db.Preload("AlarmDays").First(&alarm, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, pnerrors.New("산책 알람을 찾을 수 없습니다.")
		}
		return nil, err
	}
	return &alarm, nil
}

func (r *Repository) FindByPetID(petID int64) ([]model.PetWalkAlarm, error) {
	var alarms []model.PetWalkAlarm
	err := r.db.Preload("AlarmDays").Where("pet_id = ?", petID).Find(&alarms).Error
	return alarms, err
}

func (r *Repository) Update(alarm *model.PetWalkAlarm, days []model.AlarmDay) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("pet_alarm_id = ?", alarm.ID).Delete(&model.PetAlarmDay{}).Error; err != nil {
			return err
		}
		for _, day := range days {
			if err := tx.Create(&model.PetAlarmDay{PetAlarmID: alarm.ID, AlarmDay: day}).Error; err != nil {
				return err
			}
		}
		return tx.Save(alarm).Error
	})
}

func (r *Repository) UpdateStatus(id int64, status model.AlarmStatus) error {
	return r.db.Model(&model.PetWalkAlarm{}).Where("id = ?", id).Update("alarm_status", status).Error
}

func (r *Repository) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("pet_alarm_id = ?", id).Delete(&model.PetAlarmDay{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.PetWalkAlarm{}, id).Error
	})
}

func (r *Repository) DeleteAllByPetID(petID int64) error {
	return r.db.Where("pet_id = ?", petID).Delete(&model.PetWalkAlarm{}).Error
}

func (r *Repository) FindActiveAlarmsForBatch(status model.AlarmStatus, alarmTime string, day model.AlarmDay) ([]model.PetWalkAlarm, error) {
	var alarms []model.PetWalkAlarm
	err := r.db.Preload("Pet").Preload("AlarmDays").
		Joins("JOIN pet_alarm_days pad ON pad.pet_alarm_id = pet_walk_alarms.id").
		Where("pet_walk_alarms.alarm_status = ? AND pet_walk_alarms.alarm_time = ? AND pad.alarm_day = ?",
			status, alarmTime, day).
		Distinct("pet_walk_alarms.id").
		Find(&alarms).Error
	return alarms, err
}
