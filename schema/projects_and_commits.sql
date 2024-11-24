CREATE TABLE IF NOT EXISTS projects (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    project_lock_id SERIAL NOT NULL
);

CREATE TABLE IF NOT EXISTS commits (
    project_id TEXT NOT NULL,
    commit_hash TEXT PRIMARY KEY NOT NULL,
    commit_author_name TEXT NOT NULL,
    commit_author_email TEXT NOT NULL,
    commit_message TEXT NOT NULL,
    commit_date TIMESTAMPTZ NOT NULL,

    commit_stage_started_at TIMESTAMPTZ,
    commit_stage_completed_at TIMESTAMPTZ,
    commit_stage_status TEXT CHECK(commit_stage_status IN ('started', 'passed', 'failed')) NOT NULL DEFAULT 'started',
    
    acceptance_stage_started_at TIMESTAMPTZ,
    acceptance_stage_completed_at TIMESTAMPTZ,
    acceptance_stage_status TEXT CHECK(acceptance_stage_status IN ('none', 'started', 'passed', 'failed')) NOT NULL DEFAULT 'none',
    
    deploy_started_at TIMESTAMPTZ,
    deploy_completed_at TIMESTAMPTZ,
    deploy_status TEXT CHECK(deploy_status IN ('none', 'started', 'passed', 'failed')) NOT NULL DEFAULT 'none',
    -- ltfc_completed TIMESTAMPTZ

    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);