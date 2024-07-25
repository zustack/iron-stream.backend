package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func CreateApp(c *fiber.Ctx) error {
	var payload database.App
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	payloadToClean := database.App{
		Name:        payload.Name,
		ProcessName: payload.ProcessName,
		Os:          payload.Os,
		IsActive:    payload.IsActive,
	}

	cleanInput, err := inputs.CleanAppInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	payloadToDB := database.App{
		Name:        cleanInput.Name,
		ProcessName: cleanInput.ProcessName,
		Os:          cleanInput.Os,
		IsActive:    cleanInput.IsActive,
	}

	id, err := database.CreateApp(payloadToDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}
