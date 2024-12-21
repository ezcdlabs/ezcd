package cmd_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldFailGetQueuedForAcceptanceWhenNoDatabaseUrlSet(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.loadError = errors.New("failed to load service")

	command := cmd.NewGetQueuedForAcceptanceCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test"})
	err := command.Execute()

	assert.Error(t, err)
}

func TestShouldPrintHashFromGetQueuedForAcceptance(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.mockEzcdService.queuedForAcceptanceHash = "abc123"

	// Create a buffer to capture the output
	var outputBuffer bytes.Buffer

	command := cmd.NewGetQueuedForAcceptanceCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test"})
	command.SetOut(&outputBuffer)
	err := command.Execute()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	actual := strings.TrimSpace(outputBuffer.String())

	assert.Equal(t, "abc123", actual)
}
func TestShouldShowErrorWhenGetQueuedForAcceptanceHasNoCommitInQueue(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	// Create a buffer to capture the output
	var outputBuffer bytes.Buffer

	command := cmd.NewGetQueuedForAcceptanceCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test"})
	command.SetOut(&outputBuffer)
	err := command.Execute()

	assert.Error(t, err, "expected an error, got none")
}

func TestShouldFailGetQueuedForAcceptanceWhenServiceReturnsErr(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	mockServiceLoader.mockEzcdService.getQueuedForAcceptanceError = errors.New("test error")

	command := cmd.NewGetQueuedForAcceptanceCommand(mockServiceLoader)
	command.SetArgs([]string{"--project", "test"})
	err := command.Execute()

	assert.Error(t, err)
}

func TestGetQueuedForAcceptanceMissingProject(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	command := cmd.NewGetQueuedForAcceptanceCommand(mockServiceLoader)
	command.SetArgs([]string{})
	err := command.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "project")
}
