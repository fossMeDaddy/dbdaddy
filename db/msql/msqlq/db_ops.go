package msqlq

import "fmt"

func QGetExistingDbs() string {
	return fmt.Sprintf(
		`SHOW DATABASES WHERE %s NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')`,
		"`Database`",
	)
}

func QCheckDbExists(dbname string) string {
	return fmt.Sprintf(
		`SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'`,
		dbname,
	)
}
