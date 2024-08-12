package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Post("/verify", handlers.VerifyEmail)
	app.Post("/resend/email/token", handlers.ResendEmailToken)
	app.Post("/delete/account/at/register", handlers.DeleteAccountAtRegister)
	app.Put("/update/password", middleware.NormalUser, handlers.UpdatePassword)
	app.Put("request/email/token/reset/password", handlers.RequestEmailTokenResetPassword)
}
