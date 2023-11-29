package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
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

// NewPersistance creates a new blogposts service
func NewPersistance(ap *AuthorPostsMap, al *AuthorLastIDMap, logger *zerolog.Logger) (*Persistence, error) {
	return &Persistence{
		Posts:  *ap,
		LastID: *al,
		logger: logger,
	}, nil
}

// Seed the blogposts
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
		if len(posts[author.Author]) == 0 {
			p.logger.Fatal().Err(err).Msg("error getting posts from file")
			return
		}
		if err != nil {
			p.logger.Fatal().Err(err).Msg("error getting posts from file")
			return
		}
		// Add the posts to the posts slice
		p.Posts[author.Author] = posts[author.Author]
		// Add the lastID to the lastID slice to keep track of the last ID
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
