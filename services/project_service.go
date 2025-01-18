package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Project struct {
	Filename    string
	Name        string
	Description string
	Content     string
	URL         string
	Publish     bool

	// Status can be 'Done', 'Doing', 'Dump'
	Status string
}

type ProjectService struct {
	initialised bool

	Context       context.Context
	FileStorePath string
	Logger        *log.Logger
	Projects      []Project
}

func NewProjectService(
	ctx context.Context,
	path string,
	logger *log.Logger,
) *ProjectService {
	return &ProjectService{
		Context:       ctx,
		FileStorePath: path,
		Logger:        logger,
		Projects:      []Project{},
	}
}

func (ps *ProjectService) Init() error {
	if ps.initialised {
		return errors.New("ProjectService is already initialised")
	}

	data, err := os.ReadFile(ps.FileStorePath)
	if err != nil {
		ps.initialised = false
		return err
	}

	var projects []Project

	err = json.Unmarshal(data, &projects)
	if err != nil {
		return err
	}

	ps.initialised = true
	ps.Projects = projects
	return nil
}

func (ps *ProjectService) GetPublishedProjects() []Project {
	if !ps.initialised {
		return nil
	}

	var projects []Project
	for _, project := range ps.Projects {
		if project.Publish {
			projects = append(projects, project)
		}
	}

	return projects
}

func (ps *ProjectService) GetProjectByURL(url string) (p Project, ok bool) {
	if !ps.initialised {
		return
	}

	for _, project := range ps.Projects {
		if project.URL == url {
			return project, true
		}
	}

	return
}
