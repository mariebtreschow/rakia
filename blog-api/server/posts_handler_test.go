package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"rakia.ai/blog-api/v2/internal"
)

// MockPostsService is a mock implementation of the PostsService interface
type MockPostsService struct {
	mock.Mock
}

func (m *MockPostsService) CreatePosts(post internal.Post, author string) error {
	args := m.Called(post, author)
	return args.Error(0)
}

func (m *MockPostsService) GetAllPosts() ([]*internal.Post, error) {
	args := m.Called()
	return args.Get(0).([]*internal.Post), args.Error(1)
}

func (m *MockPostsService) UpdatePosts(post internal.Post, author string) error {
	args := m.Called(post, author)
	return args.Error(0)
}

func (m *MockPostsService) GetPostByID(id int) (*internal.Post, error) {
	args := m.Called(id)
	return args.Get(0).(*internal.Post), args.Error(1)
}

func (m *MockPostsService) DeletePosts(id int, author string) error {
	args := m.Called(id, author)
	return args.Error(0)
}

var logger = zerolog.New(os.Stdout)

// TestGetAllPostsHandler tests the GetAllPostsHandler function
func TestGetAllPostsHandler(t *testing.T) {

	testPost := internal.Post{
		ID:      1,
		Title:   "Test Post 1",
		Content: "Content 1",
		Author:  "Author 1",
	}

	// Create a mock instance of the PostsService
	mockPostsService := new(MockPostsService)
	mockPosts := []*internal.Post{}
	mockPosts = append(mockPosts, &testPost)
	// Create a logger instance or mock

	mockPostsService.On("GetAllPosts").Return(mockPosts, nil)

	// Create an instance of the Server with the mock service
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := server.GetAllPostsHandler()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if http.StatusOK != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse, _ := json.Marshal(mockPosts)

	if string(expectedResponse) != rr.Body.String() {
		t.Fatalf("expected %v; got %v", string(expectedResponse), rr.Body.String())
	}

}

// TestGetPostsHandler tests the GetPostsHandler function
func TestGetPostsHandler(t *testing.T) {

	testPost := internal.Post{
		ID:      1,
		Title:   "Test Post 1",
		Content: "Content 1",
		Author:  "Author 1",
	}
	// Create a mock instance of the PostsService
	mockPostsService := new(MockPostsService)
	mockPost := &testPost
	mockPostsService.On("GetPostByID", 1).Return(mockPost, nil)

	// Create an instance of the Server with the mock service
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/api/posts/1", nil)
	// Add the id parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	if err != nil {
		t.Fatal(err)
	}

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := server.GetPostsHandler()

	// Call the handler
	handler.ServeHTTP(rr, req)

	if http.StatusOK != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse, _ := json.Marshal(mockPost)

	if string(expectedResponse) != rr.Body.String() {
		t.Fatalf("expected %v; got %v", string(expectedResponse), rr.Body.String())
	}

}

// TestCreatePostsHandler tests the CreatePostsHandler function
func TestCreatePostsHandler(t *testing.T) {

	testPostCreate := internal.Post{
		Title:   "Test Post 2",
		Content: "Content 2",
		Author:  "Author 1",
	}

	// Create a mock instance of the PostsService
	mockPostsService := new(MockPostsService)
	mockPostsService.On("CreatePosts", testPostCreate, "Author 1").Return(nil)

	// Create an instance of the Server with the mock service
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	// Create a request to pass to the handler
	jsonPost, err := json.Marshal(testPostCreate)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(jsonPost))
	if err != nil {
		t.Fatal(err)
	}

	// Adding context with author value
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 1")
	req = req.WithContext(ctx)

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := server.CreatePostsHandler()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if http.StatusCreated != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusCreated, rr.Code)
	}
}

// TestUpdatePostsHandler tests the UpdatePostsHandler function
func TestUpdatePostsHandler(t *testing.T) {

	testPostUpdate := internal.Post{
		ID:      1,
		Title:   "Test Post 1",
		Content: "Content 33333",
		Author:  "Author 1",
	}
	// Create a mock instance of the PostsService
	mockPostsService := new(MockPostsService)
	mockPostsService.On("UpdatePosts", testPostUpdate, "Author 1").Return(nil)

	// Create an instance of the Server with the mock service
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	// Create a request to pass to the handler
	jsonPost, err := json.Marshal(testPostUpdate)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("PUT", "/api/posts/1", bytes.NewBuffer(jsonPost))
	if err != nil {
		t.Fatal(err)
	}
	// Add the id parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Adding context with author value
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 1")
	req = req.WithContext(ctx)

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := server.UpdatePostsHandler()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if http.StatusAccepted != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusAccepted, rr.Code)
	}
}

// TestDeletePostsHandler tests the DeletePostsHandler function
func TestDeletePostsHandler(t *testing.T) {

	// Create a mock instance of the PostsService
	mockPostsService := new(MockPostsService)
	mockPostsService.On("DeletePosts", 1, "Author 1").Return(nil)

	// Create an instance of the Server with the mock service
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	// Create a request to pass to the handler
	req, err := http.NewRequest("DELETE", "/api/posts/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the id parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Adding context with author value
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 1")
	req = req.WithContext(ctx)

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := server.DeletePostsHandler()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if http.StatusAccepted != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusAccepted, rr.Code)
	}

}

func TestGetPostNotDeletedPostFoundHandler(t *testing.T) {
	mockPostsService := new(MockPostsService)
	mockPostsService.On("GetPostByID", 1).Return(&internal.Post{}, internal.ErrPostNotFound)
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	req, _ := http.NewRequest("GET", "/api/posts/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler := server.GetPostsHandler()

	handler.ServeHTTP(rr, req)

	if http.StatusNotFound != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusNotFound, rr.Code)
	}
}

func TestForbiddenDeletedPostFoundHandler(t *testing.T) {
	mockPostsService := new(MockPostsService)
	mockPostsService.On("DeletePosts", 1, "Author 3").Return(internal.ErrAuthorNotAllowed)

	server := &Server{PostsService: mockPostsService, Logger: &logger}

	req, _ := http.NewRequest("DELETE", "/api/posts/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	// Adding context with author value
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 3")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := server.DeletePostsHandler()

	handler.ServeHTTP(rr, req)

	if http.StatusForbidden != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusForbidden, rr.Code)
	}

}

func TestGetPostNotFoundHandler(t *testing.T) {
	mockPostsService := new(MockPostsService)
	mockPostsService.On("GetPostByID", 99).Return(&internal.Post{}, internal.ErrPostNotFound)
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	req, _ := http.NewRequest("GET", "/api/posts/99", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "99"})

	rr := httptest.NewRecorder()
	handler := server.GetPostsHandler()

	handler.ServeHTTP(rr, req)

	if http.StatusNotFound != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusNotFound, rr.Code)
	}
}

func TestUpdateInvalidPostsHandler(t *testing.T) {
	invalidPostUpdate := internal.Post{
		ID:      1,
		Title:   "Valid Title",
		Content: "", // Invalid because content is empty
		Author:  "Author 1",
	}

	mockPostsService := new(MockPostsService)
	mockPostsService.On("UpdatePosts", invalidPostUpdate, "Author 1").Return(internal.ErrContentEmpty)

	server := &Server{PostsService: mockPostsService, Logger: &logger}

	jsonPost, _ := json.Marshal(invalidPostUpdate)
	req, _ := http.NewRequest("PUT", "/api/posts/1", bytes.NewBuffer(jsonPost))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := server.UpdatePostsHandler()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateInvalidPostsHandler(t *testing.T) {
	invalidPost := internal.Post{
		Title:   "", // Invalid because the title is empty
		Content: "Some Content",
		Author:  "Author 1",
	}

	mockPostsService := new(MockPostsService)
	mockPostsService.On("CreatePosts", invalidPost, "Author 1").Return(internal.ErrTitleEmpty)

	server := &Server{PostsService: mockPostsService, Logger: &logger}

	jsonPost, _ := json.Marshal(invalidPost)
	req, _ := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(jsonPost))
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler := server.CreatePostsHandler()

	handler.ServeHTTP(rr, req)

	if http.StatusBadRequest != rr.Code {
		t.Fatalf("expected %v; got %v", http.StatusBadRequest, rr.Code)
	}
}
