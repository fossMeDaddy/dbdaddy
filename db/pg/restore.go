package pg

import (
	"dbdaddy/db"
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

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
		"--verbose",
		"--clean",
		"--if-exists",
		fmt.Sprintf("--dbname=%s", db.GetPgConnUriFromViper(dbname)),
		dumpFilePath,
	)

	osCmd.Stderr = os.Stderr
	osCmd.Stdout = os.Stdout

	err := osCmd.Run()
	return err
}
