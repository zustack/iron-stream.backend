package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AppsRoutes(app *fiber.App) {
	app.Post("/apps/create", middleware.AdminUser, handlers.CreateApp)
}
