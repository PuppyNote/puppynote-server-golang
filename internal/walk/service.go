package walk

import (
	"fmt"
	"time"

	"puppynote/internal/model"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req WalkCreateRequest) (*WalkResponse, error) {
	photos := make([]model.WalkPhoto, 0, len(req.PhotoKeys))
	for _, key := range req.PhotoKeys {
		photos = append(photos, model.WalkPhoto{ImageKey: key})
	}

	w := &model.Walk{
		PetID:     req.PetID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
		Memo:      req.Memo,
		Photos:    photos,
	}
	if err := s.repo.Save(w); err != nil {
		return nil, err
	}

	return toWalkResponse(w), nil
}

func (s *Service) GetDetail(walkID int64) (*WalkResponse, error) {
	w, err := s.repo.FindByID(walkID)
	if err != nil {
		return nil, err
	}
	return toWalkResponse(w), nil
}

func (s *Service) GetByDate(petID int64, date time.Time) ([]WalkResponse, error) {
	walks, err := s.repo.FindByPetIDAndDate(petID, date)
	if err != nil {
		return nil, err
	}
	responses := make([]WalkResponse, 0, len(walks))
	for _, w := range walks {
		responses = append(responses, *toWalkResponse(&w))
	}
	return responses, nil
}

func (s *Service) GetCalendar(petID int64, year, month int) ([]WalkCalendarResponse, error) {
	walks, err := s.repo.FindByPetIDAndMonth(petID, year, month)
	if err != nil {
		return nil, err
	}

	walkDates := make(map[string]bool)
	for _, w := range walks {
		dateStr := w.StartTime.Format("2006-01-02")
		walkDates[dateStr] = true
	}

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	daysInMonth := start.AddDate(0, 1, -1).Day()

	responses := make([]WalkCalendarResponse, 0, daysInMonth)
	for day := 1; day <= daysInMonth; day++ {
		dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		responses = append(responses, WalkCalendarResponse{
			Date:    dateStr,
			HasWalk: walkDates[dateStr],
		})
	}
	return responses, nil
}

func (s *Service) Delete(walkID int64) error {
	return s.repo.Delete(walkID)
}

func toWalkResponse(w *model.Walk) *WalkResponse {
	photoUrls := make([]string, 0, len(w.Photos))
	for _, p := range w.Photos {
		photoUrls = append(photoUrls, p.ImageKey)
	}
	return &WalkResponse{
		WalkID:    w.ID,
		PetID:     w.PetID,
		StartTime: w.StartTime,
		EndTime:   w.EndTime,
		Latitude:  w.Latitude,
		Longitude: w.Longitude,
		Location:  w.Location,
		Memo:      w.Memo,
		PhotoUrls: photoUrls,
	}
}
