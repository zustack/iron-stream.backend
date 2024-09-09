package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func PolicyRoutes(app *fiber.App) {
	app.Delete("/policy/:id", middleware.AdminUser, handlers.DeletePolicy)
	app.Get("/policy", handlers.GetPolicy)
	app.Post("/policy", middleware.AdminUser, handlers.CreatePolicy)
}
