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
    created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL

    -- commit_phase_started TIMESTAMPTZ,
    -- commit_phase_completed TIMESTAMPTZ,
    -- commit_phase_status TEXT CHECK(commit_phase_status IN ('pass', 'fail')),
    
    -- acceptance_phase_started TIMESTAMPTZ,
    -- acceptance_phase_completed TIMESTAMPTZ,
    -- acceptance_phase_status TEXT CHECK(acceptance_phase_status IN ('pass', 'fail')),
    
    -- released_to_production TIMESTAMPTZ,
    -- ltfc_completed TIMESTAMPTZ
);