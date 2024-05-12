package cmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/lib"
	checkoutCmd "dbdaddy/src-cmd/checkout"
	configCmd "dbdaddy/src-cmd/config"
	deleteCmd "dbdaddy/src-cmd/delete"
	listCmd "dbdaddy/src-cmd/list"
	seedMeCmd "dbdaddy/src-cmd/seedmedaddy"
	statusCmd "dbdaddy/src-cmd/status"
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
	if lib.IsFirstTimeUser() {
		cmd.Println("Daddy's home baby.")
		cmd.Println("I'll create a global config for ya, let me know your database url here")

		configFilePath := constants.GetGlobalConfigPath()
		lib.EnsureDirExists(constants.GetGlobalDirPath())
		lib.InitConfigFile(viper.GetViper(), configFilePath, true)

		dbUrlPrompt := promptui.Prompt{
			Label: "Press enter to open the config file in a CLI-based text editor",
		}

		_, err := dbUrlPrompt.Run()
		if err != nil {
			cmd.Println("Cancelling initialization...")
			os.Exit(1)
		}
		lib.OpenConfigFileAt(configFilePath)

		return
	}

	cmd.Help()
}

func Execute() {
	if !lib.IsFirstTimeUser() {
		configFilePath, _ := lib.FindConfigFilePath()
		lib.InitConfigFile(viper.GetViper(), configFilePath, false)

		_, err := db.ConnectDB()
		if err != nil {
			rootCmd.Println("Error occured while connecting to db!\n", err.Error())
			rootCmd.Println("\nTry reviewing your database credentials in the config file, if not present, try creating one. (hint: 'config -h')")
			return
		}
	}

	rootCmd.AddCommand(checkoutCmd.Init())
	rootCmd.AddCommand(statusCmd.Init())
	rootCmd.AddCommand(deleteCmd.Init())
	rootCmd.AddCommand(configCmd.Init())
	rootCmd.AddCommand(seedMeCmd.Init())
	rootCmd.AddCommand(listCmd.Init())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
