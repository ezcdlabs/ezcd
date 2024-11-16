package ezcd

import (
	"time"
)

func (s *EzcdService) SetClock(clock Clock) {
	s.clock = clock
}

// RealClock is a normal clock implementation
type RealClock struct{}

// Now returns the current local time
func (rc RealClock) Now() *time.Time {
	now := time.Now()
	return &now
}
