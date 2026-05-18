package pettip

import (
	"puppynote/internal/model"
	"puppynote/pkg/middleware"
	"puppynote/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/api/v1/pet-tips/random", middleware.JWTAuth(), h.getRandom)
}

func (h *Handler) getRandom(c *gin.Context) {
	var tip model.PetTip
	if err := h.db.Order("RAND()").First(&tip).Error; err != nil {
		response.OK(c, nil)
		return
	}
	response.OK(c, tip)
}
