package execCmd

import (
	"bufio"
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	"dbdaddy/libUtils"
	"dbdaddy/middlewares"
	"fmt"
	"math"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var (
	outFileFlag    string
	queryFlag      string
	noTruncateFlag bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "exec",
	Short: "Run SQL queries directly from CLI and get query outputs in CSV format",
	Run:   cmdRunFn,
}

func runQuery(cmd *cobra.Command, query string) error {
	var wg sync.WaitGroup

	if len(query) == 0 {
		return nil
	}

	query = strings.ToLower(strings.TrimRight(query, ";"))

	if !strings.Contains(query, "limit") {
		query += " limit 50"
		cmd.Println()
		cmd.Println("No 'LIMIT' clause found in query, limiting to 50 items.")
	}

	results, err := db_int.GetRows(query)
	if err != nil {
		return err
	}

	if len(results.Columns) == 0 {
		cmd.Println("Query ran successfully.")
		cmd.Println()
		return nil
	}

	cmd.Println("ROW COUNT:", results.RowCount)
	cmd.Println()

	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		cmd.Println("ERROR: couldn't get the terminal size, terminal table printing might not be ideal...")
		w = math.MaxInt32
	}

	wg.Add(1)
	go (func() {
		defer wg.Done()

		if len(outFileFlag) > 0 {
			csvFile := lib.GetCsvString(results)
			outFilePath, err := libUtils.GetAbsolutePathFor(outFileFlag)
			if err != nil {
				cmd.Println("unexpected error occured while writing CSV file", err)
				return
			}

			fileErr := os.WriteFile(outFilePath, []byte(csvFile), 0644)
			if fileErr != nil {
				cmd.Println("unexpected error occured while writing CSV file", fileErr)
				return
			}

			cmd.Println("CSV output saved to file:", outFilePath)
		}
	})()

	formattedOutput := lib.GetFormattedColumns(results)
	outputW := strings.Index(formattedOutput, fmt.Sprintln()) + 1

	if w >= outputW {
		cmd.Println(formattedOutput)
	} else {
		tmpDir, tmpDirErr := lib.FindTmpDirPath()
		if tmpDirErr != nil {
			return err
		}

		tmpFilePath := path.Join(tmpDir, constants.TextQueryOutput)
		err := os.WriteFile(tmpFilePath, []byte(formattedOutput), 0644)
		if err != nil {
			return err
		}

		cmd.Println("Formatted output too long to display here. See temp. file:", tmpFilePath)
	}

	wg.Wait()
	return nil
}

func run(cmd *cobra.Command, args []string) {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)
	err := lib.SwitchDB(viper.GetViper(), currBranch, func() error {
		if queryFlag != "" {
			return runQuery(cmd, queryFlag)
		}

		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("> ")

			query, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return err
			}

			query = strings.Trim(query, " ")
			query = strings.Trim(query, "\n")
			query = strings.ToLower(query)

			if query == ".tables" {
				tables, err := db_int.ListTablesInDb()
				if err != nil {
					cmd.PrintErrln("Unexpectedly, coudln't read list of tables", err)
					return err
				}

				cmd.Println("Listing tables in database:", currBranch)
				for i, table := range tables {
					cmd.Printf("%d. %s.%s\n", i+1, table.Schema, table.Name)
				}
				cmd.Println()

				continue
			}

			if query == ".exit" {
				return nil
			}

			qErr := runQuery(cmd, query)
			if qErr != nil {
				cmd.PrintErrln(qErr)
			}
		}
	})
	if err != nil {
		cmd.PrintErrln("error occured:", err)
	}
}

func Init() *cobra.Command {
	cmd.Flags().StringVarP(&queryFlag, "query", "q", "", "Enter text here to execute a one-off query")
	cmd.Flags().StringVarP(&outFileFlag, "output", "o", "", "Specify CSV-formatted output file path here. supplied along with '-q'")
	cmd.Flags().BoolVar(&noTruncateFlag, "no-truncate", false, "In table formatted output, row items are truncated by default, use this to disable it.")

	return cmd
}
