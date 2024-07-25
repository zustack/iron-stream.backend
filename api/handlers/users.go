package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	payloadToClean := database.User{
		Username: payload.Username,
		Password: payload.Password,
		Pc:       payload.Pc,
	}

	cleanInput, err := inputs.CleanLoginInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := database.GetUserByUsername(cleanInput.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No se encontro el usuario con el nombre de usuario proporcionado.",
		})
	}

	if user.Pc != cleanInput.Pc {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Esta cuenta esta registrada en otra computadora.",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cleanInput.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "La contraseña es incorrecta.",
		})
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	expDuration := time.Hour * 24 * 30
	claims["sub"] = user.ID
	claims["exp"] = now.Add(expDuration).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(utils.GetEnv("SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ocurrio un error al generar el token de autenticación.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tokenString,
	})
}

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
