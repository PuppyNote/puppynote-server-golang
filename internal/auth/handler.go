package auth

import (
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
	g := rg.Group("/api/v1/auth")
	{
		g.POST("/login", h.login)
		g.POST("/oauth/login", h.oauthLogin)
		g.POST("/refresh", h.refresh)
		g.POST("/password/email/send", h.sendPasswordEmail)
		g.POST("/password/reset", h.resetPassword)
	}
}

func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.NormalLogin(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) oauthLogin(c *gin.Context) {
	var req LoginOauthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.OAuthLogin(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) refresh(c *gin.Context) {
	var req TokenRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	res, err := h.svc.Refresh(req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) sendPasswordEmail(c *gin.Context) {
	var req PasswordEmailSendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.SendPasswordResetEmail(req); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) resetPassword(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.ResetPassword(req); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}
