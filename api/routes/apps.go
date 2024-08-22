package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func AppsRoutes(app *fiber.App) {
  app.Put("apps/update/status/:id/:isActive", middleware.AdminUser, handlers.UpdateAppStatus)
	app.Get("apps/active", middleware.NormalUser, handlers.GetActiveApps)
  app.Get("apps/by/id/:id", middleware.AdminUser, handlers.GetAppByID)
  app.Get("apps/admin", middleware.AdminUser, handlers.GetAdminApps)
  app.Put("apps/update", middleware.AdminUser, handlers.UpdateApp)
  app.Delete("apps/delete/:id", middleware.AdminUser, handlers.DeleteApp)
  app.Post("apps/create", middleware.AdminUser, handlers.CreateApp)
}
