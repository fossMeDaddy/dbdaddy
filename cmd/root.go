package cmd

import (
	constants "dbdaddy/const"
	"dbdaddy/lib"
	checkoutCmd "dbdaddy/src-cmd/checkout"
	configCmd "dbdaddy/src-cmd/config"
	deleteCmd "dbdaddy/src-cmd/delete"
	"dbdaddy/src-cmd/dumpCmd"
	dumpMeCmd "dbdaddy/src-cmd/dumpmedaddy"
	execCmd "dbdaddy/src-cmd/exec"
	inspectMeCmd "dbdaddy/src-cmd/inspectme"
	listCmd "dbdaddy/src-cmd/list"
	restoreCmd "dbdaddy/src-cmd/restore"
	statusCmd "dbdaddy/src-cmd/status"
	studioCmd "dbdaddy/src-cmd/studio"
	"fmt"
	"os"

	_ "github.com/jackc/pgx"
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

func Execute() {
	if lib.IsFirstTimeUser() {
		fmt.Println("Daddy's home baby.")
		fmt.Println("I'll create a global config for ya, let me know your database url here")

		configFilePath := constants.GetGlobalConfigPath()
		lib.DirExistsCreate(constants.GetGlobalDirPath())
		lib.InitConfigFile(viper.GetViper(), configFilePath, true)

		dbUrlPrompt := promptui.Prompt{
			Label: "Press enter to open the config file in a CLI-based text editor",
		}

		_, err := dbUrlPrompt.Run()
		if err != nil {
			fmt.Println("Cancelling initialization...")
			os.Exit(1)
		}

		lib.OpenFileInEditor(configFilePath)
	} else {
		configFilePath, _ := lib.FindConfigFilePath()
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

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
