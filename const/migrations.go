package constants

// Migration action
// always use the prefixes: "CREATE_" OR "DROP_" (underscore is necessary)
const (
	MigActionCreateTable = "CREATE_TABLE"
	MigActionDropTable   = "DROP_TABLE"

	MigActionCreateCol = "CREATE_COLUMN"
	MigActionDropCol   = "DROP_COLUMN"

	MigActionCreateType = "CREATE_TYPE"
	MigActionDropType   = "DROP_TYPE"

	MigActionCreateConstraint = "CREATE_CONSTRAINT"
	MigActionDropConstraint   = "DROP_CONSTRAINT"
)

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
