package handlers

import (
	"fmt"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetForbiddenApps(c *fiber.Ctx) error {
  time.Sleep(4000 * time.Millisecond)
	user := c.Locals("user").(*database.User)
  if user.SpecialApps {
    // get the user_apps del usuario
    userIdStr := fmt.Sprintln(user.ID)
	  userAppIDs, err := database.GetUserAppsIds(userIdStr)
	  if err != nil {
		  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			  "error": err.Error(),
		  })
    }
    // get all the apps with the userAppIds
    apps, err := database.GetAppsByIds(userAppIDs)
    if err != nil {
      return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": err.Error(),
      })
    }
    return c.JSON(apps)
  } 

  apps, err := database.GetActiveApps()
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(apps)
}

func UpdateAppStatus(c *fiber.Ctx) error {
  appId := c.Params("id")
  isActive := c.Params("isActive")
  err := database.UpdateAppStatus(appId, isActive)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(fiber.StatusOK)
}

func GetActiveApps(c *fiber.Ctx) error {
  apps, err := database.GetActiveApps()
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(apps)
}

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

func GetAdminApps(c *fiber.Ctx) error {
  time.Sleep(2000 * time.Millisecond)
	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"

	isActiveParam := c.Query("a", "")

	apps, err := database.GetAdminApps(searchParam, isActiveParam)
  fmt.Println("admin apps", apps)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(apps)
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

func DeleteApp(c *fiber.Ctx) error {
  time.Sleep(2000 * time.Millisecond)
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
