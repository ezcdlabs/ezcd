package ezcd_test

import (
	"time"
)

// mockClock is a mock implementation of the Clock interface
type mockClock struct {
	CurrentTime time.Time
}

// Now returns the current time set in the MockClock
func (m *mockClock) Now() *time.Time {
	now := m.CurrentTime
	return &now
}

// SetTime sets the current time in the MockClock
func (m *mockClock) waitUntil(t time.Time) {
	m.CurrentTime = t
}

func newMockClock() *mockClock {
	return &mockClock{
		CurrentTime: time.Now().UTC(), // TODO - should this be a fixed time?
	}
}
