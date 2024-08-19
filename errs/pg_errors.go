package errs

import "fmt"

var ErrPgDumpCmdNotFound = fmt.Errorf("'pg_dump' was not found")
var ErrPsqlCmdNotFound = fmt.Errorf("'psql' was not found")

func PgDumpCmdNotFoundMsg() string {
	return "'pg_dump' was not found in PATH, please make sure 'pg_dump' is installed on your machine"
}

func PsqlCmdNotFoundMsg() string {
	return "'psql' was not found in PATH, please make sure 'psql' is installed on your machine"
}
