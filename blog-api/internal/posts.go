package internal

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type PostData struct {
	Data []Post `json:"posts"` // from json file
}

// store in memory json for a blogpost
type Persistence struct {
	Posts  PostData
	logger *zerolog.Logger
}

func NewPersistance(logger *zerolog.Logger) (*Persistence, error) {
	blogPosts, err := getPostsFromFile()
	if err != nil {
		logger.Fatal().Err(err).Msg("error getting blog posts from file")
		return nil, err
	}
	return &Persistence{
		Posts:  *blogPosts,
		logger: logger,
	}, nil
}

// seed the blogposts slice with data from the json file in resources/blog_data.json
func getPostsFromFile() (*PostData, error) {
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

	// Define a variable to hold the data from the JSON file
	var posts PostData

	// Unmarshal the JSON data into the posts slice
	if err := json.Unmarshal(byteValue, &posts); err != nil {
		return nil, err
	}
	return &posts, nil
}

// CreatePosts creates a new blogpost
func (p *Persistence) CreatePosts(ctx context.Context, post Post) error {
	// Add the post to the posts slice
	p.Posts.Data = append(p.Posts.Data, post)

	return nil
}

// Get all posts for the author
func (p *Persistence) GetAllPosts(ctx context.Context, author string) ([]Post, error) {
	// Get all posts for the author
	var posts []Post

	// Ok to use range here because the number of posts will be small
	for _, post := range p.Posts.Data {
		if post.Author == author {
			posts = append(posts, post)
		}
	}
	return posts, nil
}
