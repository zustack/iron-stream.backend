package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func FilesRoutes(app *fiber.App) {
	app.Get("/files", middleware.NormalUser, handlers.GetFiles)
	app.Delete("/files", middleware.AdminUser, handlers.DeleteFile)
	app.Post("/files", middleware.AdminUser, handlers.CreateFile)
}
