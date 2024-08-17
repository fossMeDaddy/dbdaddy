package configCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	"dbdaddy/libUtils"
	"fmt"
	"path"

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
	}

	dir, _ := path.Split(configWritePath)
	if _, err := libUtils.EnsureDirExists(dir); err != nil {
		panic("Unexpected error occured!\n" + err.Error())
	}

	if overrideExisting {
		v := viper.New()
		lib.InitConfigFile(v, configWritePath, false)

		if err := v.WriteConfigAs(configWritePath); err != nil {
			cmd.PrintErrln("Error occured while writing config file!\n" + err.Error())
			return
		}
	} else {
		if err := viper.SafeWriteConfigAs(configWritePath); err != nil {
			cmd.PrintErrln(err.Error())
			return
		}
	}

	libUtils.OpenFileInEditor(configWritePath)
}

func InitCreateCmd() *cobra.Command {
	createCmd.Flags().BoolVarP(&createInGlobalNamespace, "global", "g", false, fmt.Sprintf("write config file at '%s' for global use, writes to '%s' by default.", constants.GetGlobalConfigPath(), constants.GetLocalConfigPath()))
	createCmd.Flags().BoolVarP(&overrideExisting, "force", "f", false, "use force. override any existing config files present.")

	return createCmd
}
