package constants

const (
	// Self
	SelfConfigDirName  = ".dbdaddy"
	SelfConfigFileName = "dbdaddy.config.json"
	SelfDbName         = "__daddys_home"

	//  Config keys
	DbConfigOriginsKey       = "origins"
	DbConfigConnSubkey       = "connection"
	DbConfigDriverKey        = DbConfigConnSubkey + ".driver"
	DbConfigHostKey          = DbConfigConnSubkey + ".host"
	DbConfigPortKey          = DbConfigConnSubkey + ".port"
	DbConfigDbNameKey        = DbConfigConnSubkey + ".dbname"
	DbConfigUserKey          = DbConfigConnSubkey + ".user"
	DbConfigPassKey          = DbConfigConnSubkey + ".password"
	DbConfigParamsKey        = DbConfigConnSubkey + ".params"
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
	SoftDeleteSuffix = "__DELETED"
)

var DriverDumpDirNames = map[string]string{
	DbDriverPostgres: PgDumpDir,
	DbDriverMySQL:    MySqlDumpDir,
	DbDriverSqlite:   SqliteDumpDir,
}

var SupportedDrivers = []string{DbDriverPostgres, DbDriverMySQL}
