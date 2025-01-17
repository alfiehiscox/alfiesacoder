package services

import (
	"context"
	"errors"
	"log"
	"os"
	"path"

	"github.com/yuin/goldmark"
)

type ContentType interface {
	GetURL() string
}

type ExtractionFunction[T ContentType] func(
	data []byte,
	md goldmark.Markdown,
	filename string,
) (T, error)

type ContentService[T ContentType] struct {
	initialised bool

	Context       context.Context
	Logger        *log.Logger
	ContextPath   string
	Content       map[string]T
	BaseURL       string
	Markdown      goldmark.Markdown
	ExtractFromMD ExtractionFunction[T]
}

func NewContentService[T ContentType](
	ctx context.Context,
	path string,
	logger *log.Logger,
	md goldmark.Markdown,
	extract_from_md ExtractionFunction[T],
) *ContentService[T] {
	return &ContentService[T]{
		Logger:        logger,
		Context:       ctx,
		ContextPath:   path,
		Content:       make(map[string]T),
		Markdown:      md,
		ExtractFromMD: extract_from_md,
	}
}

func (cs *ContentService[T]) Init() error {
	if cs.initialised {
		return errors.New("ContentService is already initialised")
	}

	contentEntries, err := os.ReadDir(cs.ContextPath)
	if err != nil {
		cs.initialised = false
		return err
	}

	for _, entry := range contentEntries {

		data, err := os.ReadFile(path.Join(cs.ContextPath, entry.Name()))
		if err != nil {
			cs.initialised = false
			return err
		}

		if cs.ExtractFromMD == nil {
			return errors.New("No extraction function provided")
		}

		content, err := cs.ExtractFromMD(data, cs.Markdown, entry.Name())
		if err != nil {
			return err
		}

		cs.Content[content.GetURL()] = content
	}

	cs.initialised = true
	return nil
}

func (cs *ContentService[T]) GetContentByURL(url string) (c T, ok bool) {
	if !cs.initialised {
		return
	}

	c, ok = cs.Content[url]
	return
}

func (cs *ContentService[T]) GetContent() (contents []T) {
	if !cs.initialised {
		return
	}

	for _, content := range cs.Content {
		contents = append(contents, content)
	}

	return contents
}
