package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"rakia.ai/blog-api/v2/internal"
)

func writeJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

type PostsService interface {
	CreatePosts(post internal.Post, author string) error
	GetAllPosts() ([]*internal.Post, error)
	UpdatePosts(post internal.Post, author string) error
	GetPostByID(id int) (*internal.Post, error)
	DeletePosts(id int, author string) error
}

type AuthorsService interface {
	ValidAuthor(username string, password string) (bool, error)
}

type Server struct {
	Router         *mux.Router
	PostsService   PostsService
	AuthorsService AuthorsService
	Logger         *zerolog.Logger
}

func NewLogger() *zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	return &logger
}

func NewServer(router *mux.Router, blogs PostsService, authors AuthorsService, logger *zerolog.Logger) *Server {
	return &Server{
		Router:         router,
		PostsService:   blogs,
		AuthorsService: authors,
		Logger:         logger,
	}
}

func (s *Server) Routes() {

	// Login Author and get a JWT
	s.Router.HandleFunc("/login", s.LoginHandler()).Methods("POST")

	api := s.Router.PathPrefix("/api").Subrouter()

	// Authenticated routes
	api.Use(Middleware(*s.Logger))

	// Create a new post for an author
	api.HandleFunc("/posts", s.CreatePostsHandler()).Methods("POST")
	// Get one post for an author
	api.HandleFunc("/posts/{id}", s.GetPostsHandler()).Methods("GET")
	// Get all posts for an author
	api.HandleFunc("/posts", s.GetAllPostsHandler()).Methods("GET")
	// Update a post for an author
	api.HandleFunc("/posts/{id}", s.UpdatePostsHandler()).Methods("PUT")
	// Delete a post for an author
	api.HandleFunc("/posts/{id}", s.DeletePostsHandler()).Methods("DELETE")

}
