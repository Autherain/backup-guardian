CREATE TABLE sync_runs (
    id TEXT PRIMARY KEY,
    job_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    started_at DATETIME,
    finished_at DATETIME,
    error_message TEXT,
    files_transferred INTEGER DEFAULT 0,
    bytes_transferred INTEGER DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
