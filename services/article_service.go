package services

import (
	"bytes"
	"errors"

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

func (a Article) GetURL() string {
	return a.URL
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

	article.Filename = filename
	article.Title = title
	article.Description = description
	article.Content = content.String()
	article.URL = url

	return article, nil
}
