package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/honeycombio/beeline-go/wrappers/hnygorilla"
	"rakia.ai/blog-api/v2/internal"
	"rakia.ai/blog-api/v2/server"
)

func main() {
	fs := flag.NewFlagSet("blog_api", flag.ExitOnError)

	var (
		listenPort = fs.String("port", "8080", "port to listen on")
		wait       = fs.Duration("graceful_timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	)

	fs.Parse(os.Args[1:])

	// Logger for the server
	logger := server.NewLogger()

	// Initialize the author posts map
	// This is a map of author names to a map of post IDs to posts
	p := make(map[string]map[int]internal.Post)
	p["Author 1"] = make(map[int]internal.Post)

	// Create new blog posts service
	posts, err := internal.NewPostsService(&p, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error creating blog posts service")
	}

	logger.Info().Msg("seeding blog posts")
	// Seed the blog posts
	posts.Seed()

	a := make(internal.AuthorPassword)
	// Create a new author service
	logger.Info().Msg("creating author service")
	authors, err := internal.NewAuthorService(&a, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error creating author service")
	}

	// Seed the authors
	authors.Seed()

	// Create a new mux router
	router := mux.NewRouter()

	// Create a new server
	logger.Info().Msg("creating server")
	s := server.NewServer(router, posts, authors, logger)

	s.Routes()
	s.Router.Use(hnygorilla.Middleware)

	// Create a new server
	srv := &http.Server{
		Addr:         ":" + *listenPort,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.Router,
	}

	logger.Info().Msgf("starting server on port: %s", *listenPort)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Err(err).Msg("ListenAndServe failed")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-c // Wait for interrupt signal to gracefully shutdown the server

	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()

	logger.Info().Msg("shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Err(err).Msg("server shutdown failed")
	}
	logger.Info().Msg("server exited properly")
	os.Exit(0)
}
