package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func VideosRoutes(app *fiber.App) {
	app.Put("/videos/s_review/:id/:s_review", middleware.AdminUser, handlers.UpdateVideoSReview)
	app.Put("/history/watch", middleware.NormalUser, handlers.WatchNewVideo)
	app.Get("/history/current/:courseId", middleware.NormalUser, handlers.GetCurrentVideo)
	app.Put("/videos", middleware.AdminUser, handlers.UpdateVideo)
	app.Delete("/videos/:id", middleware.AdminUser, handlers.DeleteVideo)
	app.Get("/videos/feed/:courseId", middleware.NormalUser, handlers.GetFeed)
	app.Get("/videos/admin/:courseId", middleware.AdminUser, handlers.GetAdminVideos)
	app.Post("/videos", middleware.AdminUser, handlers.CreateVideo)
}
