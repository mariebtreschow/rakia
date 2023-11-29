package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"rakia.ai/blog-api/v2/internal"
)

var ErrInvalidRequest = "unable to process request due to invalid information"

type PostCreate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type PostUpdate struct {
	ID      int    `json:"id"`
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

// GetAllPostsHandler gets all posts
func (s *Server) GetAllPostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info().Msg("get all posts handler called")
		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			s.Logger.Error().Msg("error getting author from context")
			writeJSONError(w, ErrInvalidRequest, http.StatusBadRequest)
			return
		}

		// Get all posts for the author
		posts, err := s.PostsService.GetAllPosts(author)
		if err != nil {
			s.Logger.Error().Err(err).Msg("error getting posts")
			writeJSONError(w, "error getting posts", http.StatusInternalServerError)
			return
		}
		// Check if there are any posts
		if len(posts) == 0 {
			s.Logger.Error().Msg("no posts found")
			writeJSONError(w, "no posts found", http.StatusNotFound)
			return
		}

		// JSON encode the posts
		jsonResponse, err := json.Marshal(posts)
		if err != nil {
			s.Logger.Error().Err(err).Msg("error marshalling posts")
			writeJSONError(w, "error getting posts", http.StatusInternalServerError)
			return
		}

		// Set the content-type header to json
		w.Header().Set("Content-Type", "application/json")

		// Send the response
		w.Write(jsonResponse)

	}
}

// GetPostsHandler gets a post
func (s *Server) GetPostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info().Msg("get posts handler called")
		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			s.Logger.Error().Msg("error getting author from context")
			writeJSONError(w, ErrInvalidRequest, http.StatusBadRequest)
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

		// JSON encode the post
		jsonResponse, err := json.Marshal(post)
		if err != nil {
			s.Logger.Error().Err(err).Msg("error marshalling post")
			writeJSONError(w, "error marshalling post", http.StatusInternalServerError)
			return
		}

		// Set the content-type header to json
		w.Header().Set("Content-Type", "application/json")

		// Send the response
		w.Write(jsonResponse)

	}
}

// CreatePostsHandler creates a new post
func (s *Server) CreatePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info().Msg("create posts handler called")
		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			s.Logger.Error().Msg("error getting author from context")
			writeJSONError(w, ErrInvalidRequest, http.StatusBadRequest)
			return
		}

		// Validate incoming post
		var postRequest PostCreate
		err := json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			s.Logger.Error().Err(err).Msg("invalid request payload")
			writeJSONError(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		// Create the post
		post := internal.Post{
			Title:   postRequest.Title,
			Content: postRequest.Content,
			Author:  postRequest.Author,
		}
		// Check if the author is empty
		if post.Author == "" {
			s.Logger.Error().Msg("author must not be empty")
			writeJSONError(w, "author must not be empty", http.StatusBadRequest)
			return
		}
		// Check if the author in the request matches the author in token
		if post.Author != author && author != "admin" {
			s.Logger.Error().Msg("mismatching authors in request and url")
			writeJSONError(w, "not allowed to create posts for another author", http.StatusBadRequest)
			return
		}

		// Save the post
		err = s.PostsService.CreatePosts(post)
		if err != nil {
			s.Logger.Error().Err(err).Msg("error creating post")
			// Handle validation errors
			if err == internal.ErrUniqueTitle || err == internal.ErrTitleEmpty || err == internal.ErrTitleInvalid || err == internal.ErrContentEmpty || err == internal.ErrContentInvalid || err == internal.ErrAuthorEmpty || err == internal.ErrContentEncoding || err == internal.ErrTitleInvalidChars {
				writeJSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSONError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Status accepted
		w.WriteHeader(http.StatusCreated)

	}
}

// UpdatePostsHandler updates a post
func (s *Server) UpdatePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info().Msg("update posts handler called")
		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			s.Logger.Error().Msg("error getting author from context")
			writeJSONError(w, ErrInvalidRequest, http.StatusBadRequest)
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

		// Validate incoming post
		var postRequest PostUpdate
		err = json.NewDecoder(r.Body).Decode(&postRequest)
		if err != nil {
			s.Logger.Error().Err(err).Msg("invalid request payload")
			writeJSONError(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		// Check if the post ID in the request matches the post ID in the URL
		if postID != postRequest.ID {
			s.Logger.Error().Msg("mismatching ids in request and url")
			writeJSONError(w, "mismatching ids in request and url", http.StatusBadRequest)
			return
		}
		// Check if the author is empty
		if postRequest.Author == "" {
			s.Logger.Error().Msg("author must not be empty")
			writeJSONError(w, "author must not be empty", http.StatusBadRequest)
			return
		}
		// Check if the author in the request matches the author in token
		if postRequest.Author != author {
			s.Logger.Error().Msg("mismatching authors in request and url")
			writeJSONError(w, "not allowed to update posts for another author", http.StatusBadRequest)
			return
		}
		var post internal.Post

		// Update the post
		post.ID = postRequest.ID // Not being update buy used to find the post
		post.Title = postRequest.Title
		post.Content = postRequest.Content
		post.Author = postRequest.Author

		// Save the updated post
		err = s.PostsService.UpdatePosts(post)
		if err != nil {
			s.Logger.Error().Err(err).Msg("error updating post")
			if err == internal.ErrUniqueTitle || err == internal.ErrTitleEmpty || err == internal.ErrTitleInvalid || err == internal.ErrContentEmpty || err == internal.ErrContentInvalid || err == internal.ErrAuthorEmpty || err == internal.ErrContentEncoding || err == internal.ErrTitleInvalidChars {
				writeJSONError(w, err.Error(), http.StatusBadRequest)
				return
			}
			if err == internal.ErrPostNotFound {
				s.Logger.Error().Err(err).Msg("post not found")
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

// DeletePostsHandler deletes a post
func (s *Server) DeletePostsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info().Msg("delete posts handler called")
		// Get the context from the request
		author, ok := r.Context().Value(ContextAuthor).(string)
		if !ok {
			s.Logger.Error().Msg("error getting author from context")
			writeJSONError(w, ErrInvalidRequest, http.StatusBadRequest)
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

		// Delete the post
		err = s.PostsService.DeletePosts(postID, author)
		if err != nil {
			if err == internal.ErrPostNotFound {
				s.Logger.Error().Err(err).Msg("post not found")
				writeJSONError(w, "post not found", http.StatusNotFound)
				return
			}
			s.Logger.Error().Err(err).Msg("error deleting post")
			writeJSONError(w, "error deleting post", http.StatusInternalServerError)
			return
		}

		// Status accepted
		w.WriteHeader(http.StatusAccepted)

	}
}
