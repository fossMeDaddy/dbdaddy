package schemaPushCmd

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/middlewares"
	"github.com/fossmedaddy/dbdaddy/sqlparser"
	"github.com/spf13/cobra"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "push",
	Short: "push & override database with changes done in schema folder",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	cwd, cwdIsProject, err := libUtils.CwdIsProject()
	if err != nil {
		cmd.PrintErrln("unexpected error occured while reading CWD")
		cmd.PrintErrln(err)
		return
	}
	if !cwdIsProject {
		cmd.PrintErrln("CWD is not a project! initialize one with 'init' command")
		return
	}

	schemaSqlStr := ""
	schemaDirPath := path.Join(cwd, constants.SchemaDirName)
	filepath.WalkDir(schemaDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".sql") {
			return nil
		}

		if fileB, err := os.ReadFile(path); err != nil {
			return err
		} else {
			schemaSqlStr += string(fileB) + fmt.Sprintln()
		}

		return nil
	})

	_, parseErr := sqlparser.GetDbSchemaFromSQL(schemaSqlStr)
	if parseErr != nil {
		cmd.PrintErrln("unexpected error occured")
		cmd.PrintErrln(parseErr)
		return
	}
}

func Init() *cobra.Command {
	return cmd
}
