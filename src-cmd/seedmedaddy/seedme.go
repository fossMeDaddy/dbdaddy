package seedMeCmd

import "github.com/spf13/cobra"

var cmd = &cobra.Command{
	Use:   "seedmedaddy",
	Short: "Seeds the current selected database branch with sample data and relations",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	// TODO: do this later
	cmd.Println("cumming soon, uwu")
}

func Init() *cobra.Command {
	// setup flags & all here...

	return cmd
}
