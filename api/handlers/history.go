package handlers

import (
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func GetUserHistory(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	history, err := database.GetUserHistory(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(history)
}

func UpdateHistory(c *fiber.Ctx) error {
	type Payload struct {
		Id     string `json:"id"`
		Resume string `json:"resume"`
	}
	var payload Payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}
	err := database.UpdateHistory(payload.Id, payload.Resume)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}
