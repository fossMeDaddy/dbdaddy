package libServer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartServer() error {
	app := fiber.New()

	app.Use(cors.New())

	defineHandlers(app)

	err := app.Listen(":42069")
	return err
}
