package cmd_test

import (
	"time"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

type mockServiceLoader struct {
	mockEzcdService *mockEzcdService
	loadError       error
}

func newMockServiceLoader() *mockServiceLoader {
	return &mockServiceLoader{
		mockEzcdService: &mockEzcdService{},
	}
}

func (m *mockServiceLoader) Load() (ezcd.Ezcd, error) {
	if m.loadError != nil {
		return nil, m.loadError
	}
	// Mock implementation
	return m.mockEzcdService, nil
}

// create a mock that implements ezcd.EzcdService
type mockEzcdService struct {
	projectName                 string
	commitHash                  string
	commitData                  ezcd.CommitData
	createProjectError          error
	commitStageStartedError     error
	commitStagePassedError      error
	acceptanceStageStartedError error
}

func (m *mockEzcdService) SetClock(clock ezcd.Clock) {
	// Mock implementation
}

func (m *mockEzcdService) CheckHealth() error {
	// Mock implementation
	return nil
}

func (m *mockEzcdService) GetDatabaseInfo() string {
	// Mock implementation
	return "mock database info"
}

func (m *mockEzcdService) GetProject(id string) (*ezcd.Project, error) {
	// Mock implementation
	return &ezcd.Project{}, nil
}

func (m *mockEzcdService) GetProjects() ([]ezcd.Project, error) {
	// Mock implementation
	return []ezcd.Project{}, nil
}

func (m *mockEzcdService) CreateProject(name string) error {
	if m.createProjectError != nil {
		return m.createProjectError
	}
	m.projectName = name
	// Mock implementation
	return nil
}

func (m *mockEzcdService) GetCommits(id string) ([]ezcd.Commit, error) {
	// Mock implementation
	return []ezcd.Commit{}, nil
}

func (m *mockEzcdService) CommitStageStarted(projectId string, commitData ezcd.CommitData) error {
	if m.commitStageStartedError != nil {
		return m.commitStageStartedError
	}
	m.projectName = projectId
	m.commitData = commitData

	// Mock implementation
	return nil
}

func (m *mockEzcdService) CommitStagePassed(projectId string, hash string) error {
	if m.commitStagePassedError != nil {
		return m.commitStagePassedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

func (m *mockEzcdService) AcceptanceStageStarted(projectId string, hash string) error {
	if m.acceptanceStageStartedError != nil {
		return m.acceptanceStageStartedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

// Mock Clock implementation
type mockClock struct{}

func (m *mockClock) Now() *time.Time {
	now := time.Now()
	return &now
}
