package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AppsRoutes(app *fiber.App) {
	app.Get("/apps", handlers.GetApps)
	app.Put("/apps/update", middleware.AdminUser, handlers.UpdateApp)
	app.Get("/apps/get/:id", middleware.AdminUser, handlers.GetAppByID)
	app.Get("/apps/normal-user/:os", middleware.NormalUser, handlers.GetAppsByOsAndIsActive)
	app.Delete("/apps/delete/:id", middleware.AdminUser, handlers.DeleteApp)
	app.Post("/apps/create", middleware.AdminUser, handlers.CreateApp)
	app.Post("/special/apps/create", middleware.AdminUser, handlers.CreateSpecialApp)

	app.Get("special/apps/get/:user_id", middleware.AdminUser, handlers.GetAdminSpecialApps)

	app.Get("/special/apps/get", middleware.AdminUser, handlers.GetUserApps)
}
