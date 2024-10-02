package constants

const (
	// Self
	SelfConfigDirName  = ".dbdaddy"
	SelfConfigFileName = "dbdaddy.config.json"
	SelfDbName         = "__daddys_home"

	//  Config keys
	DbConfigOriginsKey       = "origins"
	DbConfigConnKey          = "connection"
	DbConfigShadowConnKey    = "tmp_db_connection"
	DbConfigCurrentBranchKey = "status.currentBranch"

	// Config possible driver values
	DbDriverPostgres = "postgres"
	DbDriverMySQL    = "mysql"
	DbDriverSqlite   = "sqlite"

	// dumps
	PgDumpDir     = "pg_dumps"
	MySqlDumpDir  = "mysql_dumps"
	SqliteDumpDir = "sqlite_dumps"

	// project
	ScriptsDirName = "scripts"
	SchemaDirName  = "schema"

	// tmp
	TmpDir          = "tmp"
	TextQueryOutput = "query.out"
	CSVQueryOutput  = "query.csv"

	// misc
	UpSqlScriptComment   = "--- UP SQL (APPLY THIS)"
	DownSqlScriptComment = "--- DOWN SQL (APPLY ONLY WHEN NEED TO REVERT)"
	SoftDeleteSuffix     = "__DELETED"
	ShadowDbPrefix       = "__shadow"
)

var DriverDumpDirNames = map[string]string{
	DbDriverPostgres: PgDumpDir,
	DbDriverMySQL:    MySqlDumpDir,
	DbDriverSqlite:   SqliteDumpDir,
}

var SupportedDrivers = []string{DbDriverPostgres, DbDriverMySQL}
