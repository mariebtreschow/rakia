package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/rs/zerolog"
)

var (
	ErrPostNotFound           = fmt.Errorf("post not found")
	ErrContentEmpty           = fmt.Errorf("content must not be empty")
	ErrContentInvalid         = fmt.Errorf("content must not be longer than 1600 characters or shorter than 100 and cannot have too many special characters")
	ErrContentEncoding        = fmt.Errorf("content must be valid UTF-8")
	ErrUniqueTitle            = fmt.Errorf("title must be unique")
	ErrTitleEmpty             = fmt.Errorf("title must not be empty")
	ErrTitleInvalid           = fmt.Errorf("title must not be longer than 60 characters or shorter than 5 characters")
	ErrTitleInvalidChars      = fmt.Errorf("title must not contain too many special characters")
	ErrAuthorEmpty            = fmt.Errorf("author must not be empty")
	ErrTitleFormat            = fmt.Errorf("title must not have excessive whitespace or multiple consecutive spaces")
	ErrTitleSpammy            = fmt.Errorf("title must not contain spammy patterns or phrases")
	ErrTitleCapitalization    = fmt.Errorf("title must follow capitalization rules")
	ErrContentConsecutiveChar = fmt.Errorf("content must not have excessive consecutive identical characters")
	ErrAuthorNotFound         = fmt.Errorf("author not found")
	ErrAuthorNameInvalid      = fmt.Errorf("author name must not be longer than 70 characters or shorter than 2 characters")
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

type AuthorPostsMap map[string]map[int]Post

// Store the blogposts from the json file per author
type PostService struct {
	Posts  map[string]map[int]Post
	LastID int
	mutex  sync.Mutex // Protects access to lastID and Posts
	logger *zerolog.Logger
}

// NewPostsService creates a new blogposts service
func NewPostsService(ap *map[string]map[int]Post, logger *zerolog.Logger) (*PostService, error) {
	return &PostService{
		Posts:  *ap,
		LastID: 0,
		logger: logger,
	}, nil
}

// Add the blogposts from the json file in the resources folder to the posts slice
func (p *PostService) Seed() error {
	// Open the JSON file
	jsonFile, err := os.Open(FILEPATH)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	// Read the file content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	// Define a variable to hold the data from the JSON file
	var data PostData

	// Unmarshal the JSON data into the posts slice
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return err
	}

	// Add the posts to the posts slice
	for _, post := range data.Posts {
		if validateContent(post.Content) != nil {
			return err
		}
		if validateTitle(post.Title) != nil {
			return err
		}
		if validateAuthor(post.Author) != nil {
			return err
		}
		// If the author is not in the map, add it
		if _, ok := p.Posts[post.Author]; !ok {
			p.Posts[post.Author] = make(map[int]Post)
		}
		// Add to the map with the ID as the key
		p.Posts[post.Author][post.ID] = post
	}
	// Update the lastID
	p.LastID = len(data.Posts)
	return nil
}

// isCapitalizedProperly checks if the title follows the capitalization rules.
func isCapitalizedProperly(title string) bool {
	words := strings.Fields(title)
	for _, word := range words {
		// Discard check if its a number
		if _, err := fmt.Sscanf(word, "%f", new(float64)); err == nil {
			continue
		}
		// Check if the first letter is uppercase
		if len(word) > 1 && !unicode.IsUpper(rune(word[0])) {
			return false
		}
		// Check if the rest of the letters are lowercase
		if len(word) > 1 && unicode.IsUpper(rune(word[0])) && strings.ToLower(word[1:]) != word[1:] {
			return false
		}
	}
	return true
}

// validateTitle checks if the title is empty, too long, has too many special characters, has excessive whitespace or multiple consecutive spaces, or contains spammy patterns or phrases.
func validateTitle(title string) error {
	// Check for empty title
	if title == "" {
		return ErrTitleEmpty
	}

	// Check for title length
	if len(title) > 60 || len(title) < 5 {
		return ErrTitleInvalid
	}

	// Check for unwanted special characters
	// Disallow strings with excessive special characters
	specialCharPattern := regexp.MustCompile(`[!@#$%^&*()_+{}\[\]:;"'<,>.?/\\|~` + "`" + `]`)
	// Find all instances of the pattern in the content
	specialChars := specialCharPattern.FindAllString(title, -1)
	// Calculate the percentage of the title that is composed of special characters.
	specialCharPercentage := float64(len(specialChars)) / float64(len(title))
	if specialCharPercentage > 0.1 { // If more than 10% of the title is special characters.
		return ErrTitleInvalidChars
	}

	// Check for excessive whitespace or multiple consecutive spaces
	matched, err := regexp.MatchString("\\s{2,}", title)
	if err != nil {
		return err
	}
	if matched {
		return ErrTitleFormat
	}

	// Validate against common spammy patterns or phrases.
	spammyPatterns := []string{"buy now", "discount"} // Example patterns, extend as needed
	for _, pattern := range spammyPatterns {
		if strings.Contains(strings.ToLower(title), pattern) {
			return ErrTitleSpammy
		}
	}

	// Implement a check for word capitalization rules.
	if !isCapitalizedProperly(title) {
		return ErrTitleCapitalization
	}
	return nil
}

// validateContent checks if the content is empty or too long, and if it has too many special characters
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
	// Find all instances of the pattern in the content
	specialChars := specialCharPattern.FindAllString(content, -1)

	// Check if the number of special characters exceeds a tenth of the length of the content
	if float64(len(specialChars)) > float64(len(content))/10 {
		// If condition is true, return the ErrContentInvalid
		return ErrContentInvalid
	}

	// Check for excessive consecutive identical characters.
	for _, r := range content {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			escapedR := regexp.QuoteMeta(string(r))

			// Match 4 or more consecutive identical characters
			pattern := regexp.MustCompile(escapedR + `{4,}`)
			if pattern.FindString(content) != "" {
				return ErrContentConsecutiveChar
			}
		}
	}

	return nil
}

// validateAuthor checks if the author is empty
func validateAuthor(author string) error {
	// Check for empty author
	if author == "" {
		return ErrAuthorEmpty
	}
	if len(author) > 70 || len(author) < 2 {
		return ErrAuthorNameInvalid
	}
	return nil
}

// CreatePosts creates a new blogpost
func (p *PostService) CreatePosts(post Post, author string) error {
	// mutex.Lock() and mutex.Unlock() ensure that only one goroutine can access the map at a time
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// If admin is the author, add any posts for any author
	if author == "admin" {
		if _, ok := p.Posts[post.Author]; !ok {
			p.Posts[post.Author] = make(map[int]Post)
		}
	} else {
		// Make sure the author is in the map
		if _, ok := p.Posts[post.Author]; !ok {
			return ErrAuthorNotFound
		}
	}
	// Check if the title is unique for the author
	for _, existingPost := range p.Posts[post.Author] {
		if existingPost.Title == post.Title {
			return ErrUniqueTitle
		}
	}

	// Validation
	if err := validateTitle(post.Title); err != nil {
		return err
	}
	if err := validateContent(post.Content); err != nil {
		return err
	}
	if err := validateAuthor(post.Author); err != nil {
		return err
	}
	// Add ID, must be unique
	post.ID = p.LastID + 1
	// Increment the lastID, so the next post will have a unique ID
	p.LastID = post.ID

	// Add the post
	p.Posts[post.Author][post.ID] = post

	return nil
}

// Get all posts for the author
func (p *PostService) GetAllPosts(author string) ([]*Post, error) {
	// Create a slice of pointers to the posts
	var result []*Post

	// If admin is the author, return all posts
	if author == "admin" {
		for _, posts := range p.Posts {
			// Add all posts to the postPointers slice
			for id := range posts {
				p := posts[id]
				result = append(result, &p)
			}
		}
		// Order the posts by ID
		sort.Slice(result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		})
		return result, nil
	}

	// Make sure the author is in the map
	posts, ok := p.Posts[author]
	if !ok {
		return nil, ErrAuthorNotFound
	}

	// If the author is not admin, return only the posts for that author
	for id := range posts {
		p := posts[id]
		result = append(result, &p)
	}

	return result, nil
}

// GetPosts gets a blogpost by id
func (p *PostService) GetPostByID(id int, author string) (*Post, error) {
	// If admin is the author, return all posts
	if author == "admin" {
		// If admin is the author, return any posts of id that exists
		for _, authorPosts := range p.Posts {
			if post, exists := authorPosts[id]; exists {
				result := post
				return &result, nil
			}
		}

	} else {
		// Make sure the author is in the map
		posts, ok := p.Posts[author]
		if !ok {
			return nil, ErrAuthorNotFound
		}

		// Get the post with the matching ID and that it exists
		if post, ok := posts[id]; ok {
			result := &post
			return result, nil
		}
	}

	// If the post is not found, return ErrPostNotFound
	return nil, ErrPostNotFound
}

// UpdatePosts updates a blogpost
func (p *PostService) UpdatePosts(post Post, author string) error {
	// mutex.Lock() and mutex.Unlock() ensure that only one goroutine can access the map at a time
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Validation
	if err := validateTitle(post.Title); err != nil {
		return err
	}
	if err := validateContent(post.Content); err != nil {
		return err
	}
	if err := validateAuthor(post.Author); err != nil {
		return err
	}

	// If admin is the author, update any posts
	if author == "admin" {
		// Update any post if ID exists
		for _, posts := range p.Posts {
			if _, ok := posts[post.ID]; ok {
				posts[post.ID] = post
				return nil
			}
		}
	} else {
		// Make sure the author is in the map
		if _, ok := p.Posts[post.Author]; !ok {
			return ErrAuthorNotFound
		}
		// If the author is not admin, update only the posts for that author
		posts := p.Posts[post.Author]

		// make sure the post exists
		if _, ok := posts[post.ID]; !ok {
			return ErrPostNotFound
		}

		// Update the post if the author matches
		if posts[post.ID].Author == author {
			posts[post.ID] = post
		}

	}
	return ErrPostNotFound
}

// DeletePosts deletes a blogpost
func (p *PostService) DeletePosts(id int, author string) error {
	// Delete the post from the posts slice
	// If admin is the author, delete any posts
	if author == "admin" {
		for _, posts := range p.Posts {
			if _, ok := posts[id]; ok {
				delete(posts, id)
				return nil
			}
		}
	} else {
		// If the author is not admin, delete only the posts for that author
		if posts, ok := p.Posts[author]; ok {
			if _, ok := posts[id]; ok {
				delete(posts, id)
				return nil
			}
		}

	}
	return ErrPostNotFound
}
