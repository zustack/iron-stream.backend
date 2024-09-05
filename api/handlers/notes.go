package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func DeleteNote(c *fiber.Ctx) error {
  user := c.Locals("user").(*database.User)
  noteId := c.Params("noteId")
  owner, err := database.GetNoteOwner(noteId)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  if owner != user.ID {
    return c.Status(403).JSON(fiber.Map{
      "error": "You are not the owner of this note.",
    })
  }
  err = database.DeleteNoteById(noteId)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(204)
}

func UpdateNote(c *fiber.Ctx) error {
  user := c.Locals("user").(*database.User)
  noteId := c.Params("noteId")
  owner, err := database.GetNoteOwner(noteId)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  if owner != user.ID {
    return c.Status(403).JSON(fiber.Map{
      "error": "You are not the owner of this note.",
    })
  }
  type payload struct {
    Body string `json:"body"` 
  }

	var p payload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  if p.Body == "" {
    return c.Status(400).JSON(fiber.Map{
      "error": "Body is required.",
    })
  }

	if len(p.Body) > 255 {
    return c.Status(400).JSON(fiber.Map{
      "error": "The body should not have more than 255 characters.",
    })
	}

  err = database.UpdateNote(noteId, p.Body)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  return c.SendStatus(200)
}

func CreateNote(c *fiber.Ctx) error {
  user := c.Locals("user").(*database.User)
  courseId := c.Params("courseId")

  type payload struct {
    Body string `json:"body"` 
    VideoTitle string `json:"video_title"` 
    Time string `json:"time"` 
  }

	var p payload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cleanInput, err := inputs.CreateNote(database.Note{
    Body: p.Body,
    VideoTitle: p.VideoTitle,
    Time: p.Time,
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  err = database.CreateNote(database.Note{
    Body: cleanInput.Body,
    VideoTitle: cleanInput.VideoTitle,
    Time: cleanInput.Time,
    CourseID: courseId,
    UserID: user.ID,
  })

  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  return c.SendStatus(200)
}

func GetNotes(c *fiber.Ctx) error {
  user := c.Locals("user").(*database.User)
  courseId := c.Params("courseId")
  notes, err := database.GetNotes(courseId, user.ID)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(notes)
}
