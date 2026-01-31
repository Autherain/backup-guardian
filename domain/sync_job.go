package domain

// SyncJob represents a sync job configuration (source, destination).
// Initially configured via .env; can be extended to DB later.
type SyncJob struct {
	Name        string
	Source      string
	Destination string
}
