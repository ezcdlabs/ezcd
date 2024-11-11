package ezcd

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetProjects() ([]Project, error) {
	connStr := os.Getenv("EZCD_DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	defer conn.Close()

	rows, err := conn.QueryContext(context.Background(), `SELECT id, name FROM projects`)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects from the database: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Name); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return projects, nil
}
