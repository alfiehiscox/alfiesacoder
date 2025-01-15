package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stderr, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log := log.New(stdout, "alfiesacoder: ", log.LstdFlags)

	if err := godotenv.Load(); err != nil {
		return err
	}

	host := getenv("host")
	if host == "" {
		return errors.New("host not set in environment variables")
	}
	port := getenv("port")
	if port == "" {
		return errors.New("port not set in environment variables")
	}

	srv := NewServer(log)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: srv,
	}

	// Main Listener
	go func() {
		log.Printf("listening on %s:%s\n", host, port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and server: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	// Shutdown Listener
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Printf("shutting down server\n")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithCancel(shutdownCtx)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}

func NewServer(
	log *log.Logger,
) http.Handler {

	mux := http.NewServeMux()
	addRoutes(mux, log)

	var handler http.Handler = mux
	// middleware

	return handler
}

func addRoutes(
	mux *http.ServeMux,
	log *log.Logger,
) {
	mux.Handle("/", handleIndex(log))
}
