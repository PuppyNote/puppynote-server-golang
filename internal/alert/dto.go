package alert

import (
	"time"

	"github.com/PuppyNote/puppynote-server-golang/internal/model"
)

type AlertSettingResponse struct {
	All    model.AlertType `json:"all"`
	Walk   model.AlertType `json:"walk"`
	Friend model.AlertType `json:"friend"`
}

type AlertSettingUpdateRequest struct {
	All    model.AlertType `json:"all" binding:"required"`
	Walk   model.AlertType `json:"walk" binding:"required"`
	Friend model.AlertType `json:"friend" binding:"required"`
}

type AlertHistoryResponse struct {
	ID                   int64                        `json:"id"`
	AlertDescription     string                       `json:"alertDescription"`
	AlertHistoryStatus   model.AlertHistoryStatus     `json:"alertHistoryStatus"`
	AlertDestinationType model.AlertDestinationType   `json:"alertDestinationType"`
	AlertDestinationInfo string                       `json:"alertDestinationInfo"`
	CreatedDate          time.Time                    `json:"createdDate"`
}

type AlertHistoryListResponse struct {
	Content     []AlertHistoryResponse `json:"content"`
	Page        int                    `json:"page"`
	Size        int                    `json:"size"`
	TotalElements int64                `json:"totalElements"`
	TotalPages  int                    `json:"totalPages"`
}
