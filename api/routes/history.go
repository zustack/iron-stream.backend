package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func HistoryRoutes(app *fiber.App) {
	app.Get("/history", middleware.NormalUser, handlers.GetUserHistory)
}
