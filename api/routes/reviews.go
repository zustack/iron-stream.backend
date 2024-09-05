package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func ReviewRoutes(app *fiber.App) {
	app.Delete("/reviews/:id", middleware.AdminUser, handlers.DeleteReview)
	app.Put("/reviews/update/public/:public/:id", middleware.AdminUser, handlers.UpdatePublicStatus)
	app.Get("/reviews/admin", middleware.AdminUser, handlers.GetAdminReviews)
	app.Get("/reviews/public/:courseId", middleware.AdminUser, handlers.GetPublicReviewsByCourseId)
	app.Post("/reviews", middleware.NormalUser, handlers.CreateReview)
}
