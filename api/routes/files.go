package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func FilesRoutes(app *fiber.App) {
	app.Delete("/files/:id", middleware.AdminUser, handlers.DeleteFile)
	app.Get("/files/:videoID/:page", middleware.NormalUser, handlers.GetFiles)
	app.Post("/files", middleware.AdminUser, handlers.CreateFile)
}
