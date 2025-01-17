package services

import (
	"bytes"
	"errors"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type Project struct {
	Filename    string
	Name        string
	Description string
	Content     string
	URL         string
	Link        string
}

func (p Project) GetURL() string {
	return p.URL
}

func ProjectExtractionFunction(
	data []byte,
	md goldmark.Markdown,
	filename string,
) (
	project Project,
	err error,
) {
	parser_context := parser.NewContext()
	var content bytes.Buffer
	if err := md.Convert(data, &content, parser.WithContext(parser_context)); err != nil {
		return project, err
	}
	metaData := meta.Get(parser_context)

	metaURL := metaData["URL"]
	var url string
	if u, ok := metaURL.(string); ok {
		url = u
	} else {
		return Project{}, errors.New("URL most be defined in the meta of all content files.")
	}

	metaName := metaData["Name"]
	var name string
	if n, ok := metaName.(string); ok {
		name = n
	}

	metaDescription := metaData["Description"]
	var description string
	if d, ok := metaDescription.(string); ok {
		description = d
	}

	metaLink := metaData["Link"]
	var link string
	if l, ok := metaLink.(string); ok {
		link = l
	} else {
		link = "No Current Repository"
	}

	project.Filename = filename
	project.Name = name
	project.Description = description
	project.Content = content.String()
	project.Link = link
	project.URL = url

	return project, nil
}
