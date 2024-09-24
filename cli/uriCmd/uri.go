package uriCmd

import (
	"strings"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	useGlobalConfigFlag bool
)

var cmd = &cobra.Command{
	Use:     "uri",
	Short:   "replace the connection credentials in dbdaddy config with a newly provided db connection uri",
	Example: "uri postgresql://user:pwd@127.0.0.1:5432/dbname",
	Args:    cobra.ExactArgs(1),
	Run:     run,
}

func run(cmd *cobra.Command, args []string) {
	uri := args[0]
	uri = strings.Trim(uri, " ")

	connConfig := types.ConnConfig{}
	if cc, err := libUtils.GetConnConfigFromUri(uri); err != nil {
		cmd.PrintErrln(err)
		return
	} else {
		connConfig = cc
	}

	if _, err := db.ConnectDb(connConfig); err != nil {
		cmd.PrintErrln("error occured while trying to connect to database, please check your connection uri...")
		cmd.PrintErrln(err)
		return
	}

	configFilePath, err := libUtils.FindConfigFilePath()
	if err != nil {
		cmd.PrintErrln("error occured while trying to find config file in your device...")
		cmd.PrintErrln(err)
		return
	}

	if useGlobalConfigFlag {
		configFilePath = libUtils.GetGlobalConfigPath()
	}

	viper.Set(constants.DbConfigConnSubkey, connConfig)
	viper.Set(constants.DbConfigCurrentBranchKey, connConfig.Database)

	if err := viper.WriteConfigAs(configFilePath); err != nil {
		cmd.PrintErrln("unexpected error occured while writing to config file")
		cmd.PrintErrln(err)
		return
	}

	cmd.Println("db connection credentials changed successfully in config")
}

func Init() *cobra.Command {
	// flags here
	cmd.Flags().BoolVarP(&useGlobalConfigFlag, "global", "g", false, "explicitly use global config file")

	return cmd
}
