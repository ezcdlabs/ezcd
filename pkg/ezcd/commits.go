package ezcd

import (
	"fmt"
	"log"
	"time"
)

type Commit struct {
	Project     string
	Hash        string
	AuthorName  string
	AuthorEmail string
	Message     string
	Date        time.Time
}

type CommitData struct {
	Project     string
	Hash        string
	AuthorName  string
	AuthorEmail string
	Message     string
	Date        time.Time
}

func (s *EzcdService) GetCommits(id string) ([]Commit, error) {
	return s.db.GetCommits(id)
}

func (s *EzcdService) CommitPhaseStarted(commitData CommitData) (*Commit, error) {
	uow, err := s.db.BeginWork()
	if err != nil {
		return nil, fmt.Errorf("failed to begin unit of work: %w", err)
	}

	defer uow.Rollback()

	// we need a project-level lock because the commit might not exist so there would be no commit row to lock
	uow.WaitForProjectLock(commitData.Project)

	commit, err := uow.FindCommitForUpdate(commitData.Hash)

	log.Printf("existing commit: %v\n", commit)

	if err != nil {
		commit = &Commit{
			Project:     commitData.Project,
			Hash:        commitData.Hash,
			AuthorName:  commitData.AuthorName,
			AuthorEmail: commitData.AuthorEmail,
			Message:     commitData.Message,
			Date:        commitData.Date,
			// TODO add the commit phase started info?
		} // TODO, we need the commit to have the data from the commitData
		log.Printf("new commit: %v\n", commit)
	}

	if err := uow.SaveCommit(*commit); err != nil {
		return nil, fmt.Errorf("failed to save commit: %w", err)
	}

	if err := uow.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return commit, nil
}
