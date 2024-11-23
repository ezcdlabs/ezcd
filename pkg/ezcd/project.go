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

func (s *EzcdService) CreateProject(name string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		project := Project{ID: name, Name: name}
		if err := uow.SaveProject(project); err != nil {
			return fmt.Errorf("failed to save project: %w", err)
		}
		return nil
	})
}
