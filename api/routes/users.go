package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	app.Put("deactivate/all/courses", middleware.AdminUser, handlers.DeactivateAllCoursesForAllUsers)
	app.Put("deactivate/course/for/all/users/:id", middleware.AdminUser, handlers.DeactivateSpecificCourseForAllUsers)

	app.Put("/update/active/status", middleware.AdminUser, handlers.UpdateActiveStatusAllUsers)
	app.Put("/update/active/status/:id", middleware.AdminUser, handlers.UpdateActiveStatus)

	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Post("/verify", handlers.VerifyEmail)
	app.Post("/resend/email/token", handlers.ResendEmailToken)
	app.Post("/delete/account/at/register", handlers.DeleteAccountAtRegister)
	app.Put("/update/password", middleware.NormalUser, handlers.UpdatePassword)
	app.Put("request/email/token/reset/password", handlers.RequestEmailTokenResetPassword)
	app.Get("/users/admin", middleware.AdminUser, handlers.AdminUsers)
}
