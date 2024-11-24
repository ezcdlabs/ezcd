package ezcd_postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ezcdlabs/ezcd/pkg/ezcd"
	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	// database url
	databaseUrl string
}

// GetInfo implements ezcd.Database.
func (p *PostgresDatabase) GetInfo() string {
	return p.databaseUrl
}

// GetCommits implements ezcd.Database.
func (p *PostgresDatabase) GetCommits(id string) ([]ezcd.Commit, error) {
	connStr := p.databaseUrl
	if connStr == "" {
		return nil, fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(context.Background(), `
	SELECT 
		commit_hash,
		commit_author_name,
		commit_author_email,
		commit_message,
		commit_date,
		commit_stage_started_at,
		commit_stage_completed_at,
		commit_stage_status,
		acceptance_stage_started_at,
		acceptance_stage_completed_at,
		acceptance_stage_status,
		deploy_started_at,
		deploy_completed_at,
		deploy_status
	FROM commits WHERE project_id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query commits from the database: %w", err)
	}
	defer rows.Close()

	var commits []ezcd.Commit
	for rows.Next() {
		var commit ezcd.Commit

		if err := rows.Scan(
			&commit.Hash,
			&commit.AuthorName,
			&commit.AuthorEmail,
			&commit.Message,
			&commit.Date,
			&commit.CommitStageStartedAt,
			&commit.CommitStageCompletedAt,
			&commit.CommitStageStatus,
			&commit.AcceptanceStageStartedAt,
			&commit.AcceptanceStageCompletedAt,
			&commit.AcceptanceStageStatus,
			&commit.DeployStartedAt,
			&commit.DeployCompletedAt,
			&commit.DeployStatus,
		); err != nil {
			return nil, fmt.Errorf("failed to scan commit: %w", err)
		}
		commits = append(commits, commit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return commits, nil
}

type PostgresUnitOfWork struct {
	tx *sql.Tx
}

// FindCommitForUpdate implements ezcd.UnitOfWork.
func (u *PostgresUnitOfWork) FindCommitForUpdate(projectId string, hash string) (*ezcd.Commit, error) {
	row := u.tx.QueryRowContext(context.Background(), `
	SELECT 
		commit_hash,
		commit_author_name,
		commit_author_email,
		commit_message,
		commit_date,
		commit_stage_started_at,
		commit_stage_completed_at,
		commit_stage_status,
		acceptance_stage_started_at,
		acceptance_stage_completed_at,
		acceptance_stage_status,
		deploy_started_at,
		deploy_completed_at,
		deploy_status
	FROM commits 
	WHERE 
		project_id = $1 AND
		commit_hash = $2
	FOR UPDATE`,
		projectId, hash)

	var commit ezcd.Commit
	if err := row.Scan(
		&commit.Hash,
		&commit.AuthorName,
		&commit.AuthorEmail,
		&commit.Message,
		&commit.Date,
		&commit.CommitStageStartedAt,
		&commit.CommitStageCompletedAt,
		&commit.CommitStageStatus,
		&commit.AcceptanceStageStartedAt,
		&commit.AcceptanceStageCompletedAt,
		&commit.AcceptanceStageStatus,
		&commit.DeployStartedAt,
		&commit.DeployCompletedAt,
		&commit.DeployStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("commit not found: %w", err)
		}
		return nil, fmt.Errorf("failed to scan commit: %w", err)
	}

	return &commit, nil
}

// SaveCommit implements ezcd.UnitOfWork.
func (u *PostgresUnitOfWork) SaveCommit(commit ezcd.Commit) error {
	_, err := u.tx.ExecContext(context.Background(), `
	INSERT INTO commits (
		project_id,
		commit_hash,
		commit_author_name,
		commit_author_email,
		commit_message,
		commit_date,
		commit_stage_started_at,
		commit_stage_completed_at,
		commit_stage_status,
		acceptance_stage_started_at,
		acceptance_stage_completed_at,
		acceptance_stage_status,
		deploy_started_at,
		deploy_completed_at,
		deploy_status
	) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	ON CONFLICT (commit_hash) DO UPDATE 
	SET commit_author_name = EXCLUDED.commit_author_name,
		commit_author_email = EXCLUDED.commit_author_email,
		commit_message = EXCLUDED.commit_message,
		commit_date = EXCLUDED.commit_date,
		commit_stage_started_at = EXCLUDED.commit_stage_started_at,
		commit_stage_completed_at = EXCLUDED.commit_stage_completed_at,
		commit_stage_status = EXCLUDED.commit_stage_status,
		acceptance_stage_started_at = EXCLUDED.acceptance_stage_started_at,
		acceptance_stage_completed_at = EXCLUDED.acceptance_stage_completed_at,
		acceptance_stage_status = EXCLUDED.acceptance_stage_status,
		deploy_started_at = EXCLUDED.deploy_started_at,
		deploy_completed_at = EXCLUDED.deploy_completed_at,
		deploy_status = EXCLUDED.deploy_status`,

		commit.Project,
		commit.Hash,
		commit.AuthorName,
		commit.AuthorEmail,
		commit.Message,
		commit.Date,
		commit.CommitStageStartedAt,
		commit.CommitStageCompletedAt,
		commit.CommitStageStatus,
		commit.AcceptanceStageStartedAt,
		commit.AcceptanceStageCompletedAt,
		commit.AcceptanceStageStatus,
		commit.DeployStartedAt,
		commit.DeployCompletedAt,
		commit.DeployStatus)
	if err != nil {
		return fmt.Errorf("failed to save commit: %w", err)
	}
	return nil
}

// WaitForProjectLock implements ezcd.UnitOfWork.
func (u *PostgresUnitOfWork) WaitForProjectLock(id string) error {
	var projectID string
	var projectLockId int
	row := u.tx.QueryRowContext(context.Background(), `SELECT id, project_lock_id FROM projects WHERE id = $1`, projectID)
	err := row.Scan(&projectID, &projectLockId)

	if err != nil {
		return fmt.Errorf("failed to fetch project: %w", err)
	}

	// Acquire advisory lock
	_, err = u.tx.ExecContext(context.Background(), `SELECT pg_advisory_lock($1)`, projectLockId)
	if err != nil {
		return fmt.Errorf("failed to lock project: %w", err)
	}

	return nil
}

// FindProjectForUpdate implements ezcd.UnitOfWork.
func (u *PostgresUnitOfWork) FindProjectForUpdate(id string) (*ezcd.Project, error) {
	row := u.tx.QueryRowContext(context.Background(), `SELECT id, name FROM projects WHERE id = $1 FOR UPDATE`, id)

	var project ezcd.Project
	if err := row.Scan(&project.ID, &project.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %w", err)
		}
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return &project, nil
}

// SaveProject implements ezcd.UnitOfWork.
func (u *PostgresUnitOfWork) SaveProject(project ezcd.Project) error {
	_, err := u.tx.ExecContext(context.Background(), `
		INSERT INTO projects (id, name) 
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE 
		SET name = EXCLUDED.name`,
		project.ID, project.Name)
	if err != nil {
		return fmt.Errorf("failed to save project: %w", err)
	}
	return nil
}

// Commit implements ezcd.UnitOfWork.
func (u *PostgresUnitOfWork) Commit() error {
	return u.tx.Commit()
}

// Rollback implements ezcd.UnitOfWork.
func (u *PostgresUnitOfWork) Rollback() error {
	return u.tx.Rollback()
}

// BeginWork implements ezcd.Database.
func (p *PostgresDatabase) BeginWork() (ezcd.UnitOfWork, error) {
	connStr := p.databaseUrl
	if connStr == "" {
		return nil, fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	tx, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to begin transaction: %v", err)
	}

	return &PostgresUnitOfWork{tx: tx}, nil
}

// CheckConnection implements ezcd.Database.
func (p *PostgresDatabase) CheckConnection() error {
	connStr := p.databaseUrl
	if connStr == "" {
		return fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	row := conn.QueryRowContext(context.Background(), `SELECT 1`)
	var result int
	if err := row.Scan(&result); err != nil {
		return fmt.Errorf("failed to execute test query: %w", err)
	}

	return nil
}

// CheckProjectsTable implements ezcd.Database.
func (p *PostgresDatabase) CheckProjectsTable() error {

	connStr := p.databaseUrl
	if connStr == "" {
		return fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	row := conn.QueryRowContext(context.Background(), `
	SELECT EXISTS (
		SELECT FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_name = 'projects'
	)`)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return fmt.Errorf("failed to execute check query: %w", err)
	}

	if !exists {
		return fmt.Errorf("projects table does not exist")
	}

	return nil
}

// GetProject implements ezcd.Database.
func (p *PostgresDatabase) GetProject(id string) (*ezcd.Project, error) {
	connStr := p.databaseUrl
	if connStr == "" {
		return nil, fmt.Errorf("database connection string is required, please set EZCD_DATABASE_URL")
	}

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	defer conn.Close()

	row := conn.QueryRowContext(context.Background(), `SELECT id, name FROM projects WHERE id = $1`, id)

	var project ezcd.Project
	if err := row.Scan(&project.ID, &project.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %w", err)
		}
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	return &project, nil
}

// GetProjects implements ezcd.Database.
func (p *PostgresDatabase) GetProjects() ([]ezcd.Project, error) {
	connStr := p.databaseUrl
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

	var projects []ezcd.Project
	for rows.Next() {
		var project ezcd.Project
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

func NewPostgresDatabase(databaseUrl string) ezcd.Database {
	return &PostgresDatabase{
		databaseUrl: databaseUrl,
	}
}
