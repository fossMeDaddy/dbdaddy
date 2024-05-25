package dumpListCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	"dbdaddy/middlewares"
	"path"

	"github.com/spf13/cobra"
)

var (
	useGlobalConfigFile bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "list",
	Short: "Lists available database dumps",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	configFilePath, _ := lib.FindConfigFilePath()
	if useGlobalConfigFile {
		configFilePath = constants.GetGlobalConfigPath()
	}

	dumpDbGroups, err := lib.GetDbGroupedDumpFiles(configFilePath)
	if err != nil {
		panic("Unexpected error occured!\n" + err.Error())
	}

	cmd.Printf("Listing all database backups from directory: %s\n\n", lib.GetDriverDumpDir(configFilePath))
	if len(dumpDbGroups) == 0 {
		cmd.Println("No backups found.")
	}

	for dbname, dumpFiles := range dumpDbGroups {
		cmd.Printf("%s:\n", dbname)
		for i, file := range dumpFiles {
			_, filename := path.Split(file)
			cmd.Printf("  %d. %s\n", i+1, filename)
		}
		cmd.Println()
	}
}

func Init() *cobra.Command {
	cmd.Flags().BoolVarP(&useGlobalConfigFile, "global", "g", false, "explicitly use the global config file to list the dumps")

	return cmd
}
