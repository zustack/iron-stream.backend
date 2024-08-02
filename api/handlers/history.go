package handlers

import (
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)


func GetUserHistory(c *fiber.Ctx) error {
  user := c.Locals("user").(*database.User)
  history, err := database.GetHistoryByUserID(user.ID)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(history)
}
