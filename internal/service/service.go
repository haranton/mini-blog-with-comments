package service

import (
	"blogWithComments/internal/models"
	"blogWithComments/internal/repository"
	"strconv"
)

type Service struct {
	repo *repository.Repo
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUserByLogin(login string) (*models.User, error) {
	return s.repo.GetUserByLogin(login)
}

func (s *Service) CreateUser(login, password string) *models.User {
	return s.repo.CreateUser(login, password)
}

func (s *Service) GetPosts(userID uint, limitStr, offsetStr string) *[]models.Post {

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 10
	}

	return s.repo.GetPosts(userID, limit, offset)
}

func (s *Service) CreatePost(userID uint, title string) *models.Post {
	return s.repo.CreatePost(userID, title)
}

func (s *Service) CreateComment(userID, postID uint, title string) *models.Comment {
	return s.repo.CreateComment(userID, postID, title)
}

func (s *Service) GetPostAndComments(postID uint) (*models.Post, error) {
	return s.repo.GetPostAndComments(postID)
}
