package cloneCmd

import (
	"fmt"
	"path"

	"github.com/fossmedaddy/dbdaddy/db"
	"github.com/fossmedaddy/dbdaddy/db/db_int"
	"github.com/fossmedaddy/dbdaddy/globals"
	"github.com/fossmedaddy/dbdaddy/lib/libUtils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	forceFlag      = false
	onlySchemaFlag = false
)

var cmd = &cobra.Command{
	Use:     "clone",
	Short:   "provide a uri to take a dump of the remote database into the currently connected instance",
	Example: "clone postgresql://user:pwd@localhost:5432/dbname",
	Run:     run,
	Args:    argFn,
}

func argFn(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("invalid number of arguments provided, need max: 2, min: 1 args")
	}

	return nil
}

func run(cmd *cobra.Command, args []string) {
	remoteUri := args[0]

	remoteConnConfig, err := libUtils.GetConnConfigFromUri(remoteUri)
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	dbname := ""
	if len(args) == 2 {
		dbname = args[1]
	} else {
		dbname = remoteConnConfig.Database
	}

	cmd.Println("started dump process...")

	configFilePath, _ := libUtils.FindConfigFilePath()
	dumpOutFile := path.Join(
		libUtils.GetDriverDumpDir(configFilePath, remoteConnConfig.Driver),
		libUtils.GetDumpFileName(dbname),
	)
	if err := db_int.DumpDb(dumpOutFile, remoteConnConfig, onlySchemaFlag); err != nil {
		cmd.PrintErrln("error occured while taking a dump of the remote database!")
		cmd.PrintErrln(err)
		return
	}

	cmd.Println(fmt.Sprintf("remote database dump complete, saved at: %s", dumpOutFile))

	if _, err := db.ConnectSelfDb(viper.GetViper()); err != nil {
		cmd.PrintErrln(err)
		cmd.PrintErrln("could not connect to a database to restore remote database into...")
		return
	}

	if db_int.DbExists(dbname) && !forceFlag {
		cmd.PrintErrln(fmt.Printf("database with name %s already exists, choose a different db name or use --force flag.", dbname))
		return
	}

	connConfig := globals.CurrentConnConfig
	connConfig.Database = dbname
	if err := db_int.RestoreDb(connConfig, dumpOutFile, true); err != nil {
		cmd.PrintErrln(err)
		return
	}

	originConfigKey := libUtils.GetDbConfigOriginKey(dbname)
	viper.Set(originConfigKey, remoteConnConfig)
	if err := viper.WriteConfig(); err != nil {
		cmd.PrintErrln("error occured while writing to config file")
		cmd.PrintErrln(err)
		return
	}

	cmd.Println(fmt.Sprintf("remote database '%s' cloned successfully.", dbname))
}

func Init() *cobra.Command {
	// flags here
	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "force remote database copying even if there already exists a database with same name")
	cmd.Flags().BoolVar(&onlySchemaFlag, "only-schema", false, "dump only schema, exclude data from dump")

	return cmd
}
