package configCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/lib/libUtils"

	"github.com/spf13/cobra"
)

var openGlobalConfig = false

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Opens the config file in a CLI-based text-editor, preference is given to the local config file, if it is not present, then global config file is opened",
	Run:   openCmdRun,
}

func openCmdRun(cmd *cobra.Command, args []string) {
	var configFilePath string
	if openGlobalConfig {
		configFilePath = libUtils.GetGlobalConfigPath()
		if !libUtils.Exists(configFilePath) {
			cmd.PrintErrln("File in the global context doesn't exist. please create one, run 'config -h' for more info")
			return
		}
	} else {
		tmpConfigFilePath, err := libUtils.FindConfigFilePath()
		if err != nil {
			cmd.PrintErrln("Could not find a config file! please create one. run 'config -h' for more info")
			return
		}

		configFilePath = tmpConfigFilePath
	}
	cmd.Println(fmt.Sprintf("Config file found at: %s, opening via vim...", configFilePath))

	libUtils.OpenFileInEditor(configFilePath)
}

func InitOpenCmd() *cobra.Command {
	openCmd.Flags().BoolVarP(&openGlobalConfig, "global", "g", false, "open the global config file, won't even look at the local one")

	return openCmd
}
