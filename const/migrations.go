package constants

// Migration action
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
