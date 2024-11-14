package ezcd

import "fmt"

type Project struct {
	ID   string
	Name string
}

func (s *EzcdService) GetProject(id string) (*Project, error) {
	return s.db.GetProject(id)
}

func (s *EzcdService) GetProjects() ([]Project, error) {
	return s.db.GetProjects()
}

func (s *EzcdService) CreateProject(name string) (*Project, error) {
	uow, err := s.db.BeginWork()
	if err != nil {
		return nil, fmt.Errorf("failed to begin unit of work: %w", err)
	}

	defer uow.Rollback()

	project := Project{ID: name, Name: name}
	if err := uow.SaveProject(project); err != nil {
		return nil, fmt.Errorf("failed to save project: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &project, nil
}
