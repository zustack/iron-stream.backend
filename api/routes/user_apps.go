package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserAppsRoutes(app *fiber.App) {
  app.Post("user/apps/create/:userId/:appId", middleware.AdminUser, handlers.CreateUserApp)
  app.Get("user/apps/user/apps/:userId", middleware.NormalUser, handlers.GetUserApps)
  app.Delete("user/apps/delete/user/apps/:userId/:appId", middleware.NormalUser, handlers.DeleteUserAppsByCourseIdAndUserId)
}
