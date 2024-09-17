package pg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/types"
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

	connConfig := types.ConnConfig{}
	if err := v.UnmarshalKey(constants.DbConfigConnSubkey, &connConfig); err != nil {
		return err
	}
	connConfig.Database = dbname

	// run restore command
	osCmd := exec.Command(
		"pg_restore",
		"--clean",
		"--if-exists",
		fmt.Sprintf("--dbname=%s", db.GetPgConnUriFromConnConfig(connConfig)),
		dumpFilePath,
	)

	osCmd.Stderr = os.Stderr

	err := osCmd.Run()
	return err
}
