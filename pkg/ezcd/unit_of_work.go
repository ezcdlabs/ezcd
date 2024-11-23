package ezcd

import "fmt"

func (s *EzcdService) withUnitOfWork(fn func(uow UnitOfWork) error) error {
	uow, err := s.db.BeginWork()
	if err != nil {
		return fmt.Errorf("failed to begin unit of work: %w", err)
	}
	defer uow.Rollback()

	if err := fn(uow); err != nil {
		return err
	}

	if err := uow.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
