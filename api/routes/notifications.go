package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func NotificationsRoutes(app *fiber.App) {
	app.Get("/notifications", middleware.AdminUser, handlers.GetAdminNotifications)
}
