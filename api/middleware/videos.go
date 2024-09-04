package middleware

import (
	"fmt"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Videos(c *fiber.Ctx) error {
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
		if err.Error() == "No user found with id "+fmt.Sprint(claims["sub"]) {
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

	fullPath := c.Path()
	segments := strings.Split(fullPath, "/")
	if len(segments) > 4 {
		id := segments[4]
		allowed, err := database.UserCourseExists(user.ID, id)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if !allowed {
			return c.Status(401).JSON(fiber.Map{
				"error": "You don't have permission to access this resource.",
			})
		}
		return c.Next()
	}

	return c.SendStatus(401)
}
