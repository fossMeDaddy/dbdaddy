package pg

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/types"
)

func DumpDb(outputFilePath string, connConfig types.ConnConfig, onlySchema bool) error {
	args := []string{"--clean",
		"--file=" + outputFilePath,
		"--format=directory",
		// fmt.Sprintf("--dbname=%s", db.GetPgConnUriFromConnConfig(connConfig)),
		"--username=" + connConfig.User,
		"--host=" + connConfig.Host,
		"--port=" + connConfig.Port,
		"--no-password",
		connConfig.Database,
	}
	if onlySchema {
		args = append(args, "--schema-only")
	}
	osCmd := exec.Command("pg_dump", args...)

	osCmd.Env = append(os.Environ(), "PGPASSWORD="+connConfig.Password)

	osCmd.Stdout = os.Stdout
	osCmd.Stderr = os.Stderr

	err := osCmd.Run()
	return err
}

// creates db & restores its content from given dump
func RestoreDb(connConfig types.ConnConfig, dumpFilePath string, override bool) error {
	if DbExists(connConfig.Database) {
		if override {
			if err := DeleteDb(connConfig.Database); err != nil {
				return err
			}

			if err := CreateDb(connConfig.Database); err != nil {
				return err
			}
		}
	} else {
		if err := CreateDb(connConfig.Database); err != nil {
			return err
		}
	}

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
