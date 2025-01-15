package services

import (
	"log"
)

type Article struct{}

type ArticleService struct {
	logger       *log.Logger
	context_path string
	initialised  bool
	articles     map[string]Article
}

func NewArticleService(path string, logger *log.Logger) *ArticleService {
	return &ArticleService{
		logger:       logger,
		context_path: path,
	}
}

// Graphs and Caches all articles content pages
func (as *ArticleService) Init() error { return nil }

func (as *ArticleService) GetArticles() (articles []Article) {
	for _, article := range as.articles {
		articles = append(articles, article)
	}
	return articles
}
