package family

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
	g := rg.Group("/api/v1/family-members", middleware.JWTAuth())
	{
		g.GET("", h.getFamilyMembers)
		g.GET("/search", h.searchUsers)
		g.POST("/invite", h.inviteFamilyMember)
		g.POST("/register", h.registerFamilyMember)
		g.DELETE("/:targetUserId", h.deleteFamilyMember)
	}
}

func (h *Handler) getFamilyMembers(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	petID, err := strconv.ParseInt(c.Query("petId"), 10, 64)
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
	petID, err := strconv.ParseInt(c.Query("petId"), 10, 64)
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
