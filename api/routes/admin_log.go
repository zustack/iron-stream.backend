package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminLogRoutes(app *fiber.App) {
	app.Get("/log/admin", middleware.AdminUser, handlers.GetAdminLog)
	app.Post("/log/admin", middleware.AdminUser, handlers.CreateAdminLog)
}
