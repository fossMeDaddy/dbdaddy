package serverHandlers

import (
	constants "dbdaddy/const"
	"dbdaddy/db/db_int"
	"dbdaddy/lib"
	"dbdaddy/types"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func HandleGetBranches(c *fiber.Ctx) error {
	dbs, err := db_int.GetExistingDbs()
	if err != nil {
		return err
	}
	dbs = append(dbs, constants.SelfDbName)

	return c.JSON(types.Response{
		Message: "success",
		Data:    dbs,
	})
}

func HandleGetCurrentBranch(c *fiber.Ctx) error {
	return c.JSON(types.Response{
		Message: "success",
		Data:    viper.GetString(constants.DbConfigCurrentBranchKey),
	})
}

func HandlePutBranch(c *fiber.Ctx) error {
	type reqSchema struct {
		BranchName string
	}

	var reqBody reqSchema
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	if err := lib.SetCurrentBranch(reqBody.BranchName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(types.Response{
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(types.Response{
		Message: "success",
		Data:    nil,
	})
}
