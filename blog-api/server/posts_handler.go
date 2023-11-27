package server

import (
	"encoding/json"
	"net/http"
)

type PostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type PostResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func (s *Server) GetAllPostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the context from the request
		ctx := r.Context()
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			http.Error(w, "error getting author from context", http.StatusInternalServerError)
			return
		}

		// Get all posts for the author
		posts, err := s.PostsService.GetAllPosts(ctx, author)
		if err != nil {
			http.Error(w, "error getting posts", http.StatusInternalServerError)
			return
		}

		// Write the posts to the response
		jsonResponse, err := json.Marshal(posts)
		if err != nil {
			http.Error(w, "error marshalling posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}

func (s *Server) GetPostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) CreatePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) UpdatePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *Server) DeletePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
