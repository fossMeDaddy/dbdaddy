package dumpMeCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/errs"
	"dbdaddy/lib"
	"dbdaddy/middlewares"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	useGlobalFile = false
	// outputPathArg = ""
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "dumpmedaddy",
	Short: "Takes a dump of the current database branch, this dump file can be later used to restore the data",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	configFilePath, _ := lib.FindConfigFilePath()
	if useGlobalFile {
		configFilePath = constants.GetGlobalConfigPath()
	}

	v := viper.New()
	lib.ReadConfig(v, configFilePath)

	outputFilePath := path.Join(lib.GetDriverDumpDir(configFilePath), constants.GetDumpFileName(v.GetString(constants.DbConfigCurrentBranchKey)))

	if err := db_int.DumpDb(outputFilePath, v, false); err != nil {
		if err.Error() == errs.PG_DUMP_NOT_FOUND {
			cmd.Println("Hey! we noticed you don't have 'pg_dump', then you also probably won't have 'pg_restore', we use these tools internally to perform dumps & restores... please install these tools in your OS before proceeding.")
		} else {
			cmd.PrintErrln("Unexpected error occured:", err)
		}

		return
	}

	cmd.Println("\nDumped successfully! output file:", outputFilePath)
}

func Init() *cobra.Command {
	cmd.Flags().BoolVarP(&useGlobalFile, "global", "g", false, "explicitly use global config file creds to connect to db")
	// cmd.Flags().StringVarP(&outputPathArg, "output", "o", "", "define output dump file path explicitly, by default i'll store it with me at your nearest '.dbdaddy' directory")

	return cmd
}
