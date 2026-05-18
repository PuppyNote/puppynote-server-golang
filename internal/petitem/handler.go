package petitem

import (
	"strconv"

	"puppynote/internal/model"
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
	g := rg.Group("/api/v1/pet-items", middleware.JWTAuth())
	{
		g.POST("", h.create)
		g.GET("", h.list)
		g.GET("/:petItemId", h.getDetail)
		g.PATCH("/:petItemId", h.update)
		g.DELETE("/:petItemId", h.delete)
		g.POST("/:petItemId/purchases", h.recordPurchase)
		g.GET("/:petItemId/purchases", h.getPurchaseHistory)
		g.DELETE("/purchases/:purchaseId", h.deletePurchase)
	}
}

func (h *Handler) create(c *gin.Context) {
	var req PetItemCreateRequest
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

func (h *Handler) list(c *gin.Context) {
	petID, err := strconv.ParseInt(c.Query("petId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "petId가 필요합니다.")
		return
	}

	var category *model.ItemCategory
	if cat := c.Query("category"); cat != "" {
		c := model.ItemCategory(cat)
		category = &c
	}

	res, err := h.svc.GetByPetID(petID, category)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) getDetail(c *gin.Context) {
	petItemID, err := strconv.ParseInt(c.Param("petItemId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 petItemId입니다.")
		return
	}
	res, err := h.svc.GetDetail(petItemID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) update(c *gin.Context) {
	petItemID, err := strconv.ParseInt(c.Param("petItemId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 petItemId입니다.")
		return
	}
	var req PetItemUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.Update(petItemID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) delete(c *gin.Context) {
	petItemID, err := strconv.ParseInt(c.Param("petItemId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 petItemId입니다.")
		return
	}
	if err = h.svc.Delete(petItemID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) recordPurchase(c *gin.Context) {
	petItemID, err := strconv.ParseInt(c.Param("petItemId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 petItemId입니다.")
		return
	}
	var req PetItemPurchaseCreateRequest
	_ = c.ShouldBindJSON(&req)

	res, err := h.svc.RecordPurchase(petItemID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) getPurchaseHistory(c *gin.Context) {
	petItemID, err := strconv.ParseInt(c.Param("petItemId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 petItemId입니다.")
		return
	}
	res, err := h.svc.GetPurchaseHistory(petItemID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) deletePurchase(c *gin.Context) {
	purchaseID, err := strconv.ParseInt(c.Param("purchaseId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 purchaseId입니다.")
		return
	}
	if err = h.svc.DeletePurchase(purchaseID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
