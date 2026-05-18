package walkalarm

import (
	"strconv"

	"puppynote/pkg/middleware"
	"puppynote/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/api/v1/pet-walk-alarms", middleware.JWTAuth())
	{
		g.GET("", h.getByPetID)
		g.POST("", h.create)
		g.PUT("", h.update)
		g.PATCH("/status", h.updateStatus)
		g.DELETE("/:alarmId", h.delete)
	}
}

func (h *Handler) getByPetID(c *gin.Context) {
	petID, err := strconv.ParseInt(c.Query("petId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "petId가 필요합니다.")
		return
	}
	res, err := h.svc.GetByPetID(petID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) create(c *gin.Context) {
	var req WalkAlarmCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.Create(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) update(c *gin.Context) {
	var req WalkAlarmUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.Update(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) updateStatus(c *gin.Context) {
	var req WalkAlarmStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.UpdateStatus(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) delete(c *gin.Context) {
	alarmID, err := strconv.ParseInt(c.Param("alarmId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 alarmId입니다.")
		return
	}
	if err = h.svc.Delete(alarmID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
