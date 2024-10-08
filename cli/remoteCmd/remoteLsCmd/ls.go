package remoteLsCmd

import (
	"fmt"

	"github.com/fossmedaddy/dbdaddy/constants"
	"github.com/fossmedaddy/dbdaddy/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmd = &cobra.Command{
	Use:   "ls",
	Short: "list all available remote origins",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	origins := types.DbConfigOrigins{}
	if err := viper.UnmarshalKey(constants.DbConfigOriginsKey, &origins); err != nil {
		cmd.PrintErrln("error occured while parsing origins")
		cmd.PrintErrln(err)
		return
	}

	if len(origins) == 0 {
		cmd.Println("no available origins.")
		return
	}

	cmd.Println("listing available origins for databases")

	i := 1
	for originKey, originConnConfig := range origins {
		cmd.Println(fmt.Sprintf("%d. %s <-> %s:%s/%s", i, originKey, originConnConfig.Host, originConnConfig.Port, originKey))
		i++
	}
}

func Init() *cobra.Command {
	// flags here

	return cmd
}
