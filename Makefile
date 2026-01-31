# backup-guardian Makefile
#
# Schéma : sql/schema/schema.sql (maintenu à jour avec les migrations)
# Migrations : migrations/ (goose)
#
# Workflow (ajouter une table) :
#   1. Créer migration : migrations/00002_ma_table.sql (-- +goose Up/Down)
#   2. Mettre à jour sql/schema/schema.sql avec le nouvel état
#   3. Ajouter requêtes : sql/queries/ma_table.sql
#   4. make sqlc        → génère store/sqlc/
#   5. Adapter store/ et domain/ — les migrations s'appliquent au run du runner
#
.PHONY: sqlc sqlc-vet test generate docker-build docker-up all

# --- sqlc (génère le code à partir de sql/schema/) ---
sqlc:
	sqlc generate

# Vérifie que les requêtes et le schéma sont valides
sqlc-vet:
	sqlc vet

# --- Tests ---
test:
	go test ./...

# Les migrations sont appliquées au démarrage du runner (cmd/runner).

# --- Docker ---
docker-build:
	docker compose build

docker-up:
	docker compose up -d

# --- Tout ---
# Après avoir ajouté une table : migration + schema + queries + sqlc
all: sqlc
