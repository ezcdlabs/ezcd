package ezcd_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
)

func TestShouldAddCommitToProject(t *testing.T) {
	mockDB := ezcd.NewMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	projectID := "test-id"
	commitData := ezcd.CommitData{
		Project:     projectID,
		Hash:        "abc123",
		AuthorName:  "test-author",
		AuthorEmail: "test-author-email",
		Message:     "test commit",
		Date:        time.Now(),
	}

	expectedCommit := ezcd.Commit{
		Project:     projectID,
		Hash:        commitData.Hash,
		AuthorName:  commitData.AuthorName,
		AuthorEmail: commitData.AuthorEmail,
		Message:     commitData.Message,
		Date:        commitData.Date,
	}

	_, err := service.CreateProject(projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	commit, err := service.CommitPhaseStarted(commitData)
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

	if !reflect.DeepEqual(commit, &expectedCommit) {
		t.Fatalf("expected commit from command response:\n%+v\ngot:\n%+v", expectedCommit, commit)
	}

	if !reflect.DeepEqual(commits[0], expectedCommit) {
		t.Fatalf("expected commit from get response\n%+v,\ngot\n%+v", expectedCommit, commit)
	}
}
