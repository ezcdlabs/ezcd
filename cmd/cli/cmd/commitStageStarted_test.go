package cmd_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldFailCommitStageStartedWhenNoDatabaseUrlSet(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.loadError = errors.New("failed to load service")

	command := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"ezcd", "commit-stage-started", "--project", "test", "--hash", "abc123", "--author-name", "John Doe", "--author-email", "john@example.com", "--message", "Initial commit", "--date", "2023-10-10T10:10:10Z"})
	err := command.Execute()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldCallCommitStageStarted(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	command := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test", "--hash", "abc123", "--author-name", "John Doe", "--author-email", "john@example.com", "--message", "Initial commit", "--date", "2023-10-10T10:10:10Z"})
	err := command.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assert.Equal(t, "test", mockServiceLoader.mockEzcdService.projectName)
	// make assertions for each field on commitData:
	assert.Equal(t, "abc123", mockServiceLoader.mockEzcdService.commitData.Hash)
	assert.Equal(t, "John Doe", mockServiceLoader.mockEzcdService.commitData.AuthorName)
	assert.Equal(t, "john@example.com", mockServiceLoader.mockEzcdService.commitData.AuthorEmail)
	assert.Equal(t, "Initial commit", mockServiceLoader.mockEzcdService.commitData.Message)
	assert.Equal(t, "2023-10-10T10:10:10Z", mockServiceLoader.mockEzcdService.commitData.Date.Format(time.RFC3339))
}

func TestShouldFailCommitStageStartedWhenServiceReturnsErr(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.mockEzcdService.commitStageStartedError = errors.New("test error")

	command := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test", "--hash", "abc123", "--author-name", "John Doe", "--author-email", "john@example.com", "--message", "Initial commit", "--date", "2023-10-10T10:10:10Z"})
	err := command.Execute()

	require.Error(t, err)
}

func TestCommitStageStartedMissingProject(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	x := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	x.SetArgs([]string{"--hash", "abc123", "--author-name", "John Doe", "--author-email", "john@example.com", "--message", "Initial commit", "--date", "2023-10-10T10:10:10Z"})
	err := x.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "project")
}

func TestCommitStageStartedMissingHash(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	x := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	x.SetArgs([]string{"--project", "test", "--author-name", "John Doe", "--author-email", "john@example.com", "--message", "Initial commit", "--date", "2023-10-10T10:10:10Z"})
	err := x.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "hash")
}

func TestCommitStageStartedMissingAuthorName(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	x := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	x.SetArgs([]string{"--project", "test", "--hash", "abc123", "--author-email", "john@example.com", "--message", "Initial commit", "--date", "2023-10-10T10:10:10Z"})
	err := x.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "author-name")
}

func TestCommitStageStartedMissingAuthorEmail(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	x := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	x.SetArgs([]string{"--project", "test", "--hash", "abc123", "--author-name", "John Doe", "--message", "Initial commit", "--date", "2023-10-10T10:10:10Z"})
	err := x.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "author-email")
}

func TestCommitStageStartedMissingMessage(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	x := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	x.SetArgs([]string{"--project", "test", "--hash", "abc123", "--author-name", "John Doe", "--author-email", "john@example.com", "--date", "2023-10-10T10:10:10Z"})
	err := x.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "message")
}

func TestCommitStageStartedMissingDate(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	x := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	x.SetArgs([]string{"--project", "test", "--hash", "abc123", "--author-name", "John Doe", "--author-email", "john@example.com", "--message", "Initial commit"})
	err := x.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "date")
}

func TestCommitStageStartedInvalidDate(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	x := cmd.NewCommitStageStartedCommand(mockServiceLoader)
	x.SetArgs([]string{"--project", "test", "--hash", "abc123", "--author-name", "John Doe", "--author-email", "john@example.com", "--message", "Initial commit", "--date", "invaliddate"})
	err := x.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "date")
}
