package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"blogWithComments/internal/handlers"
	"blogWithComments/internal/models"
	"blogWithComments/internal/repository"
	"blogWithComments/internal/service"
)

func setupAppForUsers(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	// In-memory SQLite DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}

	repo := repository.NewRepo(db)
	svc := service.NewService(repo)
	h := handlers.NewHandler(svc)

	r := gin.New()
	handlers.RegisterRoutes(r, h)
	return r
}

type createUserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func TestCreateUserAndGetUser(t *testing.T) {
	r := setupAppForUsers(t)

	// Create user
	payload := createUserReq{Login: "john", Password: "secret"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST /users status=%d body=%s", w.Code, w.Body.String())
	}

	// Get user by id (expect 200)
	req2 := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("GET /users/1 status=%d body=%s", w2.Code, w2.Body.String())
	}
}
