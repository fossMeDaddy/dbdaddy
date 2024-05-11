package cmd

import (
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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dbdaddy",
	Short: "DBDaddy is a helper tool to use your local databases as if they were managed by Git... in branches.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	lib.InitConfigFile(viper.GetViper())

	// connect & ping sqlDb when initializing
	_, err := db.ConnectDB()
	if err != nil {
		rootCmd.Println("Error occured while connecting to db!\n", err.Error())
		rootCmd.Println("\nTry reviewing your database credentials in the config file, if not present, try creating one. (hint: 'config -h')")
		return
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
