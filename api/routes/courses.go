package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CoursesRoutes(app *fiber.App) {
	app.Get("/courses/admin", middleware.AdminUser, handlers.GetAdminCourses)
	app.Put("/courses/update/active/:id", middleware.AdminUser, handlers.UpdateCourseActiveStatus)
	app.Post("/courses/sort", middleware.AdminUser, handlers.SortCourse)
	app.Put("/courses/update", middleware.AdminUser, handlers.UpdateCourse)

  app.Get("/courses/solo/:id", middleware.NormalUser, handlers.GetSoloCourse)
  app.Delete("/courses/delete/:id", middleware.AdminUser, handlers.DeleteCourse)
	app.Post("/courses/create", middleware.AdminUser, handlers.CreateCourse)
	app.Post("/courses/chunk/upload", middleware.AdminUser, handlers.ChunkUpload)
}
