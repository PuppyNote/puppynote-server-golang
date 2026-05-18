package alert

import (
	"math"
	"strconv"

	"puppynote/pkg/middleware"
	"puppynote/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/api/v1", middleware.JWTAuth())
	{
		g.GET("/alert-setting", h.getSetting)
		g.PATCH("/alert-setting", h.updateSetting)
		g.GET("/alertHistories", h.getHistories)
		g.GET("/alertHistories/unchecked", h.hasUnchecked)
		g.PATCH("/alertHistories/:id", h.updateHistoryStatus)
	}
}

func (h *Handler) getSetting(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	s, err := h.repo.FindSettingByUserID(user.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, AlertSettingResponse{All: s.All, Walk: s.Walk, Friend: s.Friend})
}

func (h *Handler) updateSetting(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	var req AlertSettingUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	s, err := h.repo.UpdateSetting(user.UserID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, AlertSettingResponse{All: s.All, Walk: s.Walk, Friend: s.Friend})
}

func (h *Handler) getHistories(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	histories, total, err := h.repo.FindHistoriesByUserID(user.UserID, page, size)
	if err != nil {
		_ = c.Error(err)
		return
	}

	content := make([]AlertHistoryResponse, 0, len(histories))
	for _, h := range histories {
		content = append(content, AlertHistoryResponse{
			ID:                   h.ID,
			AlertDescription:     h.AlertDescription,
			AlertHistoryStatus:   h.AlertHistoryStatus,
			AlertDestinationType: h.AlertDestinationType,
			AlertDestinationInfo: h.AlertDestinationInfo,
			CreatedDate:          h.CreatedDate,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(size)))
	response.OK(c, AlertHistoryListResponse{
		Content:       content,
		Page:          page,
		Size:          size,
		TotalElements: total,
		TotalPages:    totalPages,
	})
}

func (h *Handler) hasUnchecked(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	response.OK(c, h.repo.HasUnchecked(user.UserID))
}

func (h *Handler) updateHistoryStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 id입니다.")
		return
	}
	if err = h.repo.UpdateHistoryStatus(id); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
