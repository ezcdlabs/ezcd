package cmd_test

import (
	"errors"
	"testing"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldFailAcceptanceStageStartedWhenNoDatabaseUrlSet(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.loadError = errors.New("failed to load service")

	command := cmd.NewAcceptanceStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"ezcd", "commit-stage-passed", "--project", "test", "--hash", "abc123"})
	err := command.Execute()

	assert.Error(t, err)
}

func TestShouldCallAcceptanceStageStarted(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	command := cmd.NewAcceptanceStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test", "--hash", "abc123"})
	err := command.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, "test", mockServiceLoader.mockEzcdService.projectName)
	assert.Equal(t, "abc123", mockServiceLoader.mockEzcdService.commitHash)
}

func TestShouldFailAcceptanceStageStartedWhenServiceReturnsErr(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.mockEzcdService.acceptanceStageStartedError = errors.New("test error")

	command := cmd.NewAcceptanceStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test", "--hash", "abc123"})
	err := command.Execute()

	assert.Error(t, err)
}

func TestAcceptanceStageStartedMissingProject(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	command := cmd.NewAcceptanceStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"--hash", "abc123"})
	err := command.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "project")
}

func TestAcceptanceStageStartedMissingHash(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	command := cmd.NewAcceptanceStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test"})
	err := command.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "hash")
}
