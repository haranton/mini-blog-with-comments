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

func setupAppForPosts(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

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

func createTestUser(t *testing.T, r *gin.Engine, login, password string) {
	t.Helper()
	payload := map[string]any{"login": login, "password": password}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST /users status=%d body=%s", w.Code, w.Body.String())
	}
}

func TestCreatePostAndGetPosts(t *testing.T) {
	r := setupAppForPosts(t)
	createTestUser(t, r, "john", "secret")

	// create post
	postPayload := map[string]any{"user_id": 1, "title": "hello"}
	postBody, _ := json.Marshal(postPayload)
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST /posts status=%d body=%s", w.Code, w.Body.String())
	}

	// get posts by user
	req2 := httptest.NewRequest(http.MethodGet, "/users/1/posts", nil)
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("GET /users/1/posts status=%d body=%s", w2.Code, w2.Body.String())
	}
}

func TestCreateComment(t *testing.T) {
	r := setupAppForPosts(t)
	createTestUser(t, r, "john", "secret")

	// create post
	postPayload := map[string]any{"user_id": 1, "title": "hello"}
	postBody, _ := json.Marshal(postPayload)
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST /posts status=%d body=%s", w.Code, w.Body.String())
	}

	// create comment
	commentPayload := map[string]any{"user_id": 1, "post_id": 1, "title": "my comment"}
	commentBody, _ := json.Marshal(commentPayload)
	reqC := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewReader(commentBody))
	reqC.Header.Set("Content-Type", "application/json")
	wC := httptest.NewRecorder()
	r.ServeHTTP(wC, reqC)
	if wC.Code != http.StatusCreated {
		t.Fatalf("POST /comments status=%d body=%s", wC.Code, wC.Body.String())
	}
}

func TestCreatePost_ValidationError(t *testing.T) {
	r := setupAppForPosts(t)
	createTestUser(t, r, "john", "secret")

	// missing title
	postPayload := map[string]any{"user_id": 1}
	postBody, _ := json.Marshal(postPayload)
	req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewReader(postBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing title, got=%d body=%s", w.Code, w.Body.String())
	}
}
