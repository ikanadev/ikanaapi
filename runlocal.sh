export PORT=4000
export DATABASE=postgresql://kevv:postgres@0.0.0.0:5432/ikanadb?sslmode=disable
export MIGRATIONS_URL=github://ikanadev/ikanaapi/db/migrations

go run *.go
