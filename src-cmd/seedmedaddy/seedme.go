package seedMeCmd

import (
	"dbdaddy/middlewares"

	"github.com/spf13/cobra"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "seedmedaddy",
	Short: "Seeds the current selected database branch with sample data and relations",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	// TODO: do this later
	cmd.Println("cumming soon, uwu")
}

func Init() *cobra.Command {
	// setup flags & all here...

	return cmd
}
