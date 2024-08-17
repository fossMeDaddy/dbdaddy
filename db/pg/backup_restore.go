package pg

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func DumpDb(outputFilePath string, v *viper.Viper, onlySchema bool) error {
	args := []string{"--clean",
		"--file=" + outputFilePath,
		"--format=directory",
		"--username=" + v.GetString(constants.DbConfigUserKey),
		"--host=" + v.GetString(constants.DbConfigHostKey),
		"--port=" + v.GetString(constants.DbConfigPortKey),
		"--no-password",
		v.GetString(constants.DbConfigCurrentBranchKey),
	}
	if onlySchema {
		args = append(args, "--schema-only")
	}
	osCmd := exec.Command("pg_dump", args...)

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
