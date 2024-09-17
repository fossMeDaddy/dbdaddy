package remoteLsCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:   "ls",
	Short: "list all available remote origins",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	origins := viper.GetStringMapString(constants.DbConfigOriginsKey)

	if len(origins) == 0 {
		cmd.Println("no available origins.")
		return
	}

	cmd.Println("listing available origins")

	i := 1
	for originKey := range origins {
		cmd.Println(fmt.Sprintf("%d. %s", i, originKey))
		i++
	}
}

func Init() *cobra.Command {
	// flags here

	return cmd
}
