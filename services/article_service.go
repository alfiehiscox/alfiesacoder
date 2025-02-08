package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"time"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

const ArticlePerPageDefault = 4

type Article struct {
	Filename    string
	Title       string
	Description string
	Content     string
	URL         string
	Author      string
	DateString  string
	Date        time.Time
	Publish     bool
}

type ArticleService struct {
	initialised bool

	Context           context.Context
	Logger            *log.Logger
	ContextPath       string
	Articles          map[string]Article
	PublishedArticles []Article
	Markdown          goldmark.Markdown
	ArticlePerPage    int
	MaxPages          int
}

func NewArticleService(
	ctx context.Context,
	path string,
	logger *log.Logger,
	md goldmark.Markdown,
	per_page int,
) *ArticleService {
	return &ArticleService{
		Logger:            logger,
		Context:           ctx,
		ContextPath:       path,
		Articles:          make(map[string]Article),
		PublishedArticles: []Article{},
		Markdown:          md,
		ArticlePerPage:    per_page,
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

	as.PublishedArticles = as.getPublishedArticles()
	as.MaxPages = len(as.PublishedArticles)/as.ArticlePerPage + 1

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

func (as *ArticleService) GetPublishedArticleByURL(url string) (a Article, ok bool) {
	if !as.initialised {
		return
	}

	a, ok = as.Articles[url]
	if ok {
		ok = a.Publish
	}

	return
}

func (as *ArticleService) getPublishedArticles() (articles []Article) {
	for _, article := range as.Articles {
		if article.Publish {
			articles = append(articles, article)
		}
	}

	// Orderd by Date. Most recent first
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Date.After(articles[j].Date)
	})

	return articles
}

// 'page' needs to be zero indexed
func (as *ArticleService) GetPublishedArticlesByPage(page int) (article []Article) {
	if !as.initialised {
		return
	}

	all_articles := as.PublishedArticles

	if page > as.MaxPages || page < 0 {
		start_index := 0
		end_index := min(as.ArticlePerPage, len(all_articles))
		return all_articles[start_index:end_index]
	} else {
		start_index := page * as.ArticlePerPage
		end_index := min(page*as.ArticlePerPage+as.ArticlePerPage, len(all_articles))
		return all_articles[start_index:end_index]
	}
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
	var dateString string
	var date time.Time
	if d, ok := metaDate.(string); ok {
		dateString = d
		parsed, err := time.Parse("2006-01-02", dateString)
		if err != nil {
			msg := fmt.Sprintf("Article meta date field could not be parsed: %s\n", err)
			return Article{}, errors.New(msg)
		}
		date = parsed
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
	article.DateString = dateString
	article.URL = url
	article.Publish = publish
	article.Date = date

	return article, nil
}
