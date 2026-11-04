// Command server runs the calculator API and, when STATIC_DIR is set, serves the
// built frontend from the same port.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fullstack-calculator/backend/internal/api"
)

const (
	defaultPort = "8080"

	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 60 * time.Second

	// shutdownTimeout stays under the 10 seconds docker stop allows before it
	// kills the process.
	shutdownTimeout = 8 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      api.NewHandler(os.Getenv("STATIC_DIR")),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", server.Addr)
		serveErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	log.Print("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}
