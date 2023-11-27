package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var (
	ErrPostNotFound = fmt.Errorf("post not found")
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type PostData struct {
	Posts []Post `json:"posts"` // from json file
}

type AuthorPostsMap map[string][]Post

// Store the blogposts from the json file per author
type Persistence struct {
	Posts  AuthorPostsMap
	logger *zerolog.Logger
}

func NewPersistance(logger *zerolog.Logger) (*Persistence, error) {
	blogPosts, err := getPostsFromFile()
	if err != nil {
		logger.Fatal().Err(err).Msg("error getting blog posts from file")
		return nil, err
	}
	return &Persistence{
		Posts:  blogPosts,
		logger: logger,
	}, nil
}

// seed the blogposts slice with data from the json file in resources/blog_data.json
func getPostsFromFile() (map[string][]Post, error) {
	// Open the JSON file
	jsonFile, err := os.Open(FILEPATH)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// Read the file content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var data PostData
	// Unmarshal the JSON data into the posts slice
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return nil, err
	}

	// Define a variable to hold the data from the JSON file
	authorPosts := make(AuthorPostsMap)

	for _, post := range data.Posts {
		fmt.Println(post)
		authorPosts[post.Author] = append(authorPosts[post.Author], post)
	}

	return authorPosts, nil
}

// CreatePosts creates a new blogpost
func (p *Persistence) CreatePosts(post Post) error {
	// Add the post to the posts slice

	return nil
}

// Get all posts for the author
func (p *Persistence) GetAllPosts(author string) ([]*Post, error) {
	// Get all posts for the author
	posts := p.Posts[author]

	// Create a slice of pointers to the posts
	postPointers := make([]*Post, len(posts))
	for i := range posts {
		postPointers[i] = &posts[i]
	}

	return postPointers, nil
}

// GetPosts gets a blogpost by id
func (p *Persistence) GetPosts(id int, author string) (*Post, error) {
	// Get post for the author
	posts := p.Posts[author]
	for i := range posts {
		if posts[i].ID == id {
			return &posts[i], nil
		}
	}
	return nil, ErrPostNotFound
}

// UpdatePosts updates a blogpost
func (p *Persistence) UpdatePosts(post Post) error {
	// Update the post in the posts slice
	posts := p.Posts[post.Author]
	for i := range posts {
		if posts[i].ID == post.ID {
			posts[i] = post
		}
	}
	return ErrPostNotFound

}

// DeletePosts deletes a blogpost
func (p *Persistence) DeletePosts(id int, author string) error {
	// Delete the post from the posts slice

	posts := p.Posts[author]
	for i := range posts {
		if posts[i].ID == id {
			// Delete the post from the posts slice
			posts = append(posts[:i], posts[i+1:]...)
			return nil
		}
	}
	return ErrPostNotFound
}
