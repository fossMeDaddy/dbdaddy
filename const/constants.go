package constants

const (
	// Self
	SelfDbName = "__daddys_home"

	//  Config keys
	DbConfigDriverKey = "connection.driver"
	DbConfigHostKey   = "connection.host"
	DbConfigPortKey   = "connection.port"
	DbConfigDbNameKey = "connection.dbname"
	DbConfigUserKey   = "connection.user"
	DbConfigPassKey   = "connection.password"

	DbConfigCurrentBranchKey = "status.currentBranch"

	// Config possible driver values
	DbDriverPostgres = "postgres"
	DbDriverMySQL    = "mysql"
	DbDriverSqlite   = "sqlite"
)
