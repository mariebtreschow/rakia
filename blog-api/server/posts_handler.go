package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"rakia.ai/blog-api/v2/internal"
)

type PostCreate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type PostUpdate struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
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
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			http.Error(w, "error getting author from context", http.StatusInternalServerError)
			return
		}

		// Get all posts for the author
		posts, err := s.PostsService.GetAllPosts(author)
		if err != nil {
			writeJSONError(w, "error getting posts", http.StatusInternalServerError)
			return
		}
		if len(posts) == 0 {
			writeJSONError(w, "no posts found", http.StatusNotFound)
			return
		}

		// Write the posts to the response
		jsonResponse, err := json.Marshal(posts)
		if err != nil {
			writeJSONError(w, "error marshalling posts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}

func (s *Server) GetPostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			s.Logger.Error().Msg("error getting author from context")
			writeJSONError(w, "error getting author from context", http.StatusInternalServerError)
			return
		}

		// Get the post ID from the URL
		id, ok := mux.Vars(r)["id"]
		if !ok {
			s.Logger.Error().Msg("missing post id")
			writeJSONError(w, "missing post id", http.StatusBadRequest)
			return
		}

		// Convert the post ID from string to int
		postID, err := strconv.Atoi(id)
		if err != nil {
			s.Logger.Error().Err(err).Msg("invalid post id")
			writeJSONError(w, "invalid post id", http.StatusBadRequest)
			return
		}

		// Get the post
		post, err := s.PostsService.GetPosts(postID, author)
		if err != nil {
			s.Logger.Error().Err(err).Msg("error getting post")
			if err == internal.ErrPostNotFound {
				writeJSONError(w, "post not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, "error getting post", http.StatusInternalServerError)
			return
		}

		// return the post
		jsonResponse, err := json.Marshal(post)
		if err != nil {
			s.Logger.Error().Err(err).Msg("error marshalling post")
			writeJSONError(w, "error marshalling post", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	}
}

func (s *Server) CreatePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate incoming post
		var postRequest PostCreate
		err := json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			writeJSONError(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			writeJSONError(w, "error getting author from context", http.StatusInternalServerError)
			return
		}

		// Create the post
		post := internal.Post{
			Title:   postRequest.Title,
			Content: postRequest.Content,
		}

		// Add the author to the post
		post.Author = author

		// Save the post
		err = s.PostsService.CreatePosts(post)
		if err != nil {
			if err == internal.ErrUniqueTitle {
				writeJSONError(w, "title already exists", http.StatusBadRequest)
				return
			}
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Status accepted
		w.WriteHeader(http.StatusCreated)

	}
}

func (s *Server) UpdatePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			writeJSONError(w, "error getting author from context", http.StatusInternalServerError)
			return
		}

		// Get the post ID from the URL
		id, ok := mux.Vars(r)["id"]
		if !ok {
			writeJSONError(w, "missing post id", http.StatusBadRequest)
			return
		}

		// Convert the post ID from string to int
		postID, err := strconv.Atoi(id)
		if err != nil {
			s.Logger.Error().Err(err).Msg("invalid post id")
			writeJSONError(w, "invalid post id", http.StatusBadRequest)
			return
		}

		// validate incoming post
		var postRequest PostUpdate
		err = json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			s.Logger.Error().Err(err).Msg("invalid request payload")
			writeJSONError(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		if postID != postRequest.ID {
			writeJSONError(w, "mismatching ids in request and url", http.StatusBadRequest)
			return
		}

		var post internal.Post
		// Update the post
		post.ID = postRequest.ID // Not being update buy used to find the post
		post.Title = postRequest.Title
		post.Content = postRequest.Content

		// Save the post
		err = s.PostsService.UpdatePosts(post, author)
		if err != nil {
			if err == internal.ErrPostNotFound {
				writeJSONError(w, "post not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Status accepted
		w.WriteHeader(http.StatusAccepted)

	}
}

func (s *Server) DeletePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			writeJSONError(w, "error getting author from context", http.StatusInternalServerError)
			return
		}

		// Get the post ID from the URL
		id, ok := mux.Vars(r)["id"]
		if !ok {
			writeJSONError(w, "missing post id", http.StatusBadRequest)
			return
		}

		// Convert the post ID from string to int
		postID, err := strconv.Atoi(id)
		if err != nil {
			writeJSONError(w, "invalid post id", http.StatusBadRequest)
			return
		}

		// Delete the post
		err = s.PostsService.DeletePosts(postID, author)
		if err != nil {
			if err == internal.ErrPostNotFound {
				writeJSONError(w, "post not found", http.StatusNotFound)
				return
			}
			writeJSONError(w, "error deleting post", http.StatusInternalServerError)
			return
		}

		// Status accepted
		w.WriteHeader(http.StatusAccepted)

	}
}
