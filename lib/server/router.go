package libServer

import (
	serverHandlers "dbdaddy/lib/server/handlers"

	"github.com/gofiber/fiber/v2"
)

func defineHandlers(app *fiber.App) {
	app.Get("/branches", serverHandlers.HandleGetBranches)

	app.Get("/current-branch", serverHandlers.HandleGetCurrentBranch)
	app.Put("/current-branch", serverHandlers.HandlePutCurrentBranch)

	app.Get("/tables", serverHandlers.HandleGetTables)

	app.Get("/table/:schema/:name/rows", serverHandlers.HandleGetTableRows)
	app.Get("/table/:schema/:name/schema", serverHandlers.HandleGetTableSchema)
}
