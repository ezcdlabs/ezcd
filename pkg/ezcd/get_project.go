package ezcd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetProject(id string) (*Project, error) {
	connStr := os.Getenv("EZCD_DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	defer conn.Close()

	row := conn.QueryRowContext(context.Background(), `SELECT id, name FROM projects WHERE id = $1`, id)

	var project Project
	if err := row.Scan(&project.ID, &project.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return &project, nil
}
