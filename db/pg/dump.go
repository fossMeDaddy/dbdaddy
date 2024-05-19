package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/errs"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func PgDumpExists() bool {
	_, err := exec.LookPath("pg_dump")
	return err == nil
}

func TakeADump(outputFilePath string, v *viper.Viper) error {
	if !PgDumpExists() {
		return fmt.Errorf(errs.PG_DUMP_NOT_FOUND)
	}

	osCmd := exec.Command("pg_dump",
		"--clean",
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
