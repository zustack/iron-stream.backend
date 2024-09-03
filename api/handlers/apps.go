package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func UpdateAppEa(c *fiber.Ctx) error {
	appId := c.Params("id")
	ea := c.Params("ea")
	err := database.UpdateAppEa(appId, ea)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func UpdateAppStatus(c *fiber.Ctx) error {
	appId := c.Params("id")
	isActive := c.Params("isActive")
	err := database.UpdateAppStatus(appId, isActive)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func GetForbiddenApps(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)

	if user.SpecialApps {
		apps, err := database.GetUserApps(user.ID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(apps)
	}

	apps, err := database.GetActiveApps()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(apps)
}

func GetAdminApps(c *fiber.Ctx) error {
	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"
	isActiveParam := c.Query("a", "")

	apps, err := database.GetAdminApps(searchParam, isActiveParam)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(apps)
}

func UpdateApp(c *fiber.Ctx) error {
	var payload database.App
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cleanInput, err := inputs.UpdateApp(database.App{
		ID:          payload.ID,
		Name:        payload.Name,
		ProcessName: payload.ProcessName,
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.UpdateApp(database.App{
		ID:            cleanInput.ID,
		Name:          cleanInput.Name,
		ProcessName:   cleanInput.ProcessName,
		IsActive:      payload.IsActive,
		ExecuteAlways: payload.ExecuteAlways,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func DeleteApp(c *fiber.Ctx) error {
	id := c.Params("id")
	err := database.DeleteAppByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(204)
}

func CreateApp(c *fiber.Ctx) error {
	var payload database.App
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cleanInput, err := inputs.CreateApp(database.App{
		Name:        payload.Name,
		ProcessName: payload.ProcessName,
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateApp(database.App{
		Name:          cleanInput.Name,
		ProcessName:   cleanInput.ProcessName,
		IsActive:      payload.IsActive,
		ExecuteAlways: payload.ExecuteAlways,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}
