package ezcd

// Database interface defines the methods for interacting with the database.
type Database interface {
	CheckConnection() error
	CheckProjectsTable() error

	GetProject(id string) (*Project, error)
	GetProjects() ([]Project, error)

	BeginWork() (UnitOfWork, error)
}

// UnitOfWork interface defines the methods for managing a unit of work (transaction).
type UnitOfWork interface {
	Commit() error
	Rollback() error

	FindProjectForUpdate(id string) (*Project, error)
	SaveProject(project Project) error
}

// Ezcd interface defines the service methods for the Ezcd application.
type Ezcd interface {
	// health.go
	CheckHealth() error

	// project.go
	GetProject(id string) (*Project, error)
	GetProjects() ([]Project, error)
	CreateProject(name string) (*Project, error)
}

type EzcdService struct {
	db Database
}

// NewEzcdService initializes a new EzcdService with a database dependency.
func NewEzcdService(db Database) Ezcd {
	return &EzcdService{db: db}
}
