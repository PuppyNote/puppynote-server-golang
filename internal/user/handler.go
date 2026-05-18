package user

import (
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
	public := rg.Group("/api/v1/user")
	{
		public.POST("/email/send", h.sendVerificationEmail)
		public.POST("/signup", h.signUp)
	}

	auth := rg.Group("/api/v1/user", middleware.JWTAuth())
	{
		auth.GET("/profile", h.getProfile)
		auth.PATCH("/profile", h.updateProfile)
		auth.DELETE("/withdraw", h.withdraw)
	}
}

func (h *Handler) sendVerificationEmail(c *gin.Context) {
	var req EmailSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SendVerificationEmail(req.Email); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) signUp(c *gin.Context) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	res, err := h.svc.SignUp(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) getProfile(c *gin.Context) {
	loginUser := middleware.GetLoginUser(c)
	res, err := h.svc.GetProfile(loginUser.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) updateProfile(c *gin.Context) {
	loginUser := middleware.GetLoginUser(c)
	var req UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateProfile(loginUser.UserID, req); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) withdraw(c *gin.Context) {
	loginUser := middleware.GetLoginUser(c)
	if err := h.svc.Withdraw(loginUser.UserID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
