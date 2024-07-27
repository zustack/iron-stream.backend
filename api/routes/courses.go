package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CoursesRoutes(app *fiber.App) {
	app.Post("/courses/chunk/upload", middleware.AdminUser, handlers.ChunkUpload)
}
