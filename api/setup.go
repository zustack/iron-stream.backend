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
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
	}))
	app.Use("/web/assets/videos", middleware.AllowToWatch)
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
	return app
}
