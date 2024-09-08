package statusCmd

import (
	"dbdaddy/constants"
	"dbdaddy/middlewares"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "status",
	Short: "Check status of the current database branch",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	cmd.Println("On branch:", viper.GetString(constants.DbConfigCurrentBranchKey))
}

func Init() *cobra.Command {
	// bind flags or something here
	// ...

	return cmd
}
