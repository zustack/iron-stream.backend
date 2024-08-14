package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func GetAppByID(c *fiber.Ctx) error {
	id := c.Params("id")
	app, err := database.GetAppByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(app)
}

type AppResponse struct {
	Data       []database.App `json:"data"`
	TotalCount int           `json:"totalCount"`
}

func GetApps(c *fiber.Ctx) error {
	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"

	isActiveParam := c.Query("a", "")

	apps, err := database.GetApps(searchParam, isActiveParam)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalCount, err := database.GetAppsCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	response := AppResponse{
		Data:       apps,
		TotalCount: totalCount,
	}

	return c.JSON(response)
}

func UpdateApp(c *fiber.Ctx) error {
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
	}

	cleanInput, err := inputs.CleanAppInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	payloadToDB := database.App{
		ID:          payload.ID,
		Name:        cleanInput.Name,
		ProcessName: cleanInput.ProcessName,
		Os:          cleanInput.Os,
		IsActive:    payload.IsActive,
	}

	err = database.UpdateApp(payloadToDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func GetAppsByOsAndIsActive(c *fiber.Ctx) error {
	os := c.Params("os")
	apps, err := database.GetAppsByOsAndIsActive(os)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(apps)
}

func DeleteApp(c *fiber.Ctx) error {
	id := c.Params("id")

	err := database.DeleteAppByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

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
