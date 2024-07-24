package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	payloadToClean := database.User{
		Username:   payload.Username,
		Password:   payload.Password,
		Email:      payload.Email,
		Name:       payload.Name,
		Surname:    payload.Surname,
		EmailToken: payload.EmailToken,
		Pc:         payload.Pc,
		Os:         payload.Os,
	}

	cleanInput, err := inputs.CleanRegisterInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	code := utils.GenerateCode()

	payloadToDB := database.User{
		Username:   cleanInput.Username,
		Password:   cleanInput.Password,
		Email:      cleanInput.Email,
		Name:       cleanInput.Name,
		Surname:    cleanInput.Surname,
		IsAdmin:    cleanInput.IsAdmin,
		EmailToken: code,
		Pc:         cleanInput.Pc,
		Os:         cleanInput.Os,
	}

	id, err := database.CreateUser(payloadToDB)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "El nombre de usuario esta tomado.",
			})
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "El correo electrónico esta tomado.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	subjet := "Verifica tu correo electrónico en Iron Stream"
	err = utils.SendEmail(code, cleanInput.Email, subjet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if cleanInput.Username == "admin" {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":       id,
			"pc":       cleanInput.Pc,
			"os":       cleanInput.Os,
			"is_admin": cleanInput.IsAdmin,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}
