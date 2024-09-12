package initCmd

import (
	"fmt"
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
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
	cwd, cwdIsProject, err := libUtils.CwdIsProject()
	if err != nil {
		cmd.Println(err)
		return
	}

	if cwdIsProject {
		cmd.PrintErrln("CWD is already a project, cannot re-initialize.")
		return
	} else if _, dirname := path.Split(cwd); dirname == constants.SelfConfigDirName {
		cmd.PrintErrln(
			fmt.Sprintf(
				"Can't create a project in %s, please choose any other directory",
				constants.SelfConfigDirName,
			),
		)
		return
	} else if configPath, _ := libUtils.FindConfigDirPath(); configPath == path.Join(cwd, constants.SelfConfigDirName) {
		cmd.PrintErrln(
			fmt.Sprintf(
				"A local '%s' configuration exists in this directory, please choose another directory",
				constants.SelfConfigDirName,
			),
		)
		return
	}

	v := viper.New()
	existingConfigFilePath, _ := libUtils.FindConfigFilePath()
	if err := lib.ReadConfig(v, existingConfigFilePath); err != nil {
		cmd.PrintErrln("unexpected error occured!")
		cmd.PrintErrln(err)
		return
	}
	configFilePath := path.Join(cwd, constants.SelfConfigFileName)
	if err := v.WriteConfigAs(configFilePath); err != nil {
		cmd.PrintErrln("unexpected error occured!")
		cmd.PrintErrln(err)
		return
	}

	_, migErr := libUtils.EnsureDirExists(path.Join(constants.MigDirName))
	_, schemaErr := libUtils.EnsureDirExists(path.Join(constants.SchemaDirName))
	_, scriptsErr := libUtils.EnsureDirExists(path.Join(constants.ScriptsDirName))
	if migErr != nil {
		cmd.PrintErrln(migErr)
	}
	if schemaErr != nil {
		cmd.PrintErrln(schemaErr)
	}
	if scriptsErr != nil {
		cmd.PrintErrln(scriptsErr)
	}
	if !slices.Contains([]error{migErr, schemaErr, scriptsErr}, nil) {
		return
	}

	cmd.Println(fmt.Sprintf("Opening config file at '%s' via vim...", configFilePath))
	libUtils.OpenFileInEditor(configFilePath)

	cmd.Println(fmt.Sprintf("Created project at: %s", cwd))
	cmd.Println("run 'help init' to know more about each of the newly created directories.")
	cmd.Println("happy databasing!")
}

func Init() *cobra.Command {
	return cmd
}
