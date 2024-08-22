package handlers

import (
	"fmt"
	"iron-stream/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func DeleteUserAppsByCourseIdAndUserId(c *fiber.Ctx) error {
	userID := c.Params("userId")
	appId := c.Params("appId")
	err := database.DeleteUserAppsByAppIdAndUserId(userID, appId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetUserApps(c *fiber.Ctx) error {
  time.Sleep(2000 * time.Millisecond)
  userId := c.Params("userId")
	q := c.Query("q", "")
	q = "%" + q + "%"

	// Obtener los IDs de las apps del usuario
	userAppIDs, err := database.GetUserAppsIds(userId)
  fmt.Println("user apps ids", userAppIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Obtener todas las apps activas que coincidan con la búsqueda
	apps, err := database.GetAdminApps(q, "")
  fmt.Println("admin apps", apps)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Crear un set de IDs de apps del usuario para fácil consulta
	appIDSet := make(map[int64]bool)
	for _, id := range userAppIDs {
		appIDSet[id] = true
	}

	// Crear una estructura que contenga la app junto con el campo "exists"
	type ExistingApp struct {
		database.App
		Exists bool `json:"exists"`
	}

	appsWithExistence := make([]ExistingApp, len(apps))
	for i, app := range apps {
		appsWithExistence[i] = ExistingApp{
			App:    app,
			Exists: appIDSet[app.ID], // Si el ID está en el set, existe
		}
	}

	// Preparar la respuesta
	response := struct {
		Data []ExistingApp `json:"data"`
	}{
		Data: appsWithExistence,
	}

	// Enviar la respuesta JSON
	return c.JSON(response)
}

func CreateUserApp(c *fiber.Ctx) error {
	userId := c.Params("userId")
	appId := c.Params("appId")
	err := database.CreateUserApp(userId, appId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
