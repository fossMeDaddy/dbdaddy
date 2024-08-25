package constants

// Default SQL migration file indentation
const (
	MigFileIndentation = "    "
)

// collection of migrations for each db
const (
	MigDirName = "migrations"
)

// migration dir files e.g. "migrations/dbname/version_string/*"
const (
	MigDirStateFile        = "state.json"
	MigDirActiveSignalFile = "active"
	MigDirInfoFile         = "info.md"
	MigDirUpSqlFile        = "up.sql"
	MigDirDownSqlFile      = "down.sql"
)
