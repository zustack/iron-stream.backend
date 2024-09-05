package routes

import (
	"iron-stream/api/handlers"
	"iron-stream/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func NotesRoutes(app *fiber.App) {
  app.Put("notes/:noteId", middleware.NormalUser, handlers.UpdateApp)
  app.Delete("notes/:noteId", middleware.NormalUser, handlers.DeleteNote)
	app.Get("notes", middleware.NormalUser, handlers.GetNotes)
	app.Post("notes", middleware.NormalUser, handlers.CreateNote)
}
