package schemaCmd

import (
	"github.com/fossmedaddy/dbdaddy/cli/schemaCmd/schemaPushCmd"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/fossmedaddy/dbdaddy/middlewares"
	"github.com/spf13/cobra"
)

var cmdRunFn = middlewares.Apply(run, middlewares.CheckConnection)

var cmd = &cobra.Command{
	Use:   "schema",
	Short: "manage schema in a project environment",
	Run:   cmdRunFn,
}

func run(cmd *cobra.Command, args []string) {
	_, cwdIsProject, err := libUtils.CwdIsProject()
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	if !cwdIsProject {
		cmd.PrintErrln("CWD is not a project! please run 'init [postgresql://...|mysql://|sqlite://]' to initialize a new project")
		return
	}

	cmd.Help()
}

func Init() *cobra.Command {
	cmd.AddCommand(schemaPushCmd.Init())

	return cmd
}
