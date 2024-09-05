package handlers

import (
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

/*
func CreateNotification(c *fiber.Ctx) error {
  nType := c.Params("type")
  err := database.CreateNotification(nType)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(200)
}
*/

func GetAdminNotifications(c *fiber.Ctx) error {
  userN, err := database.GetNotifications("user")
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  reviewN, err := database.GetNotifications("review")
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

	response := struct {
		UserN int             `json:"user_n"`
    ReviewN int            `json:"review_n"`
	}{
    UserN: userN,
    ReviewN: reviewN,
	}

  return c.JSON(response)
}
