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

type AuthorPassword map[string]string

type AuthorService struct {
	authors AuthorPassword
	logger  *zerolog.Logger
}

// NewAuthorService creates a new author service
func NewAuthorService(a *AuthorPassword, logger *zerolog.Logger) (*AuthorService, error) {
	// Add the authors from the json file in the resources folder to the authors slice
	return &AuthorService{
		authors: *a,
		logger:  logger,
	}, nil
}

// convertAuthorToPassword replaces "Author" with "password" in the given string.
func convertAuthorToPassword(author string) string {
	return strings.Replace(author, "Author ", "password", 1)
}

// Seed adds the authors from the json file in the resources folder to the authors slice
func (a *AuthorService) Seed() error {
	// Add the authors from the json file in the resources folder to the authors slice
	// Open the JSON file
	file, err := os.Open(FILEPATH)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %w", err)
	}
	defer file.Close()

	// Decode the JSON file
	var data AuthorData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}

	// Add the authors to the authors slice
	for _, author := range data.Authors {
		a.authors[author.Author] = convertAuthorToPassword(author.Author)
	}

	// Add admin user
	a.authors["admin"] = "admin"
	return nil
}

// ValidAuthor returns the author id if the username and password are valid
func (a *AuthorService) ValidAuthor(username string, password string) (bool, error) {
	if val, ok := (a.authors)[username]; ok {
		if val == password {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}
