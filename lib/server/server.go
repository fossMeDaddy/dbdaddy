package libServer

import "github.com/gofiber/fiber/v2"

func StartServer() error {
	app := fiber.New()

	defineHandlers(app)

	err := app.Listen(":42069")
	return err
}
