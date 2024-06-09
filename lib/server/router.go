package libServer

import (
	serverHandlers "dbdaddy/lib/server/handlers"

	"github.com/gofiber/fiber/v2"
)

func defineHandlers(app *fiber.App) {
	app.Get("/branches", serverHandlers.HandleGetBranches)

	app.Get("/current-branch", serverHandlers.HandleGetCurrentBranch)
	app.Put("/current-branch", serverHandlers.HandlePutBranch)

	app.Get("/tables", serverHandlers.HandleGetTables)

	app.Get("/table/:schema/:name", serverHandlers.HandleGetTable)
}
