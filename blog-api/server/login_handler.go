package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"rakia.ai/blog-api/v2/internal"
)

// JWT Secret Key (should be kept secret and not hardcoded, this is just for demo purposes)
var jwtKey = []byte("my_secret_key")

// JWT Expiration Time (30 minutes)
var expirationTime = time.Now().Add(time.Minute * 30)

// Claims struct for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (s *Server) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var credentials internal.Author

		// Decode the incoming JSON payload
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			writeJSONError(w, "invalid request payload", http.StatusBadRequest)
			return
		}

		// Validate the author's credentials
		valid, err := s.AuthorsService.ValidAuthor(credentials.Author, credentials.Password)
		if err != nil {
			writeJSONError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		if !valid {
			writeJSONError(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		// Create the JWT claims, which includes the username and expiry time
		claims := &Claims{
			Username: credentials.Author,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
				Subject:   credentials.Author,
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Create the JWT string
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			writeJSONError(w, "failed to create token", http.StatusInternalServerError)
			return
		}

		// Construct the response object
		response := LoginResponse{
			Token: tokenString,
		}

		// Marshal the response object to JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			writeJSONError(w, "failed to marshal response", http.StatusInternalServerError)
			return
		}

		// Set the content-type header to json
		w.Header().Set("Content-Type", "application/json")

		// Send the response
		w.Write(jsonResponse)
	}
}
