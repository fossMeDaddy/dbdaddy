package main

import (
	"fmt"
	"os"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/lib"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	checkoutCmd "github.com/fossmedaddy/dbdaddy/src-cmd/checkout"
	configCmd "github.com/fossmedaddy/dbdaddy/src-cmd/config"
	deleteCmd "github.com/fossmedaddy/dbdaddy/src-cmd/delete"
	"github.com/fossmedaddy/dbdaddy/src-cmd/dumpCmd"
	dumpMeCmd "github.com/fossmedaddy/dbdaddy/src-cmd/dumpmedaddy"
	execCmd "github.com/fossmedaddy/dbdaddy/src-cmd/exec"
	inspectMeCmd "github.com/fossmedaddy/dbdaddy/src-cmd/inspectme"
	listCmd "github.com/fossmedaddy/dbdaddy/src-cmd/list"
	migrationsCmd "github.com/fossmedaddy/dbdaddy/src-cmd/migrations"
	restoreCmd "github.com/fossmedaddy/dbdaddy/src-cmd/restore"
	statusCmd "github.com/fossmedaddy/dbdaddy/src-cmd/status"
	studioCmd "github.com/fossmedaddy/dbdaddy/src-cmd/studio"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dbdaddy",
	Short: "DBDaddy is a helper tool to use your local databases as if they were managed by Git... in branches.",
	Run:   rootCmdRun,
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func main() {
	if lib.IsFirstTimeUser() {
		fmt.Println("Daddy's home baby.")
		fmt.Println("I'll create a global config for ya, let me know your database url here")

		configFilePath := constants.GetGlobalConfigPath()
		libUtils.EnsureDirExists(constants.GetGlobalDirPath())
		lib.InitConfigFile(viper.GetViper(), configFilePath, true)

		dbUrlPrompt := promptui.Prompt{
			Label: "Press enter to open the config file in a CLI-based text editor",
		}

		_, err := dbUrlPrompt.Run()
		if err != nil {
			fmt.Println("Cancelling initialization...")
			os.Exit(1)
		}

		libUtils.OpenFileInEditor(configFilePath)
	} else {
		configFilePath, _ := libUtils.FindConfigFilePath()
		lib.ReadConfig(viper.GetViper(), configFilePath)
		lib.EnsureSupportedDbDriver()
	}

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

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
