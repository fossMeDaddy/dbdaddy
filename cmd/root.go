package cmd

import (
	constants "dbdaddy/const"
	"dbdaddy/db"
	"dbdaddy/lib"
	checkoutCmd "dbdaddy/src-cmd/checkout"
	deleteCmd "dbdaddy/src-cmd/delete"
	statusCmd "dbdaddy/src-cmd/status"
	"fmt"
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
	viper.SetConfigName("dbdaddy.config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	// FIRST detect .env file & prompt user with auto-detected file & correct the detection
	// THEN, prompt him with the auto-detected url & correct detection
	// populate the config file with correct driver, user, passwords & db names
	// check for supported drivers & display unsupported error messages

	// in case of key update in an updated cli

	// setup default values for PG connection
	viper.SetDefault(constants.DbConfigHostKey, "localhost")
	viper.SetDefault(constants.DbConfigPortKey, 5432)
	viper.SetDefault(constants.DbConfigUserKey, "postgres")
	viper.SetDefault(constants.DbConfigPassKey, "postgres")
	viper.SetDefault(constants.DbConfigDriverKey, constants.DbDriverPostgres)
	viper.SetDefault(constants.DbConfigDbNameKey, "postgres")
	viper.SetDefault(constants.DbConfigCurrentBranchKey, "postgres")

	configNotFoundErr := viper.ReadInConfig()
	if configNotFoundErr != nil {
		// TODO: if not found, try searching at ~/.dbdaddy/dbdaddy.config.json

		// TODO: this thing, later
		// driverSelect := promptui.Select{
		// 	Label: "Choose your database driver",
		// 	Items: []string{constants.DbDriverPostgres, constants.DbDriverMySQL, constants.DbDriverSqlite},
		// }
		// _, selectedDriver, _ := driverSelect.Run()
		// if selectedDriver == "" {
		// 	viper.Set(constants.DbConfigDbNameKey, constants.DbDriverPostgres)
		// } else {
		// 	viper.Set(constants.DbConfigDbNameKey, selectedDriver)
		// }

		// if selectedDriver != constants.DbDriverPostgres {
		// 	rootCmd.Println(fmt.Sprintf("%s driver is not yet supported!", selectedDriver))
		// }

		lib.FindCwdDBCreds()

		writeErr := viper.SafeWriteConfig()
		if writeErr != nil {
			fmt.Println("Terribly sorry, this is probably a skill issue!\nI can't generate your dbdaddy.config.json")
			fmt.Println(writeErr)
			return
		}
	}

	// connect & ping sqlDb when initializing
	_, err := db.ConnectDB()
	if err != nil {
		fmt.Println("Error occured while connecting to db!\n", err)
		return
	}

	rootCmd.AddCommand(checkoutCmd.Init())
	rootCmd.AddCommand(statusCmd.Init())
	rootCmd.AddCommand(deleteCmd.Init())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
