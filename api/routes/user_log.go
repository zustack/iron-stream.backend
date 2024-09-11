package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserLogRoutes(app *fiber.App) {
	app.Get("/log/user/:userID", middleware.AdminUser, handlers.GetUserLog)
	app.Post("/log/user/logout", middleware.NormalUser, handlers.Logout)
	app.Post("/log/user/found/apps", middleware.NormalUser, handlers.FoundApps)
}
