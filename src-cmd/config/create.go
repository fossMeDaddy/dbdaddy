package configCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createInGlobalNamespace = false
var overrideExisting = false

var createCmd = &cobra.Command{
	Use:   "create",
	Short: fmt.Sprintf("Create a new config file in the current working dir. Use '-g or --global' to create config file at %s", constants.GetGlobalConfigPath()),
	Run:   createCmdRun,
}

func createCmdRun(cmd *cobra.Command, args []string) {
	configWritePath := constants.GetLocalConfigPath()
	if createInGlobalNamespace {
		configWritePath = constants.GetGlobalConfigPath()
		if err := lib.EnsureDirExists(configWritePath); err != nil {
			panic("Unexpected error occured!\n" + err.Error())
		}
	}

	if overrideExisting {
		lib.InitConfigFile(viper.New())

		if err := viper.WriteConfigAs(configWritePath); err != nil {
			cmd.PrintErrln("Error occured while writing config file!\n" + err.Error())
			return
		}
	} else {
		if err := viper.SafeWriteConfigAs(configWritePath); err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
	}
}

func InitCreateCmd() *cobra.Command {
	createCmd.Flags().BoolVarP(&createInGlobalNamespace, "global", "g", false, fmt.Sprintf("write config file at '%s' for global use, writes to '%s' by default.", constants.GetGlobalConfigPath(), constants.GetLocalConfigPath()))
	createCmd.Flags().BoolVarP(&overrideExisting, "force", "f", false, "use force. override any existing config files present.")

	return createCmd
}
