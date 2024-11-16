package ezcd_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

func TestShouldAddCommitToProject(t *testing.T) {
	mockDB := ezcd.NewMockDatabase()
	mockClock := ezcd.NewMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)

	projectID := "test-id"
	commitData := exampleCommitData("test commit", startTime)

	mockClock.WaitUntil(pointA)

	expectedCommit := ezcd.Commit{
		Project:              projectID,
		Hash:                 commitData.Hash,
		AuthorName:           commitData.AuthorName,
		AuthorEmail:          commitData.AuthorEmail,
		Message:              commitData.Message,
		Date:                 commitData.Date,
		CommitStageStatus:    ezcd.StatusStarted,
		CommitStageStartedAt: &pointA,
	}

	_, err := service.CreateProject(projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = service.CommitStageStarted(projectID, commitData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	commits, err := service.GetCommits(projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}

	if !reflect.DeepEqual(commits[0], expectedCommit) {
		t.Fatalf("expected commit from get response\n%+v,\ngot\n%+v", expectedCommit, commits[0])
	}
}

func TestShouldAddCommitAndSetItToPassed(t *testing.T) {
	mockDB := ezcd.NewMockDatabase()
	mockClock := ezcd.NewMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.WaitUntil(pointA)

	service.CommitStageStarted("project1", commitData)

	mockClock.WaitUntil(pointB)

	err := service.CommitStagePassed("project1", commitData.Hash)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	commits, err := service.GetCommits("project1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}

	if commits[0].CommitStageStatus != ezcd.StatusPassed {
		t.Fatalf("expected commit to be passed, got %v", commits[0].CommitStageStatus)
	}

	if *commits[0].CommitStageCompletedAt != pointB {
		t.Fatalf("expected commit completed at %v, got %v", pointB, commits[0].CommitStageCompletedAt)
	}
}

func exampleCommitData(message string, date time.Time) ezcd.CommitData {
	return ezcd.CommitData{
		Hash:        "abc123",
		AuthorName:  "test-author",
		AuthorEmail: "test-author-email",
		Message:     message,
		Date:        date,
	}
}
