package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"
	"unicode/utf8"

	"github.com/rs/zerolog"
)

var (
	ErrPostNotFound      = fmt.Errorf("post not found")
	ErrContentEmpty      = fmt.Errorf("content must not be empty")
	ErrContentInvalid    = fmt.Errorf("content must not be longer than 1600 characters or shorter than 100")
	ErrContentEncoding   = fmt.Errorf("content must be valid UTF-8")
	ErrUniqueTitle       = fmt.Errorf("title must be unique")
	ErrTitleEmpty        = fmt.Errorf("title must not be empty")
	ErrTitleInvalid      = fmt.Errorf("title must not be longer than 50 characters or shorter than 5 characters")
	ErrTitleInvalidChars = fmt.Errorf("title must not contain special characters")
	ErrAuthorEmpty       = fmt.Errorf("author must not be empty")
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

type AuthorLastIDMap map[string]int

// Store the blogposts from the json file per author
type Persistence struct {
	Posts  AuthorPostsMap
	LastID AuthorLastIDMap
	mutex  sync.Mutex // Protects access to lastID and Posts
	logger *zerolog.Logger
}

func NewPersistance(ap *AuthorPostsMap, al *AuthorLastIDMap, logger *zerolog.Logger) (*Persistence, error) {
	return &Persistence{
		Posts:  *ap,
		LastID: *al,
		logger: logger,
	}, nil
}

func (p *Persistence) Seed() {
	// Add the authors from the json file in the resources folder to the authors slice
	authors, err := getAuthors()
	if err != nil {
		p.logger.Fatal().Err(err).Msg("error getting authors from file")
		return
	}
	for _, author := range authors {
		// Add the posts from the json file in the resources folder to the posts slice
		posts, err := getPostsFromFile()

		if err != nil {
			p.logger.Fatal().Err(err).Msg("error getting posts from file")
			return
		}
		p.Posts[author.Author] = posts[author.Author]
		p.LastID[author.Author] = posts[author.Author][len(posts[author.Author])-1].ID
	}
}

// Add the blogposts from the json file in the resources folder to the posts slice
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

	// Add the posts to the posts slice
	for _, post := range data.Posts {
		if validateContent(post.Content) != nil {
			return nil, err
		}
		if validateTitle(post.Title) != nil {
			return nil, err
		}
		if validateAuthor(post.Author) != nil {
			return nil, err
		}
		authorPosts[post.Author] = append(authorPosts[post.Author], post)
	}
	return authorPosts, nil
}

func validateTitle(title string) error {
	if title == "" {
		return ErrTitleEmpty
	}

	if len(title) > 60 || len(title) < 5 {
		return ErrTitleInvalid
	}

	// Check for unwanted special characters
	// This regular expression allows letters, numbers, spaces, hyphens, and underscores
	matched, err := regexp.MatchString("^[a-zA-Z0-9\\-\\_\\s]+$", title)
	if err != nil {
		return err
	}
	if !matched {
		return ErrTitleInvalidChars
	}
	return nil
}

func validateContent(content string) error {
	if content == "" {
		return ErrContentEmpty
	}
	if len(content) > 1600 || len(content) < 100 {
		return ErrContentInvalid
	}
	// Check for proper UTF-8 encoding and make sure its letters not just symbols
	if !utf8.ValidString(content) {
		return ErrContentEncoding
	}
	// Disallow strings with excessive special characters
	specialCharPattern := regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;"'<,>.?/\\|~` + "`" + `]`)
	if len(specialCharPattern.FindAllString(content, -1)) > len(content)/2 { // Example condition
		return ErrContentInvalid
	}

	return nil
}

func validateAuthor(author string) error {
	if author == "" {
		return ErrAuthorEmpty
	}
	return nil
}

// CreatePosts creates a new blogpost
func (p *Persistence) CreatePosts(post Post) error {
	// mutex.Lock() and mutex.Unlock() ensure that only one goroutine can access the map at a time
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Check if the title is unique
	for _, existingPost := range p.Posts[post.Author] {
		if existingPost.Title == post.Title {
			return ErrUniqueTitle
		}
	}

	// Validate the post
	errTitle := validateTitle(post.Title)
	if errTitle != nil {
		return errTitle
	}

	errContent := validateContent(post.Content)
	if errContent != nil {
		return errContent
	}

	errAuthor := validateAuthor(post.Author)
	if errAuthor != nil {
		return errAuthor
	}

	// Add ID, must be unique
	post.ID = p.LastID[post.Author] + 1
	// Increment the lastID
	p.LastID[post.Author] = post.ID

	// Add the post to the posts slice
	p.Posts[post.Author] = append(p.Posts[post.Author], post)

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
	// Get all posts for the author
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
	// mutex.Lock() and mutex.Unlock() ensure that only one goroutine can access the map at a time
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Validate the post
	errTitle := validateTitle(post.Title)
	if errTitle != nil {
		return errTitle
	}

	errContent := validateContent(post.Content)
	if errContent != nil {
		return errContent
	}

	errAuthor := validateAuthor(post.Author)
	if errAuthor != nil {
		return errAuthor
	}
	posts := p.Posts[post.Author]
	for i := range posts {
		// Make sure the post belongs to the author
		if posts[i].ID == post.ID && posts[i].Author == post.Author {
			posts[i] = post
			return nil
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
			posts = append(posts[:i], posts[i+1:]...)
			return nil
		}
	}
	return ErrPostNotFound
}
