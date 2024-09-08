package inspectMeCmd

import (
	"dbdaddy/constants"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	"dbdaddy/lib/libUtils"
	"dbdaddy/middlewares"
	"dbdaddy/types"
	"fmt"
	"maps"
	"slices"
	"strings"

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

	err := lib.SwitchDB(viper.GetViper(), currBranch, func() error {
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
			return fmt.Errorf("unexpected error occured while fetching db schema")
		}

		tablesViewsConcat := map[string]*types.TableSchema{}
		maps.Copy(tablesViewsConcat, dbSchema.Tables)
		maps.Copy(tablesViewsConcat, dbSchema.Views)

		for tableid, tableSchema := range tablesViewsConcat {
			if !slices.Contains(selectedTables, tableid) {
				continue
			}

			isView := dbSchema.Views[tableid] != nil

			// col constraint mapping
			pKeyMapping := map[string]bool{}
			colConMapping := map[string][]*types.DbConstraint{}
			for _, con := range tableSchema.Constraints {
				if con.Type == "p" {
					pKeyMapping[con.ColName] = true
				}
				colConMapping[con.ColName] = append(colConMapping[con.ColName], con)
			}

			nColPadding := len(fmt.Sprintf("%d", len(tableSchema.Columns)))
			colNamePadding := 0
			colDefaultPadding := 0
			colDataTypePadding := 0
			colNullablePadding := 5
			for _, col := range tableSchema.Columns {
				colNamePadding = max(colNamePadding, len(getColName(col.Name, pKeyMapping[col.Name])))
				colDefaultPadding = max(colDefaultPadding, len(col.Default))
				colDataTypePadding = max(colDataTypePadding, len(col.DataType))
			}

			cmd.Println()
			cmd.Printf("TABLE: %s\n", tableid)
			if isView {
				cmd.Println("INFO: TABLE IS A VIEW")
			}

			for i, col := range tableSchema.Columns {
				cmd.Printf(
					"%0*d - %-*s %-*s DEFAULT `%-*s` NULLABLE:%-*t %s\n",
					nColPadding, i+1,
					colNamePadding, getColName(col.Name, pKeyMapping[col.Name]),
					colDataTypePadding, col.DataType,
					colDefaultPadding, col.Default,
					colNullablePadding, col.Nullable,
					getConstraintString(colConMapping[col.Name]),
				)
			}
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
