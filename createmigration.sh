# usage createmigration.sh my_migration_name

export DATABASE=postgresql://kevv:postgres@0.0.0.0:5432/ikanadb?sslmode=disable
migrate create -ext sql -dir db/migrations -seq $1
