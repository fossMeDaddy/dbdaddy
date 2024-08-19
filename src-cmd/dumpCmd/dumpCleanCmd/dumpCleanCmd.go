package dumpCleanCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	"dbdaddy/libUtils"
	"dbdaddy/middlewares"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	useGlobalConfigFile bool
	cleanAll            bool
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleans old dump files, won't touch latest dumps",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	configFilePath, _ := libUtils.FindConfigFilePath()
	if useGlobalConfigFile {
		configFilePath = constants.GetGlobalConfigPath()
	}

	dumpDbGroupedFiles, err := lib.GetDbGroupedDumpFiles(configFilePath)
	if err != nil {
		panic("Unexpected error occured!\n" + err.Error())
	}

	for _, files := range dumpDbGroupedFiles {
		var delFiles []string
		if cleanAll {
			delFiles = files
		} else {
			delFiles = files[1:]
		}

		for _, file := range delFiles {
			if err := os.RemoveAll(file); err != nil {
				panic("Unexpected erro occured while removing backups!\n" + err.Error())
			}

			cmd.Println("Removed dump:", file)
		}
	}
}

func Init() *cobra.Command {
	cmd.Flags().BoolVar(&cleanAll, "all", false, fmt.Sprintf("delete ALL dump files in the '%s' config dir", constants.SelfGlobalDirName))
	cmd.Flags().BoolVarP(&useGlobalConfigFile, "global", "g", false, "explicitly use the global config file to list the dumps")

	return cmd
}
