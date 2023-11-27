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
		listenPort = fs.String("port", "8080", "Port to listen on")
	)

	fs.Parse(os.Args[1:])

	// Logger for the server
	logger := server.NewLogger()

	// Create new blog posts service
	logger.Info().Msg("seeding blog posts")
	posts, err := internal.NewPersistance(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error seeding blog posts")
	}

	// Create a new author service
	logger.Info().Msg("creating author service")
	authors, err := internal.NewAuthorService(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("error creating author service")
	}

	// Create a new mux router TODO: add why we use mux
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
		Handler:      s.Router, // Assuming s.Router is your mux router
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Info().Msgf("starting server on port: %s", *listenPort)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Err(err).Msg("ListenAndServe failed")
		}
	}()

	<-done // Wait for interrupt signal to gracefully shutdown the server

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info().Msg("shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Err(err).Msg("server shutdown failed")
	}
	logger.Info().Msg("server exited properly")

}
