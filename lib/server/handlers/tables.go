package serverHandlers

import (
	"dbdaddy/db/db_int"
	"dbdaddy/types"

	"github.com/gofiber/fiber/v2"
)

func HandleGetTables(c *fiber.Ctx) error {
	var (
		tables []types.Table
		err    error
	)
	tables, err = db_int.ListTablesInDb()
	if err != nil {
		return err
	}

	return c.JSON(types.Response{
		Message: "success",
		Data:    tables,
	})
}
