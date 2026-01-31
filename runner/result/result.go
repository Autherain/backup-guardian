package result

import "time"

// RcloneResult holds the result of an rclone sync operation.
type RcloneResult struct {
	FilesTransferred int64
	BytesTransferred int64
	Duration         time.Duration
}
