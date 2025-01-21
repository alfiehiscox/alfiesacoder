package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"sync"

	"github.com/alfiehiscox/alfiesacoder/services"
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

	port := getenv("PORT")
	if port == "" {
		port = "8000"
	}

	var articles_per_page int
	per_page := getenv("articlePerPage")
	if per_page == "" {
		articles_per_page = services.ArticlePerPageDefault
	} else {
		parsed, err := strconv.Atoi(per_page)
		articles_per_page = parsed
		if err != nil {
			return err
		}
	}

	wd, err := getwd()
	if err != nil {
		return err
	}

	markdown := goldmark.New(
		goldmark.WithExtensions(meta.Meta),
	)

	projects := services.NewProjectService(ctx, path.Join(wd, "content", "projects.json"), log)
	if err := projects.Init(); err != nil {
		return err
	}

	articles := services.NewArticleService(
		ctx,
		path.Join(wd, "content", "articles"),
		log,
		markdown,
		articles_per_page,
	)
	if err := articles.Init(); err != nil {
		return err
	}

	statsFile := path.Join(wd, "stats.json")
	statsService, err := services.NewArticleStatsService(ctx, log, statsFile)
	if err != nil {
		return err
	}

	srv := NewServer(log, projects, articles, statsService)
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}

	// Main Listener
	go func() {
		log.Printf("listening on Port=%s\n", port)
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
	projectService *services.ProjectService,
	articleService *services.ArticleService,
	statsService *services.ArticleStatsService,
) http.Handler {

	mux := http.NewServeMux()

	addRoutes(
		mux,
		log,
		projectService,
		articleService,
		statsService,
	)

	var handler http.Handler = mux
	// middleware

	return handler
}

func addRoutes(
	mux *http.ServeMux,
	log *log.Logger,
	projectService *services.ProjectService,
	articleService *services.ArticleService,
	statsService *services.ArticleStatsService,
) {
	static := http.FileServer(http.Dir("./static"))
	mux.Handle("/", handleIndex(log, projectService, articleService, statsService))
	mux.Handle("/static/", http.StripPrefix("/static/", static))
	mux.Handle("/archive/{page}", handleArticleArchive(log, articleService, statsService))
	mux.Handle("/articles/views", handleArticleViews(log, statsService))
	mux.Handle("/articles/{title}", handleArticles(log, articleService, statsService))
}
