package ezcd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// create a new project
func CreateProject(name string) (*string, error) {

	// do the dumbest possible implementation for now
	connStr := os.Getenv("EZCD_DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	defer conn.Close()

	// generate a uuid for the project id:
	// TODO instead of using uuid, we convert the name to a slug
	id := uuid.New().String()

	_, err = conn.ExecContext(context.Background(),
		`INSERT INTO 
			projects 
			(
				id,
				name
			)
			VALUES ($1, $2)`,
		id, name,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add commit to the database with : %w", err)
	}

	return &id, nil
}
