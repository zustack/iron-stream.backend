package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func PolicyRoutes(app *fiber.App) {
  app.Delete("/policy/:id", middleware.AdminUser, handlers.DeletePolicy)
  app.Delete("/policy/:id", middleware.AdminUser, handlers.DeletePolicyItem)
  app.Get("/policy/:policyId", middleware.AdminUser, handlers.GetPolicyItems)
  app.Get("/policy", middleware.AdminUser, handlers.GetPolicys)
  app.Post("/policy/create/:policyId", middleware.AdminUser, handlers.CreatePolicy)
	app.Post("/policy/create", middleware.AdminUser, handlers.CreatePolicy)
}
