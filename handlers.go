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
	projects := projectService.GetProjects()
	articles := articleService.GetArticles()
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			index := templates.Index(projects, articles)
			index.Render(r.Context(), w)
		},
	)
}

func handleStyles() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./static/output.css")
		},
	)
}

func handleArticles(
	log *log.Logger,
	articleService *services.ArticleService,
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}

func handleProjects(
	log *log.Logger,
	projectService *services.ProjectService,
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		},
	)
}
