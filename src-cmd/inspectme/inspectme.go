package inspectMeCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	"dbdaddy/middlewares"
	"fmt"
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
	Use:     "inspectmedaddy",
	Aliases: []string{"inspectme"},
	Short:   "prints the schema of a selected table",
	Run:     cmdRunFn,
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
			dbStrTables = append(dbStrTables, fmt.Sprintf("%s.%s", dbTable.Schema, dbTable.Name))
		}

		if showAll {
			selectedTables = dbStrTables
		} else {
			prompt := promptui.Select{
				Label: "Choose table to display schema of",
				Items: slices.Concat(dbStrTables),
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

		for _, tablename := range selectedTables {
			tmp := strings.Split(tablename, ".")
			schema := tmp[0]
			table := tmp[1]

			tableSchema, err := db_int.GetTableSchema(currBranch, schema, table)
			if err != nil {
				return fmt.Errorf("Unexpected error occured while fetching table schema for %s\n"+err.Error(), tablename)
			}

			nColPadding := len(fmt.Sprintf("%d", len(tableSchema.Columns)))
			colNamePadding := 0
			colDefaultPadding := 0
			colDataTypePadding := 0
			for _, col := range tableSchema.Columns {
				colNamePadding = max(colNamePadding, len(col.Name))
				colDefaultPadding = max(colDefaultPadding, len(col.Default))
				colDataTypePadding = max(colDataTypePadding, len(col.DataType))
			}

			cmd.Println()
			cmd.Printf("TABLE: %s\n", tablename)
			for i, col := range tableSchema.Columns {
				cmd.Printf(
					"%0*d - %-*s %-*s DEFAULT %-*s NULLABLE %t\n",
					nColPadding, i+1,
					colNamePadding, col.Name,
					colDataTypePadding, col.DataType,
					colDefaultPadding, col.Default,
					col.Nullable,
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
