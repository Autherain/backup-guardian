# backup-guardian

Project made to ensure i backup my GoogleDrive to others clouds providers because me paranoid.

## Setup

1. **rclone config** : Configure your remotes with `rclone config` (creates `~/.config/rclone/rclone.conf`). No need to install rclone CLI - this project uses the rclone Go library.

2. **Environment** : Copy `.env.example` to `.env` and set `BG_SYNC_SOURCE` (e.g. `gdrive:`), `BG_SYNC_DEST` (e.g. `s3:bucket/backups`). `BG_DATA_DIR` defaults to `.` (DB at `./backup-guardian.db`).

3. **Runner** : `make run` ou `go run ./cmd/runner`. Les migrations (goose) sont appliquées automatiquement au démarrage.

Voir le `Makefile` pour les commandes sqlc (`make sqlc`) et le workflow complet.

## Docker

1. Configure `BG_SYNC_SOURCE` and `BG_SYNC_DEST` in `docker-compose.yml`.
2. Start: `docker compose up -d` (les migrations sont appliquées au démarrage du runner).

Rclone config must exist at `~/.config/rclone` (mount point in compose).
