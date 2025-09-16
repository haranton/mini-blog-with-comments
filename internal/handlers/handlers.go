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
	r.GET("/users", h.GetUser)
	r.POST("/users", h.CreateUser)
	r.GET("/users/:id/posts", h.GetPosts)

	r.POST("/posts", h.CreatePost)
	r.POST("/posts/:postid/comments", h.CreateComment)

	// Комментарии
	r.GET("/post/:postid", h.GetPostWithComments)
}

type createUserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	user := h.svc.CreateUser(req.Login, req.Password)
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	login := c.Query("login")

	user, err := h.svc.GetUserByLogin(login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login"})
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
	post := h.svc.CreatePost(req.UserID, req.Title)
	c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetPosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	limitStr, offset := getLimitAndOffset(c)

	posts := h.svc.GetPosts(uint(id), limitStr, offset)
	c.JSON(http.StatusOK, posts)
}

type createCommentReq struct {
	UserID uint   `json:"user_id" binding:"required"`
	PostID uint   `json:"post_id" binding:"required"`
	Title  string `json:"title" binding:"required"`
}

func (h *Handler) CreateComment(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req createCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment := h.svc.CreateComment(req.UserID, uint(postID), req.Title)
	c.JSON(http.StatusCreated, comment)
}

func (h *Handler) GetPostWithComments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("postid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	limitStr, offset := getLimitAndOffset(c)

	post, err := h.svc.GetPostAndComments(uint(id), limitStr, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, post)
}

func getLimitAndOffset(c *gin.Context) (string, string) {
	limitStr := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", "0")
	return limitStr, offset
}
