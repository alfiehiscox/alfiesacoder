package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/alfiehiscox/alfiesacoder/services"
	"github.com/alfiehiscox/alfiesacoder/templates"
)

func handleIndex(
	log *log.Logger,
	projectService *services.ProjectService,
	articleService *services.ArticleService,
	statsService *services.ArticleStatsService,
) http.Handler {
	projects := projectService.GetPublishedProjects()
	articles := articleService.PublishedArticles
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			statsService.IncrementView("/")
			index := templates.Index(projects, articles)
			index.Render(r.Context(), w)
		},
	)
}

func handleStatic() http.Handler {
	return http.FileServer(http.Dir("./static"))
}

func handleArticles(
	log *log.Logger,
	articleService *services.ArticleService,
	statsService *services.ArticleStatsService,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			article, ok := articleService.GetArticleByURL(r.RequestURI)
			if !ok {
				log.Printf("Error: Could not article with url %s\n", r.RequestURI)
				templates.NotFound().Render(r.Context(), w)
				return
			}
			statsService.IncrementView(article.URL)
			article_page := templates.Article(article)
			article_page.Render(r.Context(), w)
		},
	)
}

func handleArticleArchive(
	log *log.Logger,
	articleService *services.ArticleService,
	statsService *services.ArticleStatsService,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			pageString := r.PathValue("page")
			page, err := strconv.Atoi(pageString)
			if err != nil {
				log.Printf("Error: Could not get archive page: %s\n", pageString)
				templates.NotFound().Render(r.Context(), w)
				return
			}

			if page > 0 {
				page = page - 1
			}

			statsService.IncrementView(r.RequestURI)

			archive := articleService.GetPublishedArticlesByPage(page)
			archive_page := templates.ArticleArchive(page+1, articleService.MaxPages, archive)
			archive_page.Render(r.Context(), w)

		},
	)
}

func handleArticleViews(
	log *log.Logger,
	statsService *services.ArticleStatsService,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := statsService.WriteStats(w); err != nil {
				log.Printf("ERROR: Could not serve article views: %s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "505 - Internal Server Error")
			}
		},
	)
}
