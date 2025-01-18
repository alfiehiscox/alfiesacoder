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
	Author      string
	Date        string
	Publish     bool
}

type ArticleService struct {
	initialised bool

	Context     context.Context
	Logger      *log.Logger
	ContextPath string
	Articles    map[string]Article
	Markdown    goldmark.Markdown
}

func NewArticleService(
	ctx context.Context,
	path string,
	logger *log.Logger,
	md goldmark.Markdown,
) *ArticleService {
	return &ArticleService{
		Logger:      logger,
		Context:     ctx,
		ContextPath: path,
		Articles:    make(map[string]Article),
		Markdown:    md,
	}
}

func (as *ArticleService) Init() error {
	if as.initialised {
		return errors.New("ArticleService is already initialised")
	}

	articleEntries, err := os.ReadDir(as.ContextPath)
	if err != nil {
		as.initialised = false
		return err
	}

	for _, entry := range articleEntries {

		data, err := os.ReadFile(path.Join(as.ContextPath, entry.Name()))
		if err != nil {
			as.initialised = false
			return err
		}

		article, err := ArticleExtractionFunction(data, as.Markdown, entry.Name())
		if err != nil {
			return err
		}

		as.Articles[article.URL] = article

	}

	as.initialised = true
	return nil
}

func (as *ArticleService) GetArticleByURL(url string) (a Article, ok bool) {
	if !as.initialised {
		return
	}
	a, ok = as.Articles[url]
	return
}

func (as *ArticleService) GetPublishedArticles() (articles []Article) {
	if !as.initialised {
		return
	}

	for _, article := range as.Articles {
		if article.Publish {
			articles = append(articles, article)
		}
	}

	return articles
}

func ArticleExtractionFunction(
	data []byte,
	md goldmark.Markdown,
	filename string,
) (
	article Article,
	err error,
) {
	parser_context := parser.NewContext()
	var content bytes.Buffer
	if err := md.Convert(data, &content, parser.WithContext(parser_context)); err != nil {
		return article, err
	}
	metaData := meta.Get(parser_context)

	metaURL := metaData["URL"]
	var url string
	if u, ok := metaURL.(string); ok {
		url = u
	} else {
		return Article{}, errors.New("URL most be defined in the meta of all content files.")
	}

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

	metaAuthor := metaData["Author"]
	var author string
	if a, ok := metaAuthor.(string); ok {
		author = a
	}

	// Dates are stored in documents as yyyy-MM-dd
	metaDate := metaData["Date"]
	var date string
	if d, ok := metaDate.(string); ok {
		date = d
	}

	metaPublish := metaData["Publish"]
	publish := false
	if p, ok := metaPublish.(bool); ok {
		publish = p
	}

	article.Filename = filename
	article.Title = title
	article.Description = description
	article.Content = content.String()
	article.Author = author
	article.Date = date
	article.URL = url
	article.Publish = publish

	return article, nil
}
