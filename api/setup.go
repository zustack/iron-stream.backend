package api

import (
	"iron-stream/api/middleware"
	"iron-stream/api/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Setup() *fiber.App {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	}))
	app.Use("/web/uploads/videos", middleware.Videos)
	staticPath := os.Getenv("ROOT_PATH") + "/web/uploads"
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		log.Fatalf("Static path does not exist: %s", staticPath)
	}
	app.Static("/web/uploads", staticPath)
	routes.UserRoutes(app)
	routes.AppsRoutes(app)
	routes.CoursesRoutes(app)
	routes.VideosRoutes(app)
	routes.FilesRoutes(app)
	routes.HistoryRoutes(app)
	routes.UserCoursesRoutes(app)
	routes.UserAppsRoutes(app)
	routes.ReviewRoutes(app)
	routes.NotesRoutes(app)
	routes.NotificationsRoutes(app)
	routes.PolicyRoutes(app)
	routes.UserLogRoutes(app)
	routes.AdminLogRoutes(app)
	return app
}
