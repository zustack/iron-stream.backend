package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func Videos(c *fiber.Ctx) error {
  fullPath := c.Path()
  segments := strings.Split(fullPath, "/")
  if len(segments) > 4 {
		id := segments[4]
		fmt.Println("ID:", id)
    // get the course id and check in the user_courses if the user is enrolled
    // if it is not enrolled, return 401
    // if it is enrolled, continue
	  return c.Next()
	}
  // return 
  return c.SendStatus(401)
}
