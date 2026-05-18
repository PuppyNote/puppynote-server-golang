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
	g := rg.Group("/api/v1", middleware.JWTAuth())
	{
		// Pet
		g.GET("/pets", h.getMyPets)
		g.POST("/pets", h.createPet)
		g.PATCH("/pets/:petId", h.updatePet)
		g.DELETE("/pets/:petId", h.deletePet)

		// FamilyMember
		g.GET("/family-members", h.getFamilyMembers)
		g.GET("/family-members/search", h.searchUsers)
		g.POST("/family-members/invite", h.inviteFamilyMember)
		g.POST("/family-members/register", h.registerFamilyMember)
		g.DELETE("/family-members/:targetUserId", h.deleteFamilyMember)
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

func (h *Handler) getFamilyMembers(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	petIDStr := c.Query("petId")
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "petId가 필요합니다.")
		return
	}
	res, err := h.svc.GetFamilyMembers(user.UserID, petID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) searchUsers(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	email := c.Query("email")
	if email == "" {
		response.BadRequest(c, "email 쿼리 파라미터가 필요합니다.")
		return
	}
	res, err := h.svc.SearchUsersByEmail(user.UserID, email)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) inviteFamilyMember(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	var req FamilyMemberInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.InviteFamilyMember(user.UserID, req); err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, nil)
}

func (h *Handler) registerFamilyMember(c *gin.Context) {
	var req FamilyMemberRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.RegisterFamilyMember(req); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) deleteFamilyMember(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	targetUserID, err := strconv.ParseInt(c.Param("targetUserId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 targetUserId입니다.")
		return
	}
	petIDStr := c.Query("petId")
	petID, err := strconv.ParseInt(petIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "petId가 필요합니다.")
		return
	}
	if err = h.svc.DeleteFamilyMember(user.UserID, targetUserID, petID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
