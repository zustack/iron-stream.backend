package handlers

import (
	"fmt"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func GetAdminSpecialApps(c *fiber.Ctx) error {
  userId := c.Params("user_id")
  user, err := database.GetUserByID(userId)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  fmt.Println("el os", user.Os)

  userApps, err := database.GetSpecialAppsByUserId(user.Os, user.ID)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  apps, err := database.GetAppsByOs(user.Os)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  // Crear un mapa para facilitar la búsqueda de userApps
  userAppMap := make(map[string]bool)
  for _, userApp := range userApps {
    userAppMap[userApp.Name] = true
  }

  // Modificar las apps para agregar el campo "is"
  modifiedApps := make([]map[string]interface{}, len(apps))
  for i, app := range apps {
    appMap := make(map[string]interface{})

    // Copiar todos los campos de la app al nuevo mapa
    appMap["id"] = app.ID
    appMap["name"] = app.Name
    appMap["process_name"] = app.ProcessName
    appMap["os"] = app.Os
    appMap["active"] = app.IsActive
    appMap["created_at"] = app.CreatedAt

    // Agregar el campo "is" si la app está en userApps
    if userAppMap[app.Name] {
      appMap["is"] = true
    }

    modifiedApps[i] = appMap
  }

  return c.JSON(modifiedApps)
}

func GetUserApps(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)

  if user.SpecialApps {
    apps, err := database.GetSpecialAppsByUserId(user.Os, user.ID)
    if err != nil {
      return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": err.Error(),
      })
    }
    return c.JSON(apps)
  }

  apps, err := database.GetAppsByOsAndIsActive(user.Os)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(apps)
}

type AppsPayload struct {
	Apps   []database.App `json:"apps"`
	UserId int64          `json:"user_id"`
}

func CreateSpecialApp(c *fiber.Ctx) error {
	var payload AppsPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

  err := database.UpdateUserSpecialApps(payload.UserId, true)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  err = database.DeleteAllSpecialAppsByUserId(payload.UserId)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  if len(payload.Apps) == 0 {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "No se encontraron Apps.",
    })
  }

	for _, item := range payload.Apps {

		payloadToDB := database.SpecialApp{
			UserId:      payload.UserId,
			Name:        item.Name,
			ProcessName: item.ProcessName,
			Os:          item.Os,
			IsActive:    item.IsActive,
		}

		_, err := database.CreateSpecialApp(payloadToDB)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

	}

	return c.SendStatus(200)
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

type AppResponse struct {
	Data       []database.App `json:"data"`
	TotalCount int            `json:"totalCount"`
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
