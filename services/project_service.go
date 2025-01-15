package services

import "github.com/yuin/goldmark"

type Project struct{}

type ProjectService struct{}

func NewProjectService(markdown goldmark.Markdown) *ProjectService {
	return &ProjectService{}
}

func (ps ProjectService) GetProjects() []Project {
	return []Project{}
}
