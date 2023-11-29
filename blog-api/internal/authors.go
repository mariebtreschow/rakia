package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

const FILEPATH = "./resources/blog_data.json"

type Author struct {
	Author   string `json:"author"`
	Password string `json:"password"`
}

type AuthorData struct {
	Authors []Author `json:"posts"` // from json file
}

type AuthorService struct {
	authors []Author
	logger  *zerolog.Logger
}

// convertAuthorToPassword replaces "Author" with "password" in the given string.
func convertAuthorToPassword(author string) string {
	return strings.Replace(author, "Author ", "password", 1)
}

func getAuthors() ([]Author, error) {
	// Open the JSON file
	file, err := os.Open(FILEPATH)
	if err != nil {
		return nil, fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	// Decode the JSON file
	var a AuthorData
	if err := json.NewDecoder(file).Decode(&a); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	var authors []Author
	for _, author := range a.Authors {
		author.Password = convertAuthorToPassword(author.Author)
		authors = append(authors, author)

	}
	return authors, nil
}

func NewAuthorService(logger *zerolog.Logger) (*AuthorService, error) {
	// Add the authors from the json file in the resources folder to the authors slice
	authors, err := getAuthors()
	if err != nil {
		return nil, fmt.Errorf("error getting authors: %w", err)
	}
	return &AuthorService{
		authors: authors,
		logger:  logger,
	}, nil
}

// ValidAuthor returns the author id if the username and password are valid
func (a *AuthorService) ValidAuthor(username string, password string) (bool, error) {
	// Since its only 100 authors, we can just loop through them, but if we had a lot of authors
	// We would want to use a map to store the authors and their passwords
	for _, author := range a.authors {
		if author.Author == username && author.Password == password {
			return true, nil
		}
	}
	return false, nil
}
