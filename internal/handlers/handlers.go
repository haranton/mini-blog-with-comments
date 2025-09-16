package handlers

import (
	"blogWithComments/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

func RegisterRoutes(r *gin.Engine, h *Handler) {
	r.GET("/users", h.ListUsers) // ?login=
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id", h.GetUserByID)
	r.GET("/users/:id/posts", h.GetPosts)

	r.POST("/posts", h.CreatePost)
	r.GET("/posts/:postid", h.GetPostWithComments)
	r.POST("/posts/:postid/comments", h.CreateComment)
	// совместимость с старым тестом
	r.POST("/comments", h.CreateCommentLegacy)
}

type createUserReq struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.CreateUser(req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) ListUsers(c *gin.Context) {
	login := c.Query("login")
	if login == "" {
		// TODO: вернуть список всех пользователей (можно добавить репо-метод)
		c.JSON(http.StatusOK, []any{})
		return
	}
	user, err := h.svc.GetUserByLogin(login)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, []any{user})
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	user, err := h.svc.GetUserByID(uint(id64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

type createPostReq struct {
	UserID uint   `json:"user_id" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

func (h *Handler) CreatePost(c *gin.Context) {
	var req createPostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, err := h.svc.CreatePost(req.UserID, req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetPosts(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	posts, err := h.svc.GetPosts(uint(id64), limitStr, offsetStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": posts, "meta": gin.H{"limit": limitStr, "offset": offsetStr}})
}

type createCommentReq struct {
	UserID uint   `json:"user_id" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

func (h *Handler) CreateComment(c *gin.Context) {
	postID64, err := strconv.ParseUint(c.Param("postid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req createCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment, err := h.svc.CreateComment(req.UserID, uint(postID64), req.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

// Legacy route for tests expecting POST /comments with body containing post_id
func (h *Handler) CreateCommentLegacy(c *gin.Context) {
	var body struct {
		UserID uint   `json:"user_id" binding:"required"`
		PostID uint   `json:"post_id" binding:"required"`
		Title  string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment, err := h.svc.CreateComment(body.UserID, body.PostID, body.Title)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) GetPostWithComments(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("postid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	post, err := h.svc.GetPostAndComments(uint(id64), limitStr, offsetStr)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}
