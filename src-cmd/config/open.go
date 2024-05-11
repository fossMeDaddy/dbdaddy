package configCmd

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	"fmt"
	"os"
	"os/exec"

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
		configFilePath = constants.GetGlobalConfigPath()
		if !lib.FileExists(configFilePath) {
			cmd.PrintErrln("File in the global context doesn't exist. please create one, run 'config -h' for more info")
			return
		}
	} else {
		tmpConfigFilePath, err := lib.GetConfigFilePath()
		if err != nil {
			cmd.PrintErrln("Could not find a config file! please create one. run 'config -h' for more info")
			return
		}

		configFilePath = tmpConfigFilePath
	}
	cmd.Println(fmt.Sprintf("Config file found at: %s, opening via vim...", configFilePath))

	vimOsCmd := exec.Command("vim", configFilePath)
	vimOsCmd.Stdin = os.Stdin
	vimOsCmd.Stdout = os.Stdout
	vimOsCmd.Stderr = os.Stderr

	vimErr := vimOsCmd.Run()
	if vimErr != nil {
		cmd.PrintErrln("Failed to open vim, trying nano...")

		nanoOsCmd := exec.Command("nano", configFilePath)
		nanoOsCmd.Stdin = os.Stdin
		nanoOsCmd.Stdout = os.Stdout
		nanoOsCmd.Stderr = os.Stderr

		nanoErr := nanoOsCmd.Run()
		if nanoErr != nil {
			cmd.PrintErrln("Holy shit bro?! wtf are you using for an OS? no vim, no nano, is this what, an alpine docker container?")
			cmd.PrintErr("nano command gave the error:\n" + nanoErr.Error())
			cmd.PrintErr("vim command gave the error:\n" + vimErr.Error())
			os.Exit(1)
		}
	}
}

func InitOpenCmd() *cobra.Command {
	openCmd.Flags().BoolVarP(&openGlobalConfig, "global", "g", false, "open the global config file, won't even look at the local one")

	return openCmd
}
