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

// getAuthors returns the authors from the json file in the resources folder
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

// Seed adds the authors from the json file in the resources folder to the authors slice
func (a *AuthorService) Seed() {
	// Add the authors from the json file in the resources folder to the authors slice
	authors, err := getAuthors()
	if err != nil {
		a.logger.Fatal().Err(err).Msg("error getting authors from file")
		return
	}

	// Convert the authors to a map
	authorMap := make(map[string]string)
	for _, author := range authors {
		authorMap[author.Author] = author.Password
	}
	// Add admin user
	authorMap["admin"] = "admin"
	a.authors = authorMap
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
