package initCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/spf13/cobra"
)

var cmdManual = fmt.Sprintf(`
init command creates a bunch of folders:
1. %s
2. %s
3. %s

%s:
contains all the generated migrations in the database selected in this project.

%s:
contains .sql files written by you to be run on the selected database in this project.
let's take an example:
	folder structure:
	|- %s/
	|- %s/
		|- aggregations/
			|- aggregate-users.sql
		|- test.sql
	|- %s/

here, you'll run the 'aggregate-users.sql' file by running: 'run aggregations/aggregate-users'
and the 'test.sql' file by running: 'run test'

%s:
contains your database schema written by you on the selected database in this project.
each table gets its own .sql file and in case of postgres, each schema gets its own folder name,
that way, you can define multiple tables inside a schema and multiple schemas inside your database.

note:
for sqlite & mysql, there are no folders in the %s folder to represent tables directly.

let's take an example of a postgres schema defined in this structure:
	folder structure:
	|- %s/
	|- %s/
	|- %s/
		|- auth/
			|- users.sql
		|- public/
			|- profiles.sql
			|- posts.sql

here, auth & public are schemas and the filenames like 'users', 'profiles', 'posts' represent
table names inside their respective schemas in postgres.

`,
	constants.MigDirName,
	constants.ScriptsDirName,
	constants.SchemaDirName,

	constants.MigDirName,

	constants.ScriptsDirName,
	constants.MigDirName,
	constants.ScriptsDirName,
	constants.SchemaDirName,

	constants.SchemaDirName,
	constants.SchemaDirName,
	constants.MigDirName,
	constants.ScriptsDirName,
	constants.SchemaDirName,
)

var cmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a project structure in CWD to work like a sophisticated 10x dev",
	Long:  cmdManual,
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	_, cwdIsProject, err := libUtils.CwdIsProject()
	if err != nil {
		cmd.Println(err)
		return
	}
	if cwdIsProject {
		cmd.Println("CWD is already a project, cannot re-initialize.")
		return
	}

	// configFilePath, _ := libUtils.FindConfigFilePath()
	// copy from this path
	// write to cwd config json

	// create required directories (queries, migrations, schema, dbdaddy.config.json)
	// instruct to run "help init" to know more about each of the generated folders
}

func Init() *cobra.Command {
	// flags here

	return cmd
}
