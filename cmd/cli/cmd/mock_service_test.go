package cmd_test

import (
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
	acceptanceStagePassedError  error
	acceptanceStageFailedError  error
	deployStartedError          error
	deployPassedError           error
	deployFailedError           error
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

func (m *mockEzcdService) AcceptanceStagePassed(projectId string, hash string) error {
	m.methodCalled = "AcceptanceStagePassed"
	if m.acceptanceStagePassedError != nil {
		return m.acceptanceStagePassedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

func (m *mockEzcdService) AcceptanceStageFailed(projectId string, hash string) error {
	m.methodCalled = "AcceptanceStageFailed"
	if m.acceptanceStageFailedError != nil {
		return m.acceptanceStageFailedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

// DeployStarted implements ezcd.Ezcd.
func (m *mockEzcdService) DeployStarted(projectId string, hash string) error {
	m.methodCalled = "DeployStarted"
	if m.deployStartedError != nil {
		return m.deployStartedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

// DeployPassed implements ezcd.Ezcd.
func (m *mockEzcdService) DeployPassed(projectId string, hash string) error {
	m.methodCalled = "DeployPassed"
	if m.deployPassedError != nil {
		return m.deployPassedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}

// DeployFailed implements ezcd.Ezcd.
func (m *mockEzcdService) DeployFailed(projectId string, hash string) error {
	m.methodCalled = "DeployFailed"
	if m.deployFailedError != nil {
		return m.deployFailedError
	}
	m.projectName = projectId
	m.commitHash = hash
	return nil
}
