package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserCoursesRoutes(app *fiber.App) {
	app.Delete("/user/courses/solo/:userId/:courseId", middleware.AdminUser, handlers.DeleteUserCoursesByCourseIdAndUserId)
	app.Delete("/user/courses/solo/:courseId", middleware.NormalUser, handlers.DeleteUserCoursesByCourseId)
	app.Delete("/user/courses/all", middleware.NormalUser, handlers.DeleteAllUserCourses)
  app.Get("/user/courses/:userId", middleware.NormalUser, handlers.GetUserCourses)
	app.Post("/user/courses/:userId/:courseId", middleware.NormalUser, handlers.CreateUserCourse)
}
