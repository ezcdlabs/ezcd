package ezcd_test

import (
	"errors"
	"testing"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

func TestShouldReturnHealthyWhenChecksPass(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	err := service.CheckHealth()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestShouldReturnUnhealthyWhenConnectionFails(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	mockDB.CheckConnectionError = errors.New("failed to connect")

	err := service.CheckHealth()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldReturnUnhealthyWhenProjectsTableIsMissing(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	mockDB.CheckProjectsTableError = errors.New("failed to find projects table")

	err := service.CheckHealth()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
