package handlers

import (
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func DeleteUserAppsByCourseIdAndUserId(c *fiber.Ctx) error {
	userID := c.Params("userId")
	appId := c.Params("appId")
	err := database.DeleteUserAppsByAppIdAndUserId(userID, appId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(204)
}

func GetEnrolledUserApps(c *fiber.Ctx) error {
	userId := c.Params("userId")
	q := c.Query("q", "")
	q = "%" + q + "%"

	apps, err := database.GetEnrolledUserApps(userId, q)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(apps)
}

func CreateUserApp(c *fiber.Ctx) error {
	userId := c.Params("userId")
	appId := c.Params("appId")
	err := database.CreateUserApp(userId, appId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}
