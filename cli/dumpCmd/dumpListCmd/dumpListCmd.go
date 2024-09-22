package dumpListCmd

import (
	"path"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/middlewares"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	useGlobalConfigFile bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lists available database dumps",
	Run:     cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	configFilePath, _ := libUtils.FindConfigFilePath()
	if useGlobalConfigFile {
		configFilePath = libUtils.GetGlobalConfigPath()
	}

	dumpDbGroups, err := lib.GetDbGroupedDumpFiles(configFilePath)
	if err != nil {
		panic("Unexpected error occured!\n" + err.Error())
	}

	dumpDir := libUtils.GetDriverDumpDir(configFilePath, viper.GetString(constants.DbConfigDriverKey))
	cmd.Printf("Listing all database backups from directory: %s\n\n", dumpDir)
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
