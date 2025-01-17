package main

import (
	"log"
	"net/http"

	"github.com/alfiehiscox/alfiesacoder/services"
	"github.com/alfiehiscox/alfiesacoder/templates"
)

func handleIndex(
	log *log.Logger,
	projectService *services.ContentService[services.Project],
	articleService *services.ContentService[services.Article],
) http.Handler {
	projects := projectService.GetContent()
	articles := articleService.GetContent()
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
	articleService *services.ContentService[services.Article],
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			article, ok := articleService.GetContentByURL(r.RequestURI)
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

func handleProjects(
	log *log.Logger,
	projectService *services.ContentService[services.Project],
) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			project, ok := projectService.GetContentByURL(r.RequestURI)
			if !ok {
				log.Printf("Error: Could not project with url %s\n", r.RequestURI)
				templates.NotFound().Render(r.Context(), w)
				return
			}
			project_page := templates.Project(project)
			project_page.Render(r.Context(), w)
		},
	)
}
