package home

import (
	"puppynote/internal/model"
	"puppynote/pkg/middleware"
	"puppynote/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HomeResponse struct {
	User                UserSummary  `json:"user"`
	Pets                []PetSummary `json:"pets"`
	TodayWalkCount      int64        `json:"todayWalkCount"`
	UncheckedAlertCount int64        `json:"uncheckedAlertCount"`
}

type UserSummary struct {
	NickName   string `json:"nickName"`
	ProfileUrl string `json:"profileUrl"`
}

type PetSummary struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ProfileImage string `json:"profileImage"`
}

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/api/v1/home", middleware.JWTAuth(), h.getHome)
}

func (h *Handler) getHome(c *gin.Context) {
	loginUser := middleware.GetLoginUser(c)
	userID := loginUser.UserID

	var user model.User
	h.db.First(&user, userID)

	var pets []model.Pet
	h.db.Joins("JOIN family_members fm ON fm.pet_id = pets.id").
		Where("fm.user_id = ? AND fm.status = ?", userID, model.FamilyStatusDone).
		Find(&pets)

	petSummaries := make([]PetSummary, 0, len(pets))
	petIDs := make([]int64, 0, len(pets))
	for _, p := range pets {
		petSummaries = append(petSummaries, PetSummary{ID: p.ID, Name: p.Name, ProfileImage: p.ProfileImage})
		petIDs = append(petIDs, p.ID)
	}

	var walkCount int64
	if len(petIDs) > 0 {
		h.db.Model(&model.Walk{}).
			Where("pet_id IN ? AND DATE(start_time) = CURDATE()", petIDs).
			Count(&walkCount)
	}

	var alertCount int64
	h.db.Model(&model.AlertHistory{}).
		Where("user_id = ? AND alert_history_status = ?", userID, model.AlertHistoryStatusUnchecked).
		Count(&alertCount)

	response.OK(c, HomeResponse{
		User:                UserSummary{NickName: user.NickName, ProfileUrl: user.ProfileUrl},
		Pets:                petSummaries,
		TodayWalkCount:      walkCount,
		UncheckedAlertCount: alertCount,
	})
}
