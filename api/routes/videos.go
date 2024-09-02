package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func VideosRoutes(app *fiber.App) {
	app.Put("/history/watch", middleware.NormalUser, handlers.WatchNewVideo)
	app.Get("/history/current/:course_id", middleware.NormalUser, handlers.GetCurrentVideo)
	app.Put("/history/update", middleware.NormalUser, handlers.UpdateHistory)
	app.Put("/videos", middleware.AdminUser, handlers.UpdateVideo)
	app.Delete("/videos/:id", middleware.AdminUser, handlers.DeleteVideo)
	app.Get("/videos/:courseId", middleware.NormalUser, handlers.GetVideos)

  app.Get("/videos/feed/:courseId", middleware.NormalUser, handlers.GetFeed)
	app.Get("/videos/admin/:courseId", middleware.AdminUser, handlers.GetAdminVideos)
  app.Post("/videos", middleware.AdminUser, handlers.CreateVideo)
}
