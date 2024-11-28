CREATE TABLE IF NOT EXISTS projects(
    id text PRIMARY KEY NOT NULL,
    name text NOT NULL,
    project_lock_id serial NOT NULL
);

CREATE TABLE IF NOT EXISTS commits(
    project_id text NOT NULL,
    commit_hash text PRIMARY KEY NOT NULL,
    commit_author_name text NOT NULL,
    commit_author_email text NOT NULL,
    commit_message text NOT NULL,
    commit_date timestamptz NOT NULL,
    commit_stage_started_at timestamptz,
    commit_stage_completed_at timestamptz,
    commit_stage_status text CHECK (commit_stage_status IN ('started', 'passed', 'failed')) NOT NULL DEFAULT 'started',
    acceptance_stage_started_at timestamptz,
    acceptance_stage_completed_at timestamptz,
    acceptance_stage_status text CHECK (acceptance_stage_status IN ('none', 'started', 'passed', 'failed')) NOT NULL DEFAULT 'none',
    deploy_started_at timestamptz,
    deploy_completed_at timestamptz,
    deploy_status text CHECK (deploy_status IN ('none', 'started', 'passed', 'failed')) NOT NULL DEFAULT 'none',
    lead_time_completed_at timestamptz,
    created_at timestamptz DEFAULT NOW() NOT NULL
);

