-- name: CreateSyncRun :one
INSERT INTO sync_runs (id, job_name, status, started_at)
VALUES (?, ?, 'running', ?)
RETURNING *;

-- name: UpdateSyncRun :exec
UPDATE sync_runs
SET status = ?,
    finished_at = ?,
    error_message = ?,
    files_transferred = ?,
    bytes_transferred = ?
WHERE id = ?;

-- name: GetSyncRun :one
SELECT * FROM sync_runs
WHERE id = ?;

-- name: ListSyncRuns :many
SELECT * FROM sync_runs
ORDER BY created_at DESC
LIMIT ? OFFSET ?;
