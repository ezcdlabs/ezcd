package ezcd

import "time"

// Database interface defines the methods for interacting with the database.
type Database interface {
	CheckConnection() error
	CheckProjectsTable() error

	GetProject(id string) (*Project, error)
	GetProjects() ([]Project, error)
	GetCommits(id string) ([]Commit, error)

	BeginWork() (UnitOfWork, error)
}

type Clock interface {
	Now() *time.Time
}

// UnitOfWork interface defines the methods for managing a unit of work (transaction).
type UnitOfWork interface {
	Commit() error
	Rollback() error

	FindProjectForUpdate(id string) (*Project, error)
	SaveProject(project Project) error

	WaitForProjectLock(id string) error

	FindCommitForUpdate(projectId string, hash string) (*Commit, error)
	SaveCommit(commit Commit) error
}

// Ezcd interface defines the service methods for the Ezcd application.
type Ezcd interface {

	// for testing, we can override the clock that is used to get the current time
	SetClock(clock Clock)

	// health.go
	CheckHealth() error

	// project.go
	GetProject(id string) (*Project, error)
	GetProjects() ([]Project, error)
	CreateProject(name string) (*Project, error)

	// commits.go
	GetCommits(id string) ([]Commit, error)
	CommitStageStarted(projectId string, commitData CommitData) error
	CommitStagePassed(projectId string, hash string) error
}

type EzcdService struct {
	db    Database
	clock Clock
}

// NewEzcdService initializes a new EzcdService with a database dependency.
func NewEzcdService(db Database) Ezcd {
	return &EzcdService{db: db, clock: RealClock{}}
}
