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
	"path"
	"sync"

	"github.com/alfiehiscox/alfiesacoder/services"
	"github.com/joho/godotenv"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Getwd, os.Stderr, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	getwd func() (string, error),
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

	wd, err := getwd()
	if err != nil {
		return err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(meta.Meta),
	)

	projects := services.NewContentService[services.Project](
		ctx,
		path.Join(wd, "content", "projects"),
		log,
		markdown,
		services.ProjectExtractionFunction,
	)
	if err := projects.Init(); err != nil {
		return err
	}

	articles := services.NewContentService[services.Article](
		ctx,
		path.Join(wd, "content", "articles"),
		log,
		markdown,
		services.ArticleExtractionFunction,
	)
	if err := articles.Init(); err != nil {
		return err
	}

	srv := NewServer(log, projects, articles)
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
	projectService *services.ContentService[services.Project],
	articleService *services.ContentService[services.Article],
) http.Handler {

	mux := http.NewServeMux()

	addRoutes(
		mux,
		log,
		projectService,
		articleService,
	)

	var handler http.Handler = mux
	// middleware

	return handler
}

func addRoutes(
	mux *http.ServeMux,
	log *log.Logger,
	projectService *services.ContentService[services.Project],
	articleService *services.ContentService[services.Article],
) {
	mux.Handle("/", handleIndex(log, projectService, articleService))
	mux.Handle("/style.css", handleStyles())
	mux.Handle("/articles/{title}", handleArticles(log, articleService))
	mux.Handle("/projects/{name}", handleProjects(log, projectService))
}
