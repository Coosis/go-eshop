dev:
	podman compose up -d

reset-db:
	podman container exec go-eshop-db-1 psql --username "postgres" -f /var/lib/postgresql/data/sql/schema.sql

.PHONY: sqlc
sqlc:
	sqlc generate
