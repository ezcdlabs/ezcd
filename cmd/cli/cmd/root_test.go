package cmd_test

import (
	"bytes"
	"testing"

	"github.com/ezcdlabs/ezcd/cmd/cli/cmd"
	"github.com/stretchr/testify/assert"
)

func TestShouldPrintVersionNumber(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()
	// given
	c := cmd.NewRootCmd("1.0.0", mockServiceLoader)
	c.SetArgs([]string{"--version"})

	// Capture the output
	var stdout, stderr bytes.Buffer
	c.SetOut(&stdout)
	c.SetErr(&stderr)

	// when
	err := c.Execute()

	// then
	assert.NoError(t, err)
	assert.Empty(t, stderr.String())
	assert.Contains(t, stdout.String(), "1.0.0")
}

func TestShouldPrintHelp(t *testing.T) {
	mockServiceLoader := newMockServiceLoader()

	// given
	c := cmd.NewRootCmd("1.0.0", mockServiceLoader)
	c.SetArgs([]string{})

	// Capture the output
	var stdout, stderr bytes.Buffer
	c.SetOut(&stdout)
	c.SetErr(&stderr)

	// when
	err := c.Execute()

	// then
	assert.NoError(t, err)
	assert.Empty(t, stderr.String())
	assert.Contains(t, stdout.String(), "Usage")
}
