package serverHandlers

import (
	"dbdaddy/db/db_int"
	"dbdaddy/types"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func HandleGetTable(c *fiber.Ctx) error {
	tableName := c.Params("name")
	schemaName := c.Params("schema")

	tables, err := db_int.ListTablesInDb()
	if err != nil {
		return err
	}

	dbTable := types.Table{}
	isFound := false
	for _, table := range tables {
		if tableName == table.Name && schemaName == table.Schema {
			isFound = true
			dbTable = table
			break
		}
	}
	if !isFound {
		return c.Status(fiber.StatusNotFound).JSON(types.Response{
			Message: fmt.Sprintf("table '%s.%s' could not be found", schemaName, tableName),
		})
	}

	dbRows, err := db_int.GetRows(fmt.Sprintf("select * from %s.%s", dbTable.Schema, dbTable.Name))
	if err != nil {
		return err
	}

	return c.JSON(types.Response{
		Message: "success",
		Data:    dbRows,
	})
}
