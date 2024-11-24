package ezcd

import (
	"fmt"
)

func (s *EzcdService) CheckHealth() error {
	// check if the database is able to be reached:
	if err := s.db.CheckConnection(); err != nil {
		return fmt.Errorf("unable to reach database, check connection: %w", err)
	}

	// check if the database has the projects table:
	if err := s.db.CheckProjectsTable(); err != nil {
		return fmt.Errorf("database connection is unable to find projects table: %w", err)
	}

	return nil
}
