package initCmd

import (
	"fmt"
	"os"
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/lib/migrationsLib"
	"github.com/fossmedaddy/dbdaddy/types"
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

var (
	forceFlag bool
)

var cmd = &cobra.Command{
	Use:     "init",
	Short:   "initialize a project structure in CWD to work like a sophisticated 10x dev. providing a connection uri is optional.",
	Example: "init postgresql://user:pwd@localhost:5432/dbname",
	Long:    cmdManual,
	Run:     run,
	Args:    cobra.MaximumNArgs(1),
}

func run(cmd *cobra.Command, args []string) {
	if forceFlag {
		os.RemoveAll(constants.MigDirName)
		os.RemoveAll(constants.SchemaDirName)
		os.RemoveAll(constants.ScriptsDirName)
		os.Remove(constants.SelfConfigFileName)
	}

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
				"can't create a project in %s, please choose any other directory",
				constants.SelfConfigDirName,
			),
		)
		return
	}

	projectConfigFilePath := path.Join(cwd, constants.SelfConfigFileName)

	if len(args) > 0 {
		connUri := args[0]

		connConfig, err := libUtils.GetConnConfigFromUri(connUri)
		if err != nil {
			cmd.PrintErrln("error occured while parsing connection uri")
			cmd.PrintErrln(err)
			return
		}

		if _, err := db.ConnectDb(connConfig); err != nil {
			cmd.PrintErrln(err)
			return
		}

		v := viper.New()
		lib.InitConfigFile(v, cwd, false)
		v.Set(constants.DbConfigConnKey, connConfig)
		v.Set(constants.DbConfigCurrentBranchKey, connConfig.Database)

		if err := v.WriteConfigAs(projectConfigFilePath); err != nil {
			cmd.PrintErrln("unexpected error occured")
			return
		}
	} else {
		v := viper.New()

		if configDirPath, _ := libUtils.FindConfigDirPath(); configDirPath == path.Join(cwd, constants.SelfConfigDirName) {
			configFilePath := path.Join(configDirPath, constants.SelfConfigFileName)
			if err := lib.ReadConfig(v, configFilePath); err != nil {
				cmd.PrintErrln("unexpected error occured")
				cmd.PrintErrln(err)
				return
			}

			if err := v.WriteConfigAs(projectConfigFilePath); err != nil {
				cmd.PrintErrln("unexpected error occured")
				cmd.PrintErrln(err)
				return
			}

			os.Remove(configFilePath)

			var (
				migMvErr, scriptsMvErr, schemaMvErr error
			)

			migDirPath := path.Join(configDirPath, constants.MigDirName)
			if libUtils.Exists(migDirPath) {
				migMvErr = os.Rename(migDirPath, path.Join(cwd, constants.MigDirName))
			}

			scriptsDirPath := path.Join(configDirPath, constants.ScriptsDirName)
			if libUtils.Exists(scriptsDirPath) {
				scriptsMvErr = os.Rename(scriptsDirPath, path.Join(cwd, constants.ScriptsDirName))
			}

			schemaDirPath := path.Join(configDirPath, constants.SchemaDirName)
			if libUtils.Exists(schemaDirPath) {
				schemaMvErr = os.Rename(schemaDirPath, path.Join(cwd, constants.SchemaDirName))
			}

			if migMvErr == nil {
				os.RemoveAll(path.Join(migDirPath))
			}
			if scriptsMvErr == nil {
				os.RemoveAll(scriptsDirPath)
			}
			if schemaMvErr == nil {
				os.RemoveAll(schemaDirPath)
			}

			if migMvErr == nil && scriptsMvErr == nil && schemaMvErr == nil {
				os.Remove(configDirPath)
			}
		} else {
			if err := lib.InitConfigFile(v, cwd, true); err != nil {
				cmd.PrintErrln("error occured while writing config file")
				cmd.PrintErrln(err)
				return
			}
		}
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

	if len(args) == 0 {
		cmd.Println(fmt.Sprintf("opening config file at '%s' via vim...", projectConfigFilePath))
		libUtils.OpenFileInEditor(projectConfigFilePath)
	}

	if globals.DB != nil {
		dbSchema, schemaErr := db_int.GetDbSchema()
		if schemaErr != nil {
			cmd.PrintErrln("unexpected error while trying to instrospect your database")
			cmd.PrintErrln(schemaErr)
			return
		}

		changes := migrationsLib.DiffDBSchema(dbSchema, &types.DbSchema{})
		schemaSql, schemaErr := migrationsLib.GetSQLFromDiffChanges(changes)
		if schemaErr != nil {
			cmd.PrintErrln("unexpected error occured")
			cmd.PrintErrln(schemaErr)
			return
		}

		schemaFilePath := path.Join(cwd, constants.SchemaDirName, "schema.sql")
		if err := os.WriteFile(schemaFilePath, []byte(schemaSql), 0644); err != nil {
			cmd.PrintErrln("error occured while introspecting and writing to schema")
			cmd.PrintErrln(err)
			return
		}

		cmd.Println(fmt.Sprintf("database schema introspected successfully, output written to: %s", schemaFilePath))
	}

	cmd.Println(fmt.Sprintf("Created project at: %s", cwd))
	cmd.Println("run 'help init' to know more about each of the newly created directories.")
	cmd.Println("happy databasing!")
}

func Init() *cobra.Command {
	// flags here
	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "remove project files from CWD (if any) and re-initialize project")

	return cmd
}
