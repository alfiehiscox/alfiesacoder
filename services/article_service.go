package services

import (
	"bytes"
	"context"
	"errors"
	"log"
	"os"
	"path"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type Article struct {
	Filename    string
	Title       string
	Description string
	Content     string
	URL         string
}

type ArticleService struct {
	Context     context.Context
	Logger      *log.Logger
	ContextPath string
	Initialised bool
	Articles    map[string]Article
	Markdown    goldmark.Markdown
	BaseURL     string
}

func NewArticleService(
	ctx context.Context,
	path string,
	logger *log.Logger,
	markdown goldmark.Markdown,
) *ArticleService {
	return &ArticleService{
		Logger:      logger,
		Context:     ctx,
		ContextPath: path,
		Initialised: false,
		Articles:    make(map[string]Article),
		Markdown:    markdown,
	}
}

// Graphs and Caches all articles content pages
func (as *ArticleService) Init() error {
	if as.Initialised {
		return errors.New("ArticleService is already initialised")
	}

	articleEntries, err := os.ReadDir(as.ContextPath)
	if err != nil {
		as.Initialised = false
		return err
	}

	for _, entry := range articleEntries {

		var article Article
		data, err := os.ReadFile(path.Join(as.ContextPath, entry.Name()))
		if err != nil {
			as.Initialised = false
			return err
		}

		parser_context := parser.NewContext()
		var content bytes.Buffer
		if err := as.Markdown.Convert(data, &content, parser.WithContext(parser_context)); err != nil {
			return err
		}
		metaData := meta.Get(parser_context)

		metaTitle := metaData["Title"]
		var title string
		if t, ok := metaTitle.(string); ok {
			title = t
		}

		metaDescription := metaData["Description"]
		var description string
		if d, ok := metaDescription.(string); ok {
			description = d
		}

		metaURL := metaData["URL"]
		var url string
		if u, ok := metaURL.(string); ok {
			url = u
		}

		article.Filename = entry.Name()
		article.Title = title
		article.Description = description
		article.Content = content.String()
		article.URL = url
		as.Articles[article.Filename] = article

	}

	as.Initialised = true
	return nil
}

func (as *ArticleService) GetArticles() (articles []Article) {
	if !as.Initialised {
		return
	}

	for _, article := range as.Articles {
		articles = append(articles, article)
	}

	return articles
}

func (as *ArticleService) GetArticleByURL(url string) (a Article, ok bool) {
	if !as.Initialised {
		return
	}

	for _, article := range as.Articles {
		if article.URL == url {
			return article, true
		}
	}

	return
}
