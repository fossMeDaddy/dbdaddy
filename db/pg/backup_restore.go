package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/errs"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func DumpDb(outputFilePath string, v *viper.Viper) error {
	osCmd := exec.Command("pg_dump",
		"--clean",
		"--file="+outputFilePath,
		"--format=directory",
		"--username="+v.GetString(constants.DbConfigUserKey),
		"--host="+v.GetString(constants.DbConfigHostKey),
		"--port="+v.GetString(constants.DbConfigPortKey),
		"--no-password",
		v.GetString(constants.DbConfigCurrentBranchKey))

	osCmd.Env = append(os.Environ(), "PGPASSWORD="+v.GetString(constants.DbConfigPassKey))

	osCmd.Stdout = os.Stdout
	osCmd.Stderr = os.Stderr

	err := osCmd.Run()
	return err
}
func DumpDbOnlySchema(outputFilePath string, v *viper.Viper) error {
	if !PgDumpExists() {
		return fmt.Errorf(errs.PG_DUMP_NOT_FOUND)
	}

	osCmd := exec.Command("pg_dump",
		"--clean",
		"--schema-only",
		"--file="+outputFilePath,
		"--format=directory",
		"--verbose",
		"--username="+v.GetString(constants.DbConfigUserKey),
		"--host="+v.GetString(constants.DbConfigHostKey),
		"--port="+v.GetString(constants.DbConfigPortKey),
		"--no-password",
		v.GetString(constants.DbConfigCurrentBranchKey))

	osCmd.Env = append(os.Environ(), "PGPASSWORD="+v.GetString(constants.DbConfigPassKey))

	osCmd.Stdout = os.Stdout
	osCmd.Stderr = os.Stderr

	err := osCmd.Run()
	return err
}

// creates db & restores its content from given dump
func RestoreDb(dbname string, v *viper.Viper, dumpFilePath string, override bool) error {
	if DbExists(dbname) {
		if override {
			if err := DeleteDb(dbname); err != nil {
				return err
			}

			if err := CreateDb(dbname); err != nil {
				return err
			}
		}
	} else {
		if err := CreateDb(dbname); err != nil {
			return err
		}
	}

	// run restore command
	osCmd := exec.Command(
		"pg_restore",
		"--clean",
		"--if-exists",
		fmt.Sprintf("--dbname=%s", db.GetPgConnUriFromViper(v, dbname)),
		dumpFilePath,
	)

	osCmd.Stderr = os.Stderr

	err := osCmd.Run()
	return err
}

func RestoreDbOnlySchema(dbname string, v *viper.Viper, dumpFilePath string) error {
	if err := CreateDb(dbname); err != nil {
		return err
	}

	// run restore command
	osCmd := exec.Command(
		"pg_restore",
		"--verbose",
		"--clean",
		"--if-exists",
		fmt.Sprintf("--dbname=%s", db.GetPgConnUriFromViper(v, dbname)),
		dumpFilePath,
	)

	osCmd.Stderr = os.Stderr
	osCmd.Stdout = os.Stdout

	err := osCmd.Run()
	return err
}
