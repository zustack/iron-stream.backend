package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Statistics struct {
	Date    string `json:"date"`
	Windows int    `json:"windows"`
	Mac     int    `json:"mac"`
	Linux   int    `json:"linux"`
	All     int    `json:"all"`
}

func GetUserStatistics(c *fiber.Ctx) error {
	from := c.Query("from", "")
	to := c.Query("to", "")

	var startDate, endDate time.Time
	var err error

	if from == "" || to == "" {
		// Si no se proporcionan fechas, usar el mes actual
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		endDate = startDate.AddDate(0, 1, -1)
	} else {
		// Parsear las fechas proporcionadas
		startDate, err = time.Parse("2006-01-02", from)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Formato de fecha 'from' inválido. Use YYYY-MM-DD"})
		}
		endDate, err = time.Parse("2006-01-02", to)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Formato de fecha 'to' inválido. Use YYYY-MM-DD"})
		}
	}

	var stats []Statistics
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateForDB := d.Format("02/01/2006")
		dateForJSON := d.Format("2006/01/02")
		linuxCount, _ := database.GetUserCount(dateForDB, "Linux")
		windowsCount, _ := database.GetUserCount(dateForDB, "Windows")
		macCount, _ := database.GetUserCount(dateForDB, "Mac")
		allCount := linuxCount + windowsCount + macCount
		stat := Statistics{
			Date:    dateForJSON,
			Windows: windowsCount,
			Mac:     macCount,
			Linux:   linuxCount,
			All:     allCount,
		}
		stats = append(stats, stat)
	}

	return c.JSON(stats)
}

func GetCurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	return c.Status(200).JSON(fiber.Map{
		"id":         user.ID,
		"email":      user.Email,
		"name":       user.Name,
		"surname":    user.Surname,
		"created_at": user.CreatedAt,
	})
}

func UpdateAdminStatus(c *fiber.Ctx) error {
	userId := c.Params("userId")
	isAdmin := c.Params("isAdmin")
	if userId == "1" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The user with the id 1 cannot be updated",
		})
	}
	err := database.UpdateAdminStatus(userId, isAdmin)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func UpdateSpecialAppUser(c *fiber.Ctx) error {
	userId := c.Params("userId")
	specialApps := c.Params("specialApps")
	err := database.UpdateUserSpecialApps(userId, specialApps)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func UpdateActiveStatusAllUsers(c *fiber.Ctx) error {
	active := c.Params("active")
	err := database.UpdateActiveStatusAllUsers(active)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var content string
	if active == "true" {
		content = "All users were deactivated."
	} else {
		content = "All users were activated."
	}
	l_type := "3"
	err = database.CreateAdminLog(content, l_type)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func UpdateActiveStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "1" {
		return c.Status(400).JSON(fiber.Map{
			"error": "The user with the id 1 cannot be deactivated",
		})
	}

	err := database.UpdateActiveStatus(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func AdminUsers(c *fiber.Ctx) error {
	cursor, err := strconv.Atoi(c.Query("cursor", "0"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	limit := 50
	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"
	isActiveParam := c.Query("a", "")
	isAdminParam := c.Query("admin", "")
	specialAppsParam := c.Query("special", "")
	verifiedParam := c.Query("verified", "")

	from := c.Query("from", "")
	to := c.Query("to", "")

	users, err := database.GetAdminUsers(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam, from, to, limit, cursor)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	searchCount, err := database.GetAdminUsersSearchCount(searchParam, isActiveParam, isAdminParam, specialAppsParam, verifiedParam, from, to)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalCount, err := database.GetAdminUsersCountAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
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
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := database.UpdatePassword(payload.Email, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func DeleteAccountByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	user := c.Locals("user").(*database.User)
	if !user.IsAdmin && user.Email != email {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not allowed to delete this account",
		})
	}
	err := database.DeleteAccountByEmail(email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = database.DeleteNotification(email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func ResendEmailToken(c *fiber.Ctx) error {
	email := c.Params("email")

	code := utils.GenerateCode()
	err := database.UpdateEmailToken(email, code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	subjet := "Verify your email on Iron Stream"
	err = utils.SendEmail(code, email, subjet)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func VerifyEmail(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := database.GetUserByEmail(payload.Email)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if payload.EmailToken != user.EmailToken {
		return c.Status(401).JSON(fiber.Map{
			"error": "The token provided is not valid.",
		})
	}

	err = database.VerifyAccount(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, exp, err := utils.MakeJWT(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token":    token,
		"userId":   user.ID,
		"isAdmin":  user.IsAdmin,
		"exp":      exp,
		"fullName": user.Name + " " + user.Surname,
	})
}

func Signup(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	payloadToClean := inputs.SignupInput{
		Email:    payload.Email,
		Name:     payload.Name,
		Surname:  payload.Surname,
		Password: payload.Password,
		Pc:       payload.Pc,
		Os:       payload.Os,
	}

	cleanInput, err := inputs.Signup(payloadToClean)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	code := utils.GenerateCode()

	payloadToDB := database.User{
		Email:      cleanInput.Email,
		Name:       cleanInput.Name,
		Surname:    cleanInput.Surname,
		Password:   cleanInput.Password,
		Pc:         cleanInput.Pc,
		Os:         cleanInput.Os,
		EmailToken: code,
	}

	err = database.CreateUser(payloadToDB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	subjet := "Verify your email on Iron Stream"
	err = utils.SendEmail(code, cleanInput.Email, subjet)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateNotification("user", cleanInput.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(201)
}

func Login(c *fiber.Ctx) error {
	var payload database.User
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	payloadToClean := inputs.LoginInput{
		Email:    payload.Email,
		Password: payload.Password,
		Pc:       payload.Pc,
	}

	cleanInput, err := inputs.Login(payloadToClean)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := database.GetUserByEmail(cleanInput.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pc), []byte(cleanInput.Pc))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Incorrect unique identifier, please try again.",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cleanInput.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Incorrect password, please try again.",
		})
	}

	token, exp, err := utils.MakeJWT(user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateUserLog("The user has logged in.", "1", user.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"token":    token,
		"userId":   user.ID,
		"isAdmin":  user.IsAdmin,
		"exp":      exp,
		"fullName": user.Name + " " + user.Surname,
	})
}
