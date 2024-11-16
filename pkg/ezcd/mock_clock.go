package ezcd

import (
	"time"
)

// MockClock is a mock implementation of the Clock interface
type MockClock struct {
	CurrentTime time.Time
}

// Now returns the current time set in the MockClock
func (m *MockClock) Now() *time.Time {
	now := m.CurrentTime
	return &now
}

// SetTime sets the current time in the MockClock
func (m *MockClock) WaitUntil(t time.Time) {
	m.CurrentTime = t
}

func NewMockClock() *MockClock {
	return &MockClock{
		CurrentTime: time.Now().UTC(), // TODO - should this be a fixed time?
	}
}
