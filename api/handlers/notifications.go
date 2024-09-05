package handlers

import (
	"iron-stream/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

type NsPayload struct {
	Ns []database.Notification `json:"payload"`
}

func DeleteNotifications(c *fiber.Ctx) error {
	var payload NsPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for _, item := range payload.Ns {
		if item.Info == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Info cannot be empty",
			})
		}
		err := database.DeleteNotification(item.Info)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.SendStatus(204)
}
func GetAdminNotifications(c *fiber.Ctx) error {
	time.Sleep(1 * time.Second)

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

	infoR, err := database.GetInfoNotifications("review")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	infoU, err := database.GetInfoNotifications("user")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := struct {
		UserN   int                     `json:"user_n"`
		ReviewN int                     `json:"review_n"`
		InfoN   []database.Notification `json:"info_r"`
		InfoR   []database.Notification `json:"info_u"`
	}{
		UserN:   userN,
		ReviewN: reviewN,
		InfoN:   infoR,
		InfoR:   infoU,
	}

	return c.JSON(response)
}
