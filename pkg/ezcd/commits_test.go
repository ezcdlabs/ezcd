package ezcd_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/stretchr/testify/assert"
)

func TestShouldUseCorrectStatuses(t *testing.T) {
	if ezcd.StatusStarted.String() != "started" {
		t.Fatalf("expected started status to be 'started', got %s", ezcd.StatusStarted)
	}
	if ezcd.StatusNone.String() != "none" {
		t.Fatalf("expected started status to be 'none', got %s", ezcd.StatusNone)
	}
	if ezcd.StatusPassed.String() != "passed" {
		t.Fatalf("expected started status to be 'passed', got %s", ezcd.StatusPassed)
	}
	if ezcd.StatusFailed.String() != "failed" {
		t.Fatalf("expected started status to be 'failed', got %s", ezcd.StatusFailed)
	}
}

func TestShouldAddCommitToProject(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)

	projectID := "test-id"
	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	expectedCommit := ezcd.Commit{
		Project:               projectID,
		Hash:                  commitData.Hash,
		AuthorName:            commitData.AuthorName,
		AuthorEmail:           commitData.AuthorEmail,
		Message:               commitData.Message,
		Date:                  commitData.Date,
		CommitStageStatus:     ezcd.StatusStarted,
		CommitStageStartedAt:  &pointA,
		AcceptanceStageStatus: ezcd.StatusNone,
		DeployStatus:          ezcd.StatusNone,
	}

	err := service.CreateProject(projectID)
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

func TestShouldReturnErrWhenCommitCannotBeSaved(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)

	projectID := "test-id"
	commitData := exampleCommitData("test commit", mockClock.CurrentTime)

	err := service.CreateProject(projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mockDB.SaveCommitError = fmt.Errorf("failed to save commit")
	// define error here
	err = service.CommitStageStarted(projectID, commitData)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldReturnErrWhenBeginWorkFails(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)

	projectID := "test-id"
	commitData := exampleCommitData("test commit", mockClock.CurrentTime)

	err := service.CreateProject(projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mockDB.BeginWorkError = fmt.Errorf("failed to begin work")
	// define error here
	err = service.CommitStageStarted(projectID, commitData)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldReturnErrWhenTransactionFails(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)

	projectID := "test-id"
	commitData := exampleCommitData("test commit", mockClock.CurrentTime)

	err := service.CreateProject(projectID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	mockDB.TransactionCommitError = fmt.Errorf("failed to commit db transaction")
	// define error here
	err = service.CommitStageStarted(projectID, commitData)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldAddCommitAndSetItToPassed(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)

	mockClock.waitUntil(pointB)

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

func TestShouldAddCommitAndSetItToFailed(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)

	mockClock.waitUntil(pointB)

	err := service.CommitStageFailed("project1", commitData.Hash)

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

	if commits[0].CommitStageStatus != ezcd.StatusFailed {
		t.Fatalf("expected commit to be failed, got %v", commits[0].CommitStageStatus)
	}

	if *commits[0].CommitStageCompletedAt != pointB {
		t.Fatalf("expected commit completed at %v, got %v", pointB, commits[0].CommitStageCompletedAt)
	}
}

func TestShouldFailToPassCommitStageForCommitThatDoesNotExist(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	service.CreateProject("project1")

	err := service.CommitStagePassed("project1", "non-existent-hash")

	assert.Error(t, err)
}

func TestShouldFailToFailCommitStageForCommitThatDoesNotExist(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	service.CreateProject("project1")

	err := service.CommitStageFailed("project1", "non-existent-hash")

	assert.Error(t, err)
}

func TestShouldFailToStartAcceptanceStageForCommitThatDoesNotExist(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	service.CreateProject("project1")

	err := service.AcceptanceStageStarted("project1", "non-existent-hash")

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestShouldFailToPassAcceptanceStageForCommitThatDoesNotExist(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	service.CreateProject("project1")

	err := service.AcceptanceStagePassed("project1", "non-existent-hash")

	assert.Error(t, err)
}
func TestShouldFailToFailAcceptanceStageForCommitThatDoesNotExist(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	service.CreateProject("project1")

	err := service.AcceptanceStageFailed("project1", "non-existent-hash")

	assert.Error(t, err)
}

func TestShouldStartAcceptanceStage(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)
	service.CommitStagePassed("project1", commitData.Hash)

	mockClock.waitUntil(pointB)

	err := service.AcceptanceStageStarted("project1", commitData.Hash)

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

	if commits[0].AcceptanceStageStatus != ezcd.StatusStarted {
		t.Fatalf("expected acceptance stage to be started, got %v", commits[0].AcceptanceStageStatus)
	}

	if *commits[0].AcceptanceStageStartedAt != pointB {
		t.Fatalf("expected acceptance stage started at %v, got %v", pointB, commits[0].AcceptanceStageStartedAt)
	}
}

func TestShouldPassAcceptanceStage(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)
	pointC := pointB.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)
	service.CommitStagePassed("project1", commitData.Hash)

	mockClock.waitUntil(pointB)
	err := service.AcceptanceStageStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointC)
	err = service.AcceptanceStagePassed("project1", commitData.Hash)
	assert.NoError(t, err)

	commits, err := service.GetCommits("project1")
	assert.NoError(t, err)

	assert.Len(t, commits, 1)
	assert.Equal(t, ezcd.StatusPassed, commits[0].AcceptanceStageStatus)
	assert.Equal(t, pointC, *commits[0].AcceptanceStageCompletedAt)
}

func TestShouldFailAcceptanceStage(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)
	pointC := pointB.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)
	service.CommitStagePassed("project1", commitData.Hash)

	mockClock.waitUntil(pointB)
	err := service.AcceptanceStageStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointC)
	err = service.AcceptanceStageFailed("project1", commitData.Hash)
	assert.NoError(t, err)

	commits, err := service.GetCommits("project1")
	assert.NoError(t, err)

	assert.Len(t, commits, 1)
	assert.Equal(t, ezcd.StatusFailed, commits[0].AcceptanceStageStatus)
	assert.Equal(t, pointC, *commits[0].AcceptanceStageCompletedAt)
}

func TestShouldStartDeploy(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)
	pointC := pointB.Add(time.Second * 10)
	pointD := pointC.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)
	service.CommitStagePassed("project1", commitData.Hash)

	mockClock.waitUntil(pointB)
	err := service.AcceptanceStageStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointC)
	err = service.AcceptanceStagePassed("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointD)
	err = service.DeployStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	commits, err := service.GetCommits("project1")
	assert.NoError(t, err)

	assert.Len(t, commits, 1)
	assert.Equal(t, ezcd.StatusStarted, commits[0].DeployStatus)
	assert.Equal(t, pointD, *commits[0].DeployStartedAt)
}

func TestShouldPassDeploy(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)
	pointC := pointB.Add(time.Second * 10)
	pointD := pointC.Add(time.Second * 10)
	pointE := pointD.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)
	service.CommitStagePassed("project1", commitData.Hash)

	mockClock.waitUntil(pointB)
	err := service.AcceptanceStageStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointC)
	err = service.AcceptanceStagePassed("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointD)
	err = service.DeployStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointE)
	err = service.DeployPassed("project1", commitData.Hash)
	assert.NoError(t, err)

	commits, err := service.GetCommits("project1")
	assert.NoError(t, err)

	assert.Len(t, commits, 1)
	assert.Equal(t, ezcd.StatusPassed, commits[0].DeployStatus)
	assert.Equal(t, pointE, *commits[0].DeployCompletedAt)
}

func TestShouldSetLeadTimeCompletedAtWhenPassingDeploy(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)
	pointC := pointB.Add(time.Second * 10)
	pointD := pointC.Add(time.Second * 10)
	pointE := pointD.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)
	service.CommitStagePassed("project1", commitData.Hash)

	mockClock.waitUntil(pointB)
	err := service.AcceptanceStageStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointC)
	err = service.AcceptanceStagePassed("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointD)
	err = service.DeployStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointE)
	err = service.DeployPassed("project1", commitData.Hash)
	assert.NoError(t, err)

	commits, err := service.GetCommits("project1")
	assert.NoError(t, err)

	assert.Len(t, commits, 1)
	assert.Equal(t, pointE, *commits[0].LeadTimeCompletedAt)
}

func TestShouldFailDeploy(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)
	pointC := pointB.Add(time.Second * 10)
	pointD := pointC.Add(time.Second * 10)
	pointE := pointD.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData := exampleCommitData("test commit", startTime)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData)
	service.CommitStagePassed("project1", commitData.Hash)

	mockClock.waitUntil(pointB)
	err := service.AcceptanceStageStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointC)
	err = service.AcceptanceStagePassed("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointD)
	err = service.DeployStarted("project1", commitData.Hash)
	assert.NoError(t, err)

	mockClock.waitUntil(pointE)
	err = service.DeployFailed("project1", commitData.Hash)
	assert.NoError(t, err)

	commits, err := service.GetCommits("project1")
	assert.NoError(t, err)

	assert.Len(t, commits, 1)
	assert.Equal(t, ezcd.StatusFailed, commits[0].DeployStatus)
	assert.Equal(t, pointE, *commits[0].DeployCompletedAt)
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
