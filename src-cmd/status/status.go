package statusCmd

import (
	constants "dbdaddy/const"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of the current database branch",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Println("On branch:", viper.GetString(constants.DbConfigCurrentBranchKey))
}

func Init() *cobra.Command {
	// bind flags or something here
	// ...

	return cmd
}
