package walkalarm

import (
	"time"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
	pnerrors "github.com/PuppyNote/puppynote-server-golang/pkg/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req WalkAlarmCreateRequest) (*WalkAlarmResponse, error) {
	alarmTime, err := parseTime(req.AlarmTime)
	if err != nil {
		return nil, pnerrors.New("알람 시간 형식은 HH:mm입니다.")
	}

	alarm := &model.PetWalkAlarm{
		PetID:       req.PetID,
		AlarmStatus: req.AlarmStatus,
		AlarmTime:   model.LocalTime{Time: alarmTime},
	}
	if err = s.repo.Save(alarm); err != nil {
		return nil, err
	}

	for _, day := range req.AlarmDays {
		if err = s.repo.db.Create(&model.PetAlarmDay{PetAlarmID: alarm.ID, AlarmDay: day}).Error; err != nil {
			return nil, err
		}
	}

	alarm.AlarmDays = toPetAlarmDays(alarm.ID, req.AlarmDays)
	return toResponse(alarm), nil
}

func (s *Service) GetByPetID(petID int64) ([]WalkAlarmResponse, error) {
	alarms, err := s.repo.FindByPetID(petID)
	if err != nil {
		return nil, err
	}
	responses := make([]WalkAlarmResponse, 0, len(alarms))
	for _, a := range alarms {
		responses = append(responses, *toResponse(&a))
	}
	return responses, nil
}

func (s *Service) Update(req WalkAlarmUpdateRequest) (*WalkAlarmResponse, error) {
	alarm, err := s.repo.FindByID(req.AlarmID)
	if err != nil {
		return nil, err
	}

	alarmTime, err := parseTime(req.AlarmTime)
	if err != nil {
		return nil, pnerrors.New("알람 시간 형식은 HH:mm입니다.")
	}

	alarm.AlarmStatus = req.AlarmStatus
	alarm.AlarmTime = model.LocalTime{Time: alarmTime}
	days := toPetAlarmDays(alarm.ID, req.AlarmDays)

	if err = s.repo.Update(alarm, req.AlarmDays); err != nil {
		return nil, err
	}
	alarm.AlarmDays = days
	return toResponse(alarm), nil
}

func (s *Service) UpdateStatus(req WalkAlarmStatusUpdateRequest) (*WalkAlarmResponse, error) {
	if err := s.repo.UpdateStatus(req.AlarmID, req.AlarmStatus); err != nil {
		return nil, err
	}
	alarm, err := s.repo.FindByID(req.AlarmID)
	if err != nil {
		return nil, err
	}
	return toResponse(alarm), nil
}

func (s *Service) Delete(alarmID int64) error {
	return s.repo.Delete(alarmID)
}

func parseTime(hhmm string) (time.Time, error) {
	return time.Parse("15:04", hhmm)
}

func toPetAlarmDays(alarmID int64, days []model.AlarmDay) []model.PetAlarmDay {
	result := make([]model.PetAlarmDay, 0, len(days))
	for _, d := range days {
		result = append(result, model.PetAlarmDay{PetAlarmID: alarmID, AlarmDay: d})
	}
	return result
}

func toResponse(alarm *model.PetWalkAlarm) *WalkAlarmResponse {
	days := make([]model.AlarmDay, 0, len(alarm.AlarmDays))
	for _, d := range alarm.AlarmDays {
		days = append(days, d.AlarmDay)
	}
	return &WalkAlarmResponse{
		AlarmID:     alarm.ID,
		AlarmStatus: alarm.AlarmStatus,
		AlarmDays:   days,
		AlarmTime:   alarm.AlarmTime.Format("15:04"),
	}
}
