package handlers

import (
	"fmt"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func MakeSpecialAppUser(c *fiber.Ctx) error {
  userId := c.Params("userId")
  specialApps := c.Params("specialApps")
  err := database.UpdateUserSpecialApps(userId, specialApps)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(fiber.StatusOK)
}

func DeactivateSpecificCourseForAllUsers(c *fiber.Ctx) error {
	id := c.Params("id")
	// convert id to int64
	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	users, err := database.GetUserIds()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	for _, user := range users {
		err = database.DeactivateCourseForUser(user.ID, id64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeactivateAllCoursesForAllUsers(c *fiber.Ctx) error {
	err := database.DeactivateAllCoursesForAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func UpdateActiveStatusAllUsers(c *fiber.Ctx) error {
	isActive := c.FormValue("isActive")
	// convert isActive to bool
	isActiveBool, err := strconv.ParseBool(isActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = database.UpdateActiveStatusAllUsers(isActiveBool)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func UpdateActiveStatus(c *fiber.Ctx) error {
	id := c.Params("id")
  if id == "1" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "El admin con id 1 no puede ser desactivado.",
		})
  }
	err := database.UpdateActiveStatus(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func AdminUsers(c *fiber.Ctx) error {
	time.Sleep(1000 * time.Millisecond)
	cursor, err := strconv.Atoi(c.Query("cursor", "0"))
	if err != nil {
		return c.Status(400).SendString("Invalid cursor")
	}

	limit := 50
	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"

	isActiveParam := c.Query("a", "")
	isAdminParam := c.Query("admin", "")
	specialAppsParam := c.Query("special", "")
	verifiedParam := c.Query("verified", "")

	users, err := database.GetAdminUsers(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam, limit, cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println(users)

	searchCount, err := database.GetAdminUsersCount(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam)
	if err != nil {
		fmt.Println("el error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalCount, err := database.GetAdminUsersCountAll()
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

	response := struct {
		Data       []database.User `json:"data"`
		TotalCount int             `json:"totalCount"`
		PreviousID *int            `json:"previousId"`
		NextID     *int            `json:"nextId"`
	}{
		Data:       users,
		TotalCount: searchCount,
		PreviousID: previousID,
		NextID:     nextID,
	}

	return c.JSON(response)
}

func UpdatePassword(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	fmt.Println("The new pas", payload.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	fmt.Println("The new hash password", string(hashedPassword))

	err = database.UpdatePassword(string(hashedPassword), user.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no user found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No user found with the provided email",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeleteAccountAtRegister(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	err := database.DeleteAccount(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func RequestEmailTokenResetPassword(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	code := utils.GenerateCode()

	err := database.UpdateEmailToken(payload.Email, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	subjet := "Verifica tu correo electrónico en Iron Stream"
	err = utils.SendEmail(code, payload.Email, subjet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func ResendEmailToken(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	code := utils.GenerateCode()

	subjet := "Verifica tu correo electrónico en Iron Stream"
	err := utils.SendEmail(code, payload.Email, subjet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func VerifyEmail(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	user, err := database.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "No se encontro el usuario con el email ingresado.",
		})
	}

	if payload.EmailToken != user.EmailToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "El codigo es incorrecto",
		})
	}

	err = database.VerifyAccount(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ocurrio un error inesperado y no se pudo verificar la cuenta.",
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

	tokenString, err := tokenByte.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ocurrio un error al generar el token de autenticación.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   tokenString,
		"userId":  user.ID,
		"isAdmin": user.IsAdmin,
		"exp":     now.Add(expDuration).Unix(),
	})
}

func Login(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	if payload.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El correo electrónico es requerido.",
		})
	}

	if payload.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "La contraseña es requerida.",
		})
	}

	if payload.Pc == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ocurrio un error debido a una incompatibilidad con tu sistema operativo.",
		})
	}

	user, err := database.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No se encontro el usuario con el nombre de usuario proporcionado.",
		})
	}

	fmt.Println("payload.Pc", payload.Pc, "user.Pc", user.Pc)
	if user.Pc != payload.Pc {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Esta cuenta esta registrada en otra computadora.",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
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

	tokenString, err := tokenByte.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ocurrio un error al generar el token de autenticación.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   tokenString,
		"userId":  user.ID,
		"isAdmin": user.IsAdmin,
		"exp":     now.Add(expDuration).Unix(),
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
		Password:   cleanInput.Password,
		Email:      cleanInput.Email,
		Name:       cleanInput.Name,
		Surname:    cleanInput.Surname,
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

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}
