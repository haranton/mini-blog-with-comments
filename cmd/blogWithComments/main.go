package main

import (
	"blogWithComments/internal/config"
	"blogWithComments/internal/handlers"
	"blogWithComments/internal/models"
	"blogWithComments/internal/repository"
	"blogWithComments/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {

	db := config.ConnectDB()
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	repo := repository.NewRepo(db)
	service := service.NewService(repo)
	handlerObj := handlers.NewHandler(service)

	r := gin.Default()

	handlers.RegisterRoutes(r, handlerObj)

	r.Run(":8080")

}
