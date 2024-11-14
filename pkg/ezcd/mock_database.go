package ezcd

import (
	"fmt"
	"sort"
)

type MockDatabase struct {
	Projects map[string]Project
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		Projects: make(map[string]Project),
	}
}

// CheckConnection implements Database.
func (m *MockDatabase) CheckConnection() error {
	panic("unimplemented")
}

// CheckProjectsTable implements Database.
func (m *MockDatabase) CheckProjectsTable() error {
	panic("unimplemented")
}

func (m *MockDatabase) GetProject(id string) (*Project, error) {
	project, exists := m.Projects[id]
	if !exists {
		return nil, fmt.Errorf("project with id %s not found", id)
	}
	return &project, nil
}

func (m *MockDatabase) GetProjects() ([]Project, error) {
	projects := make([]Project, 0, len(m.Projects))
	keys := make([]string, 0, len(m.Projects))
	for key := range m.Projects {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		projects = append(projects, m.Projects[key])
	}
	return projects, nil
}

func (m *MockDatabase) BeginWork() (UnitOfWork, error) {
	// Mock implementation
	return &MockUnitOfWork{db: m}, nil
}

type MockUnitOfWork struct {
	db      *MockDatabase
	commits []func()
}

func (m *MockUnitOfWork) Commit() error {
	for _, commit := range m.commits {
		commit()
	}
	return nil
}

func (m *MockUnitOfWork) Rollback() error {
	m.commits = nil
	return nil
}

func (m *MockUnitOfWork) FindProjectForUpdate(id string) (*Project, error) {
	return m.db.GetProject(id)
}

func (m *MockUnitOfWork) SaveProject(project Project) error {
	m.commits = append(m.commits, func() {
		m.db.Projects[project.ID] = project
	})
	return nil
}
