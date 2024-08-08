package msqlq

import "fmt"

func QGetExistingDbs() string {
	return fmt.Sprintf(
		`SHOW DATABASES WHERE %s NOT IN ('information_schema', 'mysql', 'performance_schema', 'sys')`,
		"`Database`",
	)
}
