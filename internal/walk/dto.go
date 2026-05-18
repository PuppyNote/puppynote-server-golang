package walk

import "time"

type WalkCreateRequest struct {
	PetID     int64      `json:"petId" binding:"required"`
	StartTime time.Time  `json:"startTime" binding:"required"`
	EndTime   time.Time  `json:"endTime" binding:"required"`
	Latitude  *float64   `json:"latitude"`
	Longitude *float64   `json:"longitude"`
	Location  string     `json:"location"`
	Memo      string     `json:"memo"`
	PhotoKeys []string   `json:"photoKeys"`
}

type WalkResponse struct {
	WalkID    int64      `json:"walkId"`
	PetID     int64      `json:"petId"`
	StartTime time.Time  `json:"startTime"`
	EndTime   time.Time  `json:"endTime"`
	Latitude  *float64   `json:"latitude"`
	Longitude *float64   `json:"longitude"`
	Location  string     `json:"location"`
	Memo      string     `json:"memo"`
	PhotoUrls []string   `json:"photoUrls,omitempty"`
}

type WalkCalendarResponse struct {
	Date    string `json:"date"`
	HasWalk bool   `json:"hasWalk"`
}
