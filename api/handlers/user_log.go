package handlers

import (
	"fmt"
	"iron-stream/internal/database"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	err := database.CreateUserLog("The user has logged out.", "1", user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func FoundApps(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)

	type app struct {
		Name string `json:"name"`
	}

	type payload struct {
		Apps       []app  `json:"apps"`
		VideoTitle string `json:"video_title"`
	}

	var p payload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var text string
	appNames := []string{}
	for _, a := range p.Apps {
		appNames = append(appNames, a.Name)
	}

	if len(appNames) == 1 {
		text = fmt.Sprintf("The app %s was open while watching the video %s.", appNames[0], p.VideoTitle)
	} else if len(appNames) > 1 {
		text = fmt.Sprintf("The apps %s where open while watching the video %s.", strings.Join(appNames, ", "), p.VideoTitle)
	} else {
		return c.Status(500).JSON(fiber.Map{
			"error": "It needs at least one app.",
		})
	}

	err := database.CreateUserLog(text, "3", user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func GetUserLog(c *fiber.Ctx) error {
	userID := c.Params("userID")
	uls, err := database.GetUserLog(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(uls)
}
