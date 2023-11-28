package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockAuthorService is a mock version of AuthorsService
type MockAuthorService struct {
	validAuthor bool
}

// ValidAuthor mocks the ValidAuthor function of the AuthorsService
func (m *MockAuthorService) ValidAuthor(username, password string) (bool, error) {
	return m.validAuthor, nil
}

func TestLoginHandler(t *testing.T) {
	// Create a new instance of our server with a mock AuthorsService
	s := Server{
		AuthorsService: &MockAuthorService{validAuthor: true},
	}

	handler := s.LoginHandler()

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	loginData := `{"author":"testauthor","password":"password"}`
	req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(loginData))
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	// Check the response body is what we expect.
	expected := `{"token":"` // We expect a token to be returned, so the response should start with this string.
	assert.Contains(t, rr.Body.String(), expected, "handler returned unexpected body")

}
