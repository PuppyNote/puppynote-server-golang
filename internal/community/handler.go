package community

import (
	"strconv"

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
	g := rg.Group("/api/v1/community", middleware.JWTAuth())
	{
		g.POST("/posts", h.createPost)
		g.GET("/posts", h.getPosts)
		g.GET("/posts/my", h.getMyPosts)
		g.GET("/posts/:postId", h.getPost)
		g.PATCH("/posts/:postId", h.updatePost)
		g.DELETE("/posts/:postId", h.deletePost)
		g.GET("/posts/hashtags", h.searchHashtags)
		g.POST("/posts/:postId/like", h.toggleLike)
	}
}

func (h *Handler) createPost(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	var req PostCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	postID, err := h.svc.CreatePost(user.UserID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, postID)
}

func (h *Handler) getPosts(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	res, err := h.svc.GetPosts(user.UserID, page, size)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) getMyPosts(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	res, err := h.svc.GetMyPosts(user.UserID, page, size)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) getPost(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 postId입니다.")
		return
	}
	res, err := h.svc.GetPost(postID, user.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) updatePost(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 postId입니다.")
		return
	}
	var req PostUpdateRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err = h.svc.UpdatePost(postID, user.UserID, req); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) deletePost(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 postId입니다.")
		return
	}
	if err = h.svc.DeletePost(postID, user.UserID); err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) searchHashtags(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.BadRequest(c, "keyword가 필요합니다.")
		return
	}
	res, err := h.svc.SearchHashtags(keyword)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) toggleLike(c *gin.Context) {
	user := middleware.GetLoginUser(c)
	postID, err := strconv.ParseInt(c.Param("postId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "잘못된 postId입니다.")
		return
	}
	res, err := h.svc.ToggleLike(postID, user.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.OK(c, res)
}
