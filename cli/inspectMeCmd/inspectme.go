package inspectMeCmd

import (
	"fmt"
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/middlewares"
	"github.com/fossmedaddy/dbdaddy/sqlwriter"
	"github.com/fossmedaddy/dbdaddy/types"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	showAll bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:     "inspect",
	Aliases: []string{"inspectme"},
	Short:   "prints the schema of a selected table in current database",
	Run:     cmdRunFn,
}

func getColName(name string, pk bool) string {
	if pk {
		return fmt.Sprintf("%s (Primary Key)", name)
	}

	return name
}

func getConstraintString(constraints []*types.DbConstraint) string {
	conStr := ""
	if len(constraints) == 0 {
		return conStr
	}

	for _, con := range constraints {
		switch con.Type {
		case "f":
			conStr += fmt.Sprintf("FK(%s.%s.%s)", con.FTableSchema, con.FTableName, con.FColName)
		case "u":
			conStr += "UNIQUE"
		case "c":
			conStr += strings.ReplaceAll(con.Syntax, fmt.Sprintln(), " ")
		}
	}

	return conStr
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	err := lib.TmpSwitchDB(currBranch, func() error {
		selectedTables := []string{}

		dbTables, err := db_int.ListTablesInDb()
		if err != nil {
			return fmt.Errorf("unexpected error occured while fetching tables from database '%s'\n%s", currBranch, err.Error())
		}
		dbStrTables := []string{}
		for _, dbTable := range dbTables {
			dbStrTables = append(dbStrTables, libUtils.GetTableId(dbTable.Schema, dbTable.Name))
		}

		if showAll {
			selectedTables = dbStrTables
		} else {
			prompt := promptui.Select{
				Label: "Choose table to display schema of",
				Items: dbStrTables,
				Searcher: func(input string, index int) bool {
					return strings.Contains(dbStrTables[index], strings.ToLower(strings.Trim(input, " ")))
				},
				StartInSearchMode: true,
			}

			_, result, err := prompt.Run()
			if err != nil {
				return err
			}

			selectedTables = append(selectedTables, result)
		}

		dbSchema, err := db_int.GetDbSchema(currBranch)
		if err != nil {
			return err
		}

		viewsPrintBuf := ""
		for _, tableid := range selectedTables {
			isView := false
			tableSchema := dbSchema.Tables[tableid]
			if tableSchema == nil {
				tableSchema = dbSchema.Views[tableid]
				isView = true
			}

			var getSqlFn func(*types.TableSchema) (string, error)
			if !isView {
				getSqlFn = sqlwriter.GetCreateTableSQL
			} else {
				getSqlFn = sqlwriter.GetCreateViewSQL
			}
			defSql, err := getSqlFn(tableSchema)
			if err != nil {
				return err
			}

			tableConSql := ""
			if !isView {
				for _, con := range tableSchema.Constraints {
					if atConSql, err := sqlwriter.GetATCreateConstraintSQL(tableid, con); err != nil {
						return err
					} else {
						tableConSql += atConSql
					}
				}
			}

			if isView {
				viewsPrintBuf += fmt.Sprintln(fmt.Sprintf("--- VIEW: %s", tableid))
				viewsPrintBuf += fmt.Sprintln(defSql + tableConSql)
				viewsPrintBuf += fmt.Sprintln()
			} else {
				cmd.Println(fmt.Sprintf("--- TABLE: %s", tableid))
				cmd.Println(defSql + tableConSql)
				cmd.Println()
			}
		}

		if len(viewsPrintBuf) > 0 {
			cmd.Print(viewsPrintBuf)
		}

		return nil
	})
	if err != nil {
		cmd.PrintErrln(err)
	}
}

func Init() *cobra.Command {
	cmd.Flags().BoolVar(&showAll, "all", false, "print schema for all tables")

	return cmd
}
