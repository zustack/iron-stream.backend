package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")
	err := database.DeleteFileByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(204)
}

func GetFiles(c *fiber.Ctx) error {
	videoID := c.Params("videoID")
	page := c.Params("page")

	files, err := database.GetFiles(videoID, page)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	getTotalPages, err := database.GetTotalPagesByVideoId(videoID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"files":     files,
		"pageCount": getTotalPages,
	})
}

func CreateFile(c *fiber.Ctx) error {
	file, err := c.FormFile("path")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "The file is required.",
		})
	}

	cleanInput, err := inputs.CreateFile(inputs.FileInput{
		VideoID: c.FormValue("videoID"),
		Page:    c.FormValue("page"),
		File:    file,
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateFile(database.File{
		VideoID: cleanInput.VideoID,
		Page:    cleanInput.Page,
		Path:    cleanInput.Path,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}
