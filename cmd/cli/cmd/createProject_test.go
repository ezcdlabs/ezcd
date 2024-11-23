package cmd_test

import (
	"errors"
	"testing"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
)

func TestShouldFailCreateProjectWhenNoDatabaseUrlSet(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.loadError = errors.New("failed to load service")

	command := cmd.NewCreateProjectCommand(mockServiceLoader)
	command.SetArgs([]string{"test"})
	err := command.Execute()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldCallCreateProject(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockService := mockServiceLoader.mockEzcdService

	command := cmd.NewCreateProjectCommand(mockServiceLoader)
	command.SetArgs([]string{"test"})
	err := command.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mockService.projectName != "test" {
		t.Fatalf("expected project name test, got %v", mockService.projectName)
	}
}

func TestShouldFailCreateProjectWhenServiceReturnsErr(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockService := mockServiceLoader.mockEzcdService

	mockService.createProjectError = errors.New("test error")

	command := cmd.NewCreateProjectCommand(mockServiceLoader)
	command.SetArgs([]string{"test"})
	err := command.Execute()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
