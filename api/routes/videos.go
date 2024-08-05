package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func VideosRoutes(app *fiber.App) {
  app.Put("/history/watch", middleware.NormalUser, handlers.WatchNewVideo)
  app.Get("/history/current/:course_id", middleware.NormalUser, handlers.GetCurrentVideo)

	app.Put("/videos", middleware.AdminUser, handlers.UpdateVideo)
  app.Delete("/videos/:id", middleware.AdminUser, handlers.DeleteVideo)
	app.Post("/videos", middleware.AdminUser, handlers.CreateVideo)
  app.Get("/videos/:id", middleware.NormalUser, handlers.GetVideos)
}
