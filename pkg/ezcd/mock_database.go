package ezcd

import (
	"fmt"
	"sort"
)

type MockDatabase struct {
	Projects map[string]Project
	Commits  map[string]Commit
}

func NewMockDatabase() *MockDatabase {
	return &MockDatabase{
		Projects: make(map[string]Project),
		Commits:  make(map[string]Commit),
	}
}

// GetCommits implements Database.
func (m *MockDatabase) GetCommits(id string) ([]Commit, error) {
	commits := make([]Commit, 0)
	for _, commit := range m.Commits {
		if commit.Project == id {
			commits = append(commits, commit)
		}
	}
	sort.SliceStable(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})
	return commits, nil
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
	db             *MockDatabase
	pendingActions []func()
}

// FindCommitForUpdate implements UnitOfWork.
func (m *MockUnitOfWork) FindCommitForUpdate(projectId string, hash string) (*Commit, error) {
	commit, exists := m.db.Commits[hash]
	if !exists {
		return nil, fmt.Errorf("commit with id %s not found", hash)
	}
	if commit.Project != projectId {
		return nil, fmt.Errorf("commit with id %s does not belong to project %s", hash, projectId)
	}
	return &commit, nil
}

// SaveCommit implements UnitOfWork.
func (m *MockUnitOfWork) SaveCommit(commit Commit) error {
	m.pendingActions = append(m.pendingActions, func() {
		m.db.Commits[commit.Hash] = commit
	})
	return nil
}

// WaitForProjectLock implements UnitOfWork.
func (m *MockUnitOfWork) WaitForProjectLock(id string) error {
	return nil
}

func (m *MockUnitOfWork) Commit() error {
	for _, commit := range m.pendingActions {
		commit()
	}
	return nil
}

func (m *MockUnitOfWork) Rollback() error {
	m.pendingActions = nil
	return nil
}

func (m *MockUnitOfWork) FindProjectForUpdate(id string) (*Project, error) {
	return m.db.GetProject(id)
}

func (m *MockUnitOfWork) SaveProject(project Project) error {
	m.pendingActions = append(m.pendingActions, func() {
		m.db.Projects[project.ID] = project
	})
	return nil
}
