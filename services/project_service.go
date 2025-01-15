package services

type Project struct{}

type ProjectService struct{}

func NewProjectService() *ProjectService { return &ProjectService{} }

func (ps ProjectService) GetProjects() []Project {
	return []Project{}
}
