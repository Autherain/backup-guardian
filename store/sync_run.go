package store

import (
	"context"
	"database/sql"

	"github.com/eva01/backup-guardian/domain"
	"github.com/eva01/backup-guardian/internal/errors"
	"github.com/eva01/backup-guardian/store/sqlc"
)

type syncRunsStore struct {
	baseStore *Store
}

var _ domain.SyncRunsReadWriter = (*syncRunsStore)(nil)

func (s *syncRunsStore) CreateSyncRun(run *domain.SyncRun) (*domain.SyncRun, error) {
	if err := run.Validate(); err != nil {
		return nil, err
	}

	q := sqlc.New(s.baseStore.db)

	startedAt := sql.NullTime{Time: run.StartedAt, Valid: !run.StartedAt.IsZero()}
	row, err := q.CreateSyncRun(context.Background(), sqlc.CreateSyncRunParams{
		ID:        run.ID,
		JobName:   run.JobName,
		StartedAt: startedAt,
	})
	if err != nil {
		return nil, errors.MapSQLError(err)
	}

	return mapSQLcToSyncRun(&row), nil
}

func (s *syncRunsStore) UpdateSyncRun(run *domain.SyncRun) error {
	if err := run.Validate(); err != nil {
		return err
	}

	q := sqlc.New(s.baseStore.db)

	finishedAt := sql.NullTime{Time: run.FinishedAt, Valid: !run.FinishedAt.IsZero()}
	var errMsg sql.NullString
	if run.ErrorMessage != "" {
		errMsg = sql.NullString{String: run.ErrorMessage, Valid: true}
	}
	filesTransferred := sql.NullInt64{Int64: run.FilesTransferred, Valid: true}
	bytesTransferred := sql.NullInt64{Int64: run.BytesTransferred, Valid: true}

	err := q.UpdateSyncRun(context.Background(), sqlc.UpdateSyncRunParams{
		Status:           run.Status,
		FinishedAt:       finishedAt,
		ErrorMessage:     errMsg,
		FilesTransferred: filesTransferred,
		BytesTransferred: bytesTransferred,
		ID:               run.ID,
	})

	return errors.MapSQLError(err)
}

func (s *syncRunsStore) GetSyncRun(selector *domain.SyncRunSelector) (*domain.SyncRun, error) {
	q := sqlc.New(s.baseStore.db)

	row, err := q.GetSyncRun(context.Background(), selector.ID)
	if err != nil {
		return nil, errors.MapSQLError(err)
	}

	return mapSQLcToSyncRun(&row), nil
}

func (s *syncRunsStore) ListSyncRuns(selector *domain.SyncRunsSelector) ([]*domain.SyncRun, error) {
	q := sqlc.New(s.baseStore.db)

	limit := int64(50)
	if selector.Limit > 0 {
		limit = int64(selector.Limit)
	}
	offset := int64(selector.Offset)

	rows, err := q.ListSyncRuns(context.Background(), sqlc.ListSyncRunsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, errors.MapSQLError(err)
	}

	result := make([]*domain.SyncRun, len(rows))
	for i := range rows {
		result[i] = mapSQLcToSyncRun(&rows[i])
	}

	return result, nil
}

func mapSQLcToSyncRun(row *sqlc.SyncRun) *domain.SyncRun {
	run := &domain.SyncRun{
		ID:        row.ID,
		JobName:   row.JobName,
		Status:    row.Status,
		CreatedAt: row.CreatedAt,
	}

	if row.StartedAt.Valid {
		run.StartedAt = row.StartedAt.Time
	}
	if row.FinishedAt.Valid {
		run.FinishedAt = row.FinishedAt.Time
	}
	if row.ErrorMessage.Valid {
		run.ErrorMessage = row.ErrorMessage.String
	}
	if row.FilesTransferred.Valid {
		run.FilesTransferred = row.FilesTransferred.Int64
	}
	if row.BytesTransferred.Valid {
		run.BytesTransferred = row.BytesTransferred.Int64
	}

	return run
}
