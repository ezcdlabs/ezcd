package ezcd_test

import (
	"testing"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

func TestShouldGetDatabaseInfo(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	info := service.GetDatabaseInfo()

	if info != mockDB.GetInfo() {
		t.Fatalf("expected %s, got %s", mockDB.GetInfo(), info)
	}
}
