package main

import (
	"log"
	"net/http"

	"github.com/alfiehiscox/alfiesacoder/services"
	"github.com/alfiehiscox/alfiesacoder/templates"
)

func handleIndex(
	log *log.Logger,
	projectService *services.ProjectService,
	articleService *services.ArticleService,
) http.Handler {
	projects := projectService.GetPublishedProjects()
	articles := articleService.GetPublishedArticles()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
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
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			article, ok := articleService.GetArticleByURL(r.RequestURI)
			if !ok {
				log.Printf("Error: Could not article with url %s\n", r.RequestURI)
				templates.NotFound().Render(r.Context(), w)
				return
			}
			article_page := templates.Article(article)
			article_page.Render(r.Context(), w)
		},
	)
}
