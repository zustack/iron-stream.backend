package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func GetPolicy(c *fiber.Ctx) error {
	data, err := database.GetPolicy()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(data)
}

func DeletePolicy(c *fiber.Ctx) error {
	id := c.Params("id")
	err := database.DeletePolicy(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(204)
}

func CreatePolicy(c *fiber.Ctx) error {
	var p database.Policy
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  cleanInput, err := inputs.CreatePolicy(database.Policy{
    Content: p.Content,
    PType: p.PType,
  })

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
  }

	err = database.CreatePolicy(cleanInput.Content, cleanInput.PType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}
