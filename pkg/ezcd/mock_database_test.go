package ezcd_test

import (
	"fmt"
	"sort"
	"time"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

type mockDatabase struct {
	Projects                map[string]ezcd.Project
	Commits                 map[string]ezcd.Commit
	SaveProjectError        error
	SaveCommitError         error
	BeginWorkError          error
	TransactionCommitError  error
	CheckConnectionError    error
	CheckProjectsTableError error
}

func newMockDatabase() *mockDatabase {
	return &mockDatabase{
		Projects: make(map[string]ezcd.Project),
		Commits:  make(map[string]ezcd.Commit),
	}
}

func (m *mockDatabase) GetInfo() string {
	return "mock database"
}

// GetCommits implements Database.
func (m *mockDatabase) GetCommits(id string) ([]ezcd.Commit, error) {
	commits := make([]ezcd.Commit, 0)
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
func (m *mockDatabase) CheckConnection() error {
	return m.CheckConnectionError
}

// CheckProjectsTable implements Database.
func (m *mockDatabase) CheckProjectsTable() error {
	return m.CheckProjectsTableError
}

// GetProject implements Database.
func (m *mockDatabase) GetProject(id string) (*ezcd.Project, error) {
	project, exists := m.Projects[id]
	if !exists {
		return nil, fmt.Errorf("project with id %s not found", id)
	}
	return &project, nil
}

// GetProjects implements Database.
func (m *mockDatabase) GetProjects() ([]ezcd.Project, error) {
	projects := make([]ezcd.Project, 0, len(m.Projects))
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

// BeginWork implements Database.
func (m *mockDatabase) BeginWork() (ezcd.UnitOfWork, error) {
	if m.BeginWorkError != nil {
		return nil, m.BeginWorkError
	}
	return &mockUnitOfWork{db: m}, nil
}

type mockUnitOfWork struct {
	db             *mockDatabase
	pendingActions []func()
}

// FindCommitForUpdate implements UnitOfWork.
func (m *mockUnitOfWork) FindCommitForUpdate(projectId string, hash string) (*ezcd.Commit, error) {
	commit, exists := m.db.Commits[hash]
	if !exists {
		return nil, fmt.Errorf("commit with id %s not found", hash)
	}
	if commit.Project != projectId {
		return nil, fmt.Errorf("commit with id %s does not belong to project %s", hash, projectId)
	}
	return &commit, nil
}

// FindUndeployedCommitsBeforeForUpdate implements UnitOfWork.
func (m *mockUnitOfWork) FindUndeployedCommitsBeforeForUpdate(projectId string, date time.Time) ([]ezcd.Commit, error) {
	commits := make([]ezcd.Commit, 0)
	for _, commit := range m.db.Commits {
		if commit.Project == projectId && commit.LeadTimeCompletedAt == nil && commit.Date.Before(date) {
			commits = append(commits, commit)
		}
	}
	sort.SliceStable(commits, func(i, j int) bool {
		return commits[i].Date.Before(commits[j].Date)
	})
	return commits, nil
}

// SaveCommit implements UnitOfWork.
func (m *mockUnitOfWork) SaveCommit(commit ezcd.Commit) error {
	if m.db.SaveCommitError != nil {
		return m.db.SaveCommitError
	}

	m.pendingActions = append(m.pendingActions, func() {
		m.db.Commits[commit.Hash] = commit
	})
	return nil
}

// WaitForProjectLock implements UnitOfWork.
func (m *mockUnitOfWork) WaitForProjectLock(id string) error {
	return nil
}

// Commit implements UnitOfWork.
func (m *mockUnitOfWork) Commit() error {
	if m.db.TransactionCommitError != nil {
		return m.db.TransactionCommitError
	}

	for _, commit := range m.pendingActions {
		commit()
	}
	return nil
}

// Rollback implements UnitOfWork.
func (m *mockUnitOfWork) Rollback() error {
	m.pendingActions = nil
	return nil
}

// FindProjectForUpdate implements UnitOfWork.
func (m *mockUnitOfWork) FindProjectForUpdate(id string) (*ezcd.Project, error) {
	return m.db.GetProject(id)
}

// SaveProject implements UnitOfWork.
func (m *mockUnitOfWork) SaveProject(project ezcd.Project) error {
	if m.db.SaveProjectError != nil {
		return m.db.SaveProjectError
	}

	m.pendingActions = append(m.pendingActions, func() {
		m.db.Projects[project.ID] = project
	})
	return nil
}
