package constants

const (
	DbConfigDefaultPostgresUrl = "postgresql://postgres:postgres@localhost:5342/postgres"

	// Self
	SelfGlobalDirName = ".dbdaddy"
	SelfDbName        = "__daddys_home"

	//  Config keys
	DbConfigDriverKey        = "connection.driver"
	DbConfigHostKey          = "connection.host"
	DbConfigPortKey          = "connection.port"
	DbConfigDbNameKey        = "connection.dbname"
	DbConfigUserKey          = "connection.user"
	DbConfigPassKey          = "connection.password"
	DbConfigCurrentBranchKey = "status.currentBranch"

	// Config possible driver values
	DbDriverPostgres = "postgres"
	DbDriverMySQL    = "mysql"
	DbDriverSqlite   = "sqlite"

	// dumps
	PgDumpDir     = "pg_dumps"
	MySqlDumpDir  = "mysql_dumps"
	SqliteDumpDir = "sqlite_dumps"

	// migrations
	MigrationDir = "migrations"

	// tmp
	TmpDir          = "tmp"
	TextQueryOutput = "query.out"
	CSVQueryOutput  = "query.csv"
)

var DriverDumpDirNames = map[string]string{
	DbDriverPostgres: PgDumpDir,
	DbDriverMySQL:    MySqlDumpDir,
	DbDriverSqlite:   SqliteDumpDir,
}

var SupportedDrivers = []string{DbDriverPostgres}
