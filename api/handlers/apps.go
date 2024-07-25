package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"strconv"

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
	PreviousID *int           `json:"previousId"`
	NextID     *int           `json:"nextId"`
}

func GetApps(c *fiber.Ctx) error {
	cursor, err := strconv.Atoi(c.Query("cursor", "0"))
	if err != nil {
		return c.Status(400).SendString("Invalid cursor")
	}

	limit := 50
	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"

	isActiveParam := c.Query("a", "")
	isActiveParam = "%" + isActiveParam + "%"
	apps, err := database.GetApps(searchParam, isActiveParam, limit, cursor)
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

	var previousID, nextID *int
	if cursor > 0 {
		prev := cursor - limit
		if prev < 0 {
			prev = 0
		}
		previousID = &prev
	}
	if cursor+limit < totalCount {
		next := cursor + limit
		nextID = &next
	}

	response := AppResponse{
		Data:       apps,
		PreviousID: previousID,
		NextID:     nextID,
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
