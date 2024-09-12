package handlers

import (
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func CreateAdminLog(c *fiber.Ctx) error {
	var payload database.AdminLog
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err := database.CreateAdminLog(payload.Content, payload.LType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(204)
}

func GetAdminLog(c *fiber.Ctx) error {
	data, err := database.GetAdminLog()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(data)
}
