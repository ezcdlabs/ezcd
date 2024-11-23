package cmd_test

import (
	"errors"
	"testing"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestShouldFailCreateProjectWhenNoDatabaseUrlSet(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.loadError = errors.New("failed to load service")

	command := cmd.NewCreateProjectCommand(mockServiceLoader)
	command.SetArgs([]string{"test"})
	err := command.Execute()

	assert.Error(t, err)
}

func TestShouldCallCreateProject(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockService := mockServiceLoader.mockEzcdService

	command := cmd.NewCreateProjectCommand(mockServiceLoader)
	command.SetArgs([]string{"test"})
	err := command.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "test", mockService.projectName)
	assert.Equal(t, "CreateProject", mockService.methodCalled)
}

func TestShouldFailCreateProjectWhenServiceReturnsErr(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockService := mockServiceLoader.mockEzcdService

	mockService.createProjectError = errors.New("test error")

	command := cmd.NewCreateProjectCommand(mockServiceLoader)
	command.SetArgs([]string{"test"})
	err := command.Execute()

	assert.Error(t, err)
}
