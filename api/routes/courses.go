package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func CoursesRoutes(app *fiber.App) {
	app.Put("/courses/add/user", middleware.AdminUser, handlers.AddCourseToUser)
	app.Get("/courses", middleware.NormalUser, handlers.GetCourses)
	app.Get("/courses/admin", middleware.AdminUser, handlers.GetAdminCourses)
	app.Post("/courses/chunk/upload", middleware.AdminUser, handlers.ChunkUpload)
	app.Post("/courses/create", middleware.AdminUser, handlers.CreateCourse)
	app.Put("/courses/update", middleware.AdminUser, handlers.UpdateCourse)
	app.Delete("/courses/delete/:id", middleware.AdminUser, handlers.DeleteCourse)
}
