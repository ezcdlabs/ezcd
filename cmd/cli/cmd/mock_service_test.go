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
	methodCalled                string
	projectName                 string
	commitHash                  string
	commitData                  ezcd.CommitData
	createProjectError          error
	commitStageStartedError     error
	commitStagePassedError      error
	commitStageFailedError      error
	acceptanceStageStartedError error
}

func (m *mockEzcdService) SetClock(clock ezcd.Clock) {
}

func (m *mockEzcdService) CheckHealth() error {
	m.methodCalled = "CheckHealth"
	return nil
}

func (m *mockEzcdService) GetDatabaseInfo() string {
	m.methodCalled = "GetDatabaseInfo"
	return "mock database info"
}

func (m *mockEzcdService) GetProject(id string) (*ezcd.Project, error) {
	m.methodCalled = "GetProject"
	return &ezcd.Project{}, nil
}

func (m *mockEzcdService) GetProjects() ([]ezcd.Project, error) {
	m.methodCalled = "GetProjects"
	return []ezcd.Project{}, nil
}

func (m *mockEzcdService) CreateProject(name string) error {
	m.methodCalled = "CreateProject"
	if m.createProjectError != nil {
		return m.createProjectError
	}
	m.projectName = name
	// Mock implementation
	return nil
}

func (m *mockEzcdService) GetCommits(id string) ([]ezcd.Commit, error) {
	m.methodCalled = "GetCommits"
	return []ezcd.Commit{}, nil
}

func (m *mockEzcdService) CommitStageStarted(projectId string, commitData ezcd.CommitData) error {
	m.methodCalled = "CommitStageStarted"
	if m.commitStageStartedError != nil {
		return m.commitStageStartedError
	}
	m.projectName = projectId
	m.commitData = commitData
	return nil
}

func (m *mockEzcdService) CommitStageFailed(projectId string, hash string) error {
	m.methodCalled = "CommitStageFailed"
	if m.commitStageFailedError != nil {
		return m.commitStageFailedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

func (m *mockEzcdService) CommitStagePassed(projectId string, hash string) error {
	m.methodCalled = "CommitStagePassed"
	if m.commitStagePassedError != nil {
		return m.commitStagePassedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

func (m *mockEzcdService) AcceptanceStageStarted(projectId string, hash string) error {
	m.methodCalled = "AcceptanceStageStarted"
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
