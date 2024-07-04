package serverHandlers

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/types"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func HandleGetTableRows(c *fiber.Ctx) error {
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
		return fmt.Errorf("table '%s.%s' could not be found", schemaName, tableName)
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

func HandleGetTableSchema(c *fiber.Ctx) error {
	tableName := c.Params("name")
	schemaName := c.Params("schema")

	currBranch := viper.GetString(constants.DbConfigCurrentBranchKey)

	schema, err := db_int.GetTableSchema(currBranch, schemaName, tableName)
	if err != nil {
		return err
	}

	return c.JSON(types.Response{
		Message: "success",
		Data:    schema,
	})
}

func HandlePostExec(c *fiber.Ctx) error {
	type ReqSchema struct {
		Query string
	}
	var reqBody ReqSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	data, err := db_int.GetRows(reqBody.Query)
	if err != nil {
		return err
	}

	return c.JSON(types.Response{
		Message: "success",
		Data:    data,
	})
}
