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

func (r *Repo) GetPosts(userID uint, limit, offset int) ([]models.Post, error) {
	var posts []models.Post
	if err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Preload("Comments").
		Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *Repo) CreatePost(userID uint, title string) (*models.Post, error) {
	post := models.Post{
		Title:  title,
		UserID: userID,
	}

	if err := r.db.Create(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *Repo) CreateComment(userID, postID uint, title string) (*models.Comment, error) {
	comment := models.Comment{
		Title:  title,
		UserID: userID,
		PostID: postID,
	}
	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *Repo) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("login = ?", login).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) CreateUser(login, password string) (*models.User, error) {
	user := models.User{
		Login:    login,
		Password: password,
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetPostAndComments(postID uint, limit, offset int) (*models.Post, error) {
	var post models.Post
	if err := r.db.Preload("Comments").First(&post, postID).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
