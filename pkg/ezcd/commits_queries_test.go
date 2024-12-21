package ezcd_test

import (
	"testing"
	"time"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldHaveNilForQueuedForAcceptance(t *testing.T) {
	mockDB := newMockDatabase()
	service := ezcd.NewEzcdService(mockDB)

	service.CreateProject("project1")

	commit, err := service.GetQueuedForAcceptance("project1")

	require.NoError(t, err)
	require.Nil(t, commit)
}

func TestShouldShowPassedCommitAsQueuedForAcceptance(t *testing.T) {
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

	require.NoError(t, err)

	commit, err := service.GetQueuedForAcceptance("project1")

	require.NoError(t, err)
	require.NotNil(t, commit)
	assert.Equal(t, commitData.Hash, commit.Hash)
}

func TestShouldShowOnlyPassedCommitAsQueuedForAcceptance(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData1 := exampleCommitData("test commit 1", startTime)
	commitData2 := exampleCommitData("test commit 2", pointA)

	service.CommitStageStarted("project1", commitData1)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData2)

	mockClock.waitUntil(pointB)

	service.CommitStagePassed("project1", commitData1.Hash)

	commit, err := service.GetQueuedForAcceptance("project1")

	require.NoError(t, err)
	require.NotNil(t, commit)
	assert.Equal(t, commitData1.Hash, commit.Hash)
}

func TestShouldShowNewestPassedCommitAsQueuedForAcceptance(t *testing.T) {
	mockDB := newMockDatabase()
	mockClock := newMockClock()
	service := ezcd.NewEzcdService(mockDB)
	service.SetClock(mockClock)

	startTime := mockClock.CurrentTime
	pointA := startTime.Add(time.Second * 10)
	pointB := pointA.Add(time.Second * 10)
	pointC := pointB.Add(time.Second * 10)

	service.CreateProject("project1")

	commitData1 := exampleCommitData("test commit 1", startTime)
	commitData2 := exampleCommitData("test commit 2", pointA)

	service.CommitStageStarted("project1", commitData1)

	mockClock.waitUntil(pointA)

	service.CommitStageStarted("project1", commitData2)

	mockClock.waitUntil(pointB)

	service.CommitStagePassed("project1", commitData1.Hash)

	mockClock.waitUntil(pointC)

	err := service.CommitStagePassed("project1", commitData2.Hash)

	commit, err := service.GetQueuedForAcceptance("project1")

	require.NoError(t, err)
	require.NotNil(t, commit)
	assert.Equal(t, commitData2.Hash, commit.Hash)
}
