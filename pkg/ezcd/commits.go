package ezcd

import (
	"fmt"
	"time"
)

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	StatusNone    Status = "none"
	StatusStarted Status = "started"
	StatusPassed  Status = "passed"
	StatusFailed  Status = "failed"
)

type Commit struct {
	Project                string
	Hash                   string
	AuthorName             string
	AuthorEmail            string
	Message                string
	Date                   time.Time
	CommitStageStartedAt   *time.Time
	CommitStageCompletedAt *time.Time
	CommitStageStatus      Status
}

type CommitData struct {
	Hash        string
	AuthorName  string
	AuthorEmail string
	Message     string
	Date        time.Time
}

func (s *EzcdService) GetCommits(id string) ([]Commit, error) {
	return s.db.GetCommits(id)
}

func (s *EzcdService) CommitStageStarted(projectId string, commitData CommitData) error {
	uow, err := s.db.BeginWork()
	if err != nil {
		return fmt.Errorf("failed to begin unit of work: %w", err)
	}

	defer uow.Rollback()

	// we need a project-level lock because the commit might not exist so there would be no commit row to lock
	uow.WaitForProjectLock(projectId)

	commit, err := uow.FindCommitForUpdate(projectId, commitData.Hash)

	if err != nil {
		commit = &Commit{
			Project:     projectId,
			Hash:        commitData.Hash,
			AuthorName:  commitData.AuthorName,
			AuthorEmail: commitData.AuthorEmail,
			Message:     commitData.Message,
			Date:        commitData.Date,
		}
	}

	commit.CommitStageStartedAt = s.clock.Now()
	commit.CommitStageStatus = StatusStarted
	commit.CommitStageCompletedAt = nil

	if err := uow.SaveCommit(*commit); err != nil {
		return fmt.Errorf("failed to save commit: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *EzcdService) CommitStagePassed(projectId string, hash string) error {
	uow, err := s.db.BeginWork()
	if err != nil {
		return fmt.Errorf("failed to begin unit of work: %w", err)
	}

	defer uow.Rollback()

	// we need a project-level lock because the commit might not exist so there would be no commit row to lock
	uow.WaitForProjectLock(projectId)

	commit, err := uow.FindCommitForUpdate(projectId, hash)

	if err != nil {
		return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
	}

	commit.CommitStageCompletedAt = s.clock.Now()
	commit.CommitStageStatus = StatusPassed

	if err := uow.SaveCommit(*commit); err != nil {
		return fmt.Errorf("failed to save commit: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
