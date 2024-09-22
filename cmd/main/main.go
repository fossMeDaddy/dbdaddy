package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fossmedaddy/dbdaddy/cli/checkoutCmd"
	"github.com/fossmedaddy/dbdaddy/cli/cloneCmd"
	"github.com/fossmedaddy/dbdaddy/cli/configCmd"
	"github.com/fossmedaddy/dbdaddy/cli/deleteCmd"
	"github.com/fossmedaddy/dbdaddy/cli/dumpCmd"
	"github.com/fossmedaddy/dbdaddy/cli/dumpMeCmd"
	"github.com/fossmedaddy/dbdaddy/cli/execCmd"
	"github.com/fossmedaddy/dbdaddy/cli/initCmd"
	"github.com/fossmedaddy/dbdaddy/cli/inspectMeCmd"
	"github.com/fossmedaddy/dbdaddy/cli/listCmd"
	"github.com/fossmedaddy/dbdaddy/cli/migrationsCmd"
	"github.com/fossmedaddy/dbdaddy/cli/remoteCmd"
	"github.com/fossmedaddy/dbdaddy/cli/restoreCmd"
	"github.com/fossmedaddy/dbdaddy/cli/statusCmd"
	"github.com/fossmedaddy/dbdaddy/cli/studioCmd"
	"github.com/fossmedaddy/dbdaddy/cli/versionCmd"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/manifoldco/promptui"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:               "dbdaddy",
	Short:             "DBDaddy is a helper tool to use your local databases as if they were managed by Git... in branches.",
	Run:               rootCmdRun,
	PersistentPreRunE: rootPreRun,
}

func rootPreRun(cmd *cobra.Command, args []string) error {
	fullCmdPath := strings.Split(cmd.CommandPath(), " ")
	if len(fullCmdPath) <= 1 {
		return nil
	}
	cmdPath := strings.Join(fullCmdPath[1:], " ")

	// DANGEROUS, FUCKED WHEN COMMAND CHANGE
	if strings.HasPrefix(cmdPath, "help") {
		return nil
	}
	if strings.HasSuffix(cmdPath, "--help") {
		return nil
	}
	if strings.HasPrefix(cmdPath, "config") {
		return nil
	}
	if strings.HasPrefix(cmdPath, "clone") {
		return nil
	}
	if strings.HasPrefix(cmdPath, "init") {
		return nil
	}

	if lib.IsFirstTimeUser() {
		cmd.Println(fmt.Sprintf("Daddy's home baby. (version: %s)", globals.Version))
		cmd.Println("I'll create a global config for ya, let me know your database url here")

		configDirPath := libUtils.GetGlobalDirPath()
		libUtils.EnsureDirExists(configDirPath)
		lib.InitConfigFile(viper.GetViper(), configDirPath, true)

		dbUrlPrompt := promptui.Prompt{
			Label: "Press enter to open the config file in a CLI-based text editor",
		}

		_, err := dbUrlPrompt.Run()
		if err != nil {
			cmd.Println("Cancelling initialization...")
			os.Exit(1)
		}

		configFilePath := libUtils.GetGlobalConfigPath()
		libUtils.OpenFileInEditor(configFilePath)
	}

	if err := lib.EnsureSupportedDbDriver(); err != nil {
		return err
	}

	return nil
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func main() {

	if !lib.IsFirstTimeUser() {
		configFilePath, _ := libUtils.FindConfigFilePath()
		lib.ReadConfig(viper.GetViper(), configFilePath)
	}

	rootCmd.AddCommand(versionCmd.Init())
	rootCmd.AddCommand(checkoutCmd.Init())
	rootCmd.AddCommand(statusCmd.Init())
	rootCmd.AddCommand(deleteCmd.Init())
	rootCmd.AddCommand(configCmd.Init())
	rootCmd.AddCommand(dumpMeCmd.Init())
	rootCmd.AddCommand(dumpCmd.Init())
	rootCmd.AddCommand(listCmd.Init())
	rootCmd.AddCommand(restoreCmd.Init())
	rootCmd.AddCommand(inspectMeCmd.Init())
	rootCmd.AddCommand(studioCmd.Init())
	rootCmd.AddCommand(execCmd.Init())
	rootCmd.AddCommand(migrationsCmd.Init())
	rootCmd.AddCommand(initCmd.Init())
	rootCmd.AddCommand(remoteCmd.Init())
	rootCmd.AddCommand(cloneCmd.Init())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
