package middleware

import (
	"fmt"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func NormalUser(c *fiber.Ctx) error {
	tokenString := utils.ExtractTokenFromHeader(c.Get("Authorization"))

	if tokenString == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "You are not logged in.",
		})
	}

	token, err := utils.ParseAndValidateToken(tokenString)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
      "error": err.Error(),
    })
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(401).JSON(fiber.Map{
      "error": "invalid token claim",
    })
	}

	user, err := database.GetUserByID(fmt.Sprint(claims["sub"]))
	if err != nil {
    if err.Error() == "No user found with id " + fmt.Sprint(claims["sub"]) {
		  return c.Status(401).JSON(fiber.Map{
        "error": "No user found with this token",
      })
    }
		return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
	}

	c.Locals("user", &user)
	c.Locals("token", tokenString)

	return c.Next()
}
