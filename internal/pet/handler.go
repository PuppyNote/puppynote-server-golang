package pet

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
	g := rg.Group("/api/v1/pets", middleware.JWTAuth())
	{
		g.GET("", h.getMyPets)
		g.POST("", h.createPet)
		g.PATCH("/:petId", h.updatePet)
		g.DELETE("/:petId", h.deletePet)
	}
}

func (h *Handler) getMyPets(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	res, err := h.svc.GetMyPets(user.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) createPet(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	var req PetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.CreatePet(user.UserID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) updatePet(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	petID, err := strconv.ParseInt(c.Param("petId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 petId입니다.")
		return
	}
	var req PetUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err = h.svc.UpdatePet(user.UserID, petID, req); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) deletePet(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	petID, err := strconv.ParseInt(c.Param("petId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 petId입니다.")
		return
	}
	if err = h.svc.DeletePet(user.UserID, petID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
