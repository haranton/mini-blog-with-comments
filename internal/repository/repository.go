package repository

import (
	"blogWithComments/internal/models"

	"gorm.io/gorm"
)

// todo пагинация

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r Repo) GetPosts(userID uint, limit, offset int) *[]models.Post {
	var Posts []models.Post
	r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Find(&Posts)
	return &Posts
}

func (r Repo) CreatePost(userID uint, title string) *models.Post {
	Post := models.Post{
		Title:  title,
		UserID: userID,
	}

	r.db.Create(&Post)
	return &Post
}

func (r Repo) CreateComment(userID, postID uint, title string) *models.Comment {
	comment := models.Comment{
		Title:  title,
		UserID: userID,
		PostID: postID,
	}
	r.db.Create(&comment)
	return &comment
}

func (r Repo) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r Repo) CreateUser(login, password string) *models.User {
	user := models.User{
		Login:    login,
		Password: password,
	}

	r.db.Create(&user)
	return &user
}

func (r Repo) GetPostAndComments(postID uint) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Comments").First(&post, postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
