package walk

import (
	"strconv"
	"time"

	"github.com/PuppyNote/puppynote-server-golang/pkg/middleware"
	"github.com/PuppyNote/puppynote-server-golang/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/api/v1/walks", middleware.JWTAuth())
	{
		g.POST("", h.create)
		g.GET("/:walkId", h.getDetail)
		g.GET("", h.getByDate)
		g.GET("/calendar", h.getCalendar)
		g.DELETE("/:walkId", h.delete)
	}
}

func (h *Handler) create(c *gin.Context) {
	var req WalkCreateRequest
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

func (h *Handler) getDetail(c *gin.Context) {
	walkID, err := strconv.ParseInt(c.Param("walkId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 walkId입니다.")
		return
	}
	res, err := h.svc.GetDetail(walkID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) getByDate(c *gin.Context) {
	petID, err := strconv.ParseInt(c.Query("petId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "petId가 필요합니다.")
		return
	}
	dateStr := c.Query("date")
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response.BadRequest(c, "date 형식은 yyyy-MM-dd입니다.")
		return
	}
	res, err := h.svc.GetByDate(petID, date)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) getCalendar(c *gin.Context) {
	petID, err := strconv.ParseInt(c.Query("petId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "petId가 필요합니다.")
		return
	}
	yearMonthStr := c.Query("yearMonth")
	t, err := time.Parse("2006-01", yearMonthStr)
	if err != nil {
		response.BadRequest(c, "yearMonth 형식은 yyyy-MM입니다.")
		return
	}
	res, err := h.svc.GetCalendar(petID, t.Year(), int(t.Month()))
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) delete(c *gin.Context) {
	walkID, err := strconv.ParseInt(c.Param("walkId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 walkId입니다.")
		return
	}
	if err = h.svc.Delete(walkID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
