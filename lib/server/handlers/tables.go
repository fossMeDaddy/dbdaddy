package serverHandlers

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	"dbdaddy/types"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func HandleGetTables(c *fiber.Ctx) error {
	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	var (
		tables []types.Table
		err    error
	)
	err = lib.SwitchDB(viper.GetViper(), currBranch, func() error {
		tables, err = db_int.ListTablesInDb()
		return err
	})
	if err != nil {
		return err
	}

	return c.JSON(types.Response{
		Message: "success",
		Data:    tables,
	})
}
