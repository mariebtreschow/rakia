package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog"
)

type contextKey string

const (
	// ContextAuthor is the key for the author data in the request context
	ContextAuthor contextKey = "author"
)

func Middleware(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().Msg("validate author with JWT token")

			// Extract the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				writeJSONError(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := bearerToken[1]

			// Parse the token
			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				writeJSONError(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// If the token is valid, set the author in the context
			ctx := context.WithValue(r.Context(), ContextAuthor, claims.Username)

			// Call the next handler, with the new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
