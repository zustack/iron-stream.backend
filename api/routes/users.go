package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	app.Get("/users/stats", middleware.AdminUser, handlers.GetUserStatistics)
	app.Get("/users/current", middleware.NormalUser, handlers.GetCurrentUser)
	app.Put("/users/update/admin/status/:userId/:isAdmin", middleware.AdminUser, handlers.UpdateAdminStatus)
	app.Put("/users/update/special/apps/user/:userId/:specialApps", middleware.AdminUser, handlers.UpdateSpecialAppUser)
	app.Put("/users/update/all/active/status/:active", middleware.AdminUser, handlers.UpdateActiveStatusAllUsers)
	app.Put("/users/update/active/status/:id", middleware.AdminUser, handlers.UpdateActiveStatus)
	app.Get("/users/admin", middleware.AdminUser, handlers.AdminUsers)
	app.Put("/users/update/password", middleware.NormalUser, handlers.UpdatePassword)
	app.Delete("/users/delete/account/by/email/:email", middleware.NormalUser, handlers.DeleteAccountByEmail)
	app.Post("/users/resend/email/token/:email", handlers.ResendEmailToken)
	app.Post("/users/verify", handlers.VerifyEmail)
	app.Post("/users/signup", handlers.Signup)
	app.Post("/users/login", handlers.Login)
}
