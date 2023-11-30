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

func (m *MockPostsService) GetAllPosts(author string) ([]*internal.Post, error) {
	args := m.Called(author)
	return args.Get(0).([]*internal.Post), args.Error(1)
}

func (m *MockPostsService) UpdatePosts(post internal.Post, author string) error {
	args := m.Called(post, author)
	return args.Error(0)
}

func (m *MockPostsService) GetPostByID(id int, author string) (*internal.Post, error) {
	args := m.Called(id, author)
	return args.Get(0).(*internal.Post), args.Error(1)
}

func (m *MockPostsService) DeletePosts(id int, author string) error {
	args := m.Called(id, author)
	return args.Error(0)
}

var testPost = internal.Post{
	ID:      1,
	Title:   "Test Post 1",
	Content: "Content 1",
	Author:  "Author 1",
}

var testPostUpdate = internal.Post{
	ID:      1,
	Title:   "Test Post 1",
	Content: "Content 33333",
	Author:  "Author 1",
}

var testPostCreate = internal.Post{
	Title:   "Test Post 2",
	Content: "Content 2",
	Author:  "Author 1",
}

var logger = zerolog.New(os.Stdout)

// TestGetAllPostsHandler tests the GetAllPostsHandler function
func TestGetAllPostsHandler(t *testing.T) {
	// Create a mock instance of the PostsService
	mockPostsService := new(MockPostsService)
	mockPosts := []*internal.Post{}
	mockPosts = append(mockPosts, &testPost)
	// Create a logger instance or mock

	mockPostsService.On("GetAllPosts", "Author 1").Return(mockPosts, nil)

	// Create an instance of the Server with the mock service
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/api/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Adding context with author value
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 1")
	req = req.WithContext(ctx)

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetAllPostsHandler())

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expectedResponse, _ := json.Marshal(mockPosts)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())
}

// TestGetPostsHandler tests the GetPostsHandler function
func TestGetPostsHandler(t *testing.T) {
	// Create a mock instance of the PostsService
	mockPostsService := new(MockPostsService)
	mockPost := &testPost
	mockPostsService.On("GetPostByID", 1, "Author 1").Return(mockPost, nil)

	// Create an instance of the Server with the mock service
	server := &Server{PostsService: mockPostsService, Logger: &logger}

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/api/posts/1", nil)
	// Add the id parameter to the request
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	if err != nil {
		t.Fatal(err)
	}

	// Adding context with author value
	ctx := context.WithValue(req.Context(), ContextAuthor, "Author 1")
	req = req.WithContext(ctx)

	// Record the response using httptest
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetPostsHandler())

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body
	expectedResponse, _ := json.Marshal(mockPost)
	assert.JSONEq(t, string(expectedResponse), rr.Body.String())
}

// TestCreatePostsHandler tests the CreatePostsHandler function
func TestCreatePostsHandler(t *testing.T) {

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
	handler := http.HandlerFunc(server.CreatePostsHandler())

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusCreated, rr.Code)
}

// TestUpdatePostsHandler tests the UpdatePostsHandler function
func TestUpdatePostsHandler(t *testing.T) {

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
	handler := http.HandlerFunc(server.UpdatePostsHandler())

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusAccepted, rr.Code)
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
	handler := http.HandlerFunc(server.DeletePostsHandler())

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusAccepted, rr.Code)
}
