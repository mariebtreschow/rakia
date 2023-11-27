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
	CreatePosts(post internal.Post) error
	GetAllPosts(author string) ([]*internal.Post, error)
	UpdatePosts(post internal.Post, author string) error
	GetPosts(id int, author string) (*internal.Post, error)
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

	api.HandleFunc("/posts", s.CreatePostsHandler()).Methods("POST")
	api.HandleFunc("/posts/{id}", s.GetPostsHandler()).Methods("GET")
	api.HandleFunc("/posts", s.GetAllPostsHandler()).Methods("GET")
	api.HandleFunc("/posts/{id}", s.UpdatePostsHandler()).Methods("PUT")
	api.HandleFunc("/posts/{id}", s.DeletePostsHandler()).Methods("DELETE")

}
