package ezcd

import (
	"fmt"
	"sort"
	"time"
)

// Status represents the status of a commit stage.
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

// Commit represents a commit in the project.
type Commit struct {
	Project                    string
	Hash                       string
	AuthorName                 string
	AuthorEmail                string
	Message                    string
	Date                       time.Time
	CommitStageStartedAt       *time.Time
	CommitStageCompletedAt     *time.Time
	CommitStageStatus          Status
	AcceptanceStageStartedAt   *time.Time
	AcceptanceStageCompletedAt *time.Time
	AcceptanceStageStatus      Status
	DeployStartedAt            *time.Time
	DeployCompletedAt          *time.Time
	DeployStatus               Status
	LeadTimeCompletedAt        *time.Time
}

// CommitData represents the data of a commit.
type CommitData struct {
	Hash        string
	AuthorName  string
	AuthorEmail string
	Message     string
	Date        time.Time
}

// GetCommits retrieves the commits for a given project ID.
func (s *EzcdService) GetCommits(id string) ([]Commit, error) {
	return s.db.GetCommits(id)
}

func (s *EzcdService) GetQueuedForAcceptance(projectId string) (*Commit, error) {
	// dummy implementation where we just get the first commit from the commits
	commits, err := s.db.GetCommits(projectId)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits for project %v: %w", projectId, err)
	}

	// first sort the commits by commit date, newest first
	sort.SliceStable(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})

	// loop until we find a commit that has passed the commit stage
	for _, commit := range commits {
		if commit.CommitStageStatus == StatusPassed {
			// check that it hasn't also aleady started the acceptance stage
			if commit.AcceptanceStageStatus == StatusNone {
				return &commit, nil
			}

			// if it has started the acceptance stage, we can stop looking and return nil
			break
		}
	}

	return nil, nil
}

// CommitStageStarted marks the commit stage as started for a given project and commit data.
func (s *EzcdService) CommitStageStarted(projectId string, commitData CommitData) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		commit := &Commit{
			Project:               projectId,
			Hash:                  commitData.Hash,
			AuthorName:            commitData.AuthorName,
			AuthorEmail:           commitData.AuthorEmail,
			Message:               commitData.Message,
			Date:                  commitData.Date,
			CommitStageStatus:     StatusStarted,
			CommitStageStartedAt:  s.clock.Now(),
			AcceptanceStageStatus: StatusNone,
			DeployStatus:          StatusNone,
		}

		return s.saveCommit(uow, commit)
	})
}

// CommitStagePassed marks the commit stage as passed for a given project and commit hash.
func (s *EzcdService) CommitStagePassed(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)

		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		commit.CommitStageCompletedAt = s.clock.Now()
		commit.CommitStageStatus = StatusPassed

		return s.saveCommit(uow, commit)
	})
}

// CommitStagePassed marks the commit stage as passed for a given project and commit hash.
func (s *EzcdService) CommitStageFailed(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)

		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		commit.CommitStageCompletedAt = s.clock.Now()
		commit.CommitStageStatus = StatusFailed

		return s.saveCommit(uow, commit)
	})
}

// AcceptanceStageStarted marks the acceptance stage as started for a given project and commit hash.
func (s *EzcdService) AcceptanceStageStarted(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		// we need a project-level lock because the commit might not exist so there would be no commit row to lock
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)
		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		commit.AcceptanceStageStartedAt = s.clock.Now()
		commit.AcceptanceStageStatus = StatusStarted

		return s.saveCommit(uow, commit)
	})
}

// AcceptanceStagePassed marks the acceptance stage as passed for a given project and commit hash.
func (s *EzcdService) AcceptanceStagePassed(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		// we need a project-level lock because the commit might not exist so there would be no commit row to lock
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)
		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		commit.AcceptanceStageCompletedAt = s.clock.Now()
		commit.AcceptanceStageStatus = StatusPassed

		return s.saveCommit(uow, commit)
	})
}

// AcceptanceStageFailed marks the acceptance stage as failed for a given project and commit hash.
func (s *EzcdService) AcceptanceStageFailed(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		// we need a project-level lock because the commit might not exist so there would be no commit row to lock
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)
		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		commit.AcceptanceStageCompletedAt = s.clock.Now()
		commit.AcceptanceStageStatus = StatusFailed

		return s.saveCommit(uow, commit)
	})
}

// DeployStarted marks the deploy as started for a given project and commit hash.
func (s *EzcdService) DeployStarted(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		// we need a project-level lock because the commit might not exist so there would be no commit row to lock
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)
		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		commit.DeployStartedAt = s.clock.Now()
		commit.DeployStatus = StatusStarted

		return s.saveCommit(uow, commit)
	})
}

// DeployPassed marks the deploy as passed for a given project and commit hash.
func (s *EzcdService) DeployPassed(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		// we need a project-level lock because the commit might not exist so there would be no commit row to lock
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)
		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		now := s.clock.Now()
		commit.DeployCompletedAt = now
		commit.LeadTimeCompletedAt = now
		commit.DeployStatus = StatusPassed

		olderCommitsToStop, err := uow.FindUndeployedCommitsBeforeForUpdate(projectId, commit.Date)
		if err != nil {
			return fmt.Errorf("failed to query undeployed commits before commit with hash %v: %w", hash, err)
		}

		// for each commit in olderCommitsToStop:
		for _, olderCommit := range olderCommitsToStop {
			olderCommit.LeadTimeCompletedAt = now
			if err := s.saveCommit(uow, &olderCommit); err != nil {
				return fmt.Errorf("failed to save older commit with hash %v: %w", olderCommit.Hash, err)
			}
		}

		return s.saveCommit(uow, commit)
	})
}

// DeployFailed marks the deploy as Failed for a given project and commit hash.
func (s *EzcdService) DeployFailed(projectId string, hash string) error {
	return s.withUnitOfWork(func(uow UnitOfWork) error {
		// we need a project-level lock because the commit might not exist so there would be no commit row to lock
		uow.WaitForProjectLock(projectId)

		commit, err := uow.FindCommitForUpdate(projectId, hash)
		if err != nil {
			return fmt.Errorf("failed to find commit with hash %v: %w", hash, err)
		}

		commit.DeployCompletedAt = s.clock.Now()
		commit.DeployStatus = StatusFailed

		return s.saveCommit(uow, commit)
	})
}

func (s *EzcdService) saveCommit(uow UnitOfWork, commit *Commit) error {
	if err := uow.SaveCommit(*commit); err != nil {
		return fmt.Errorf("failed to save commit: %w", err)
	}
	return nil
}
