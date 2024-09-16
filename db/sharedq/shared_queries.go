package sharedq

import "fmt"

func QDeleteDb(dbname string) string {
	return fmt.Sprintf(
		`drop database %s`,
		dbname,
	)
}

func QCreateNewDb(dbname string) string {
	return fmt.Sprintf(
		`CREATE DATABASE %s`,
		dbname,
	)
}
