package handlers

import (
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func DeleteUserCoursesByCourseIdAndUserId(c *fiber.Ctx) error {
	userID := c.Params("userId")
	courseID := c.Params("courseId")
	err := database.DeleteUserCoursesByCourseIdAndUserId(userID, courseID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}


func DeleteUserCoursesByCourseId(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	err := database.DeleteUserCoursesByCourseId(courseId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}


func DeleteAllUserCourses(c *fiber.Ctx) error {
	err := database.DeleteAllUserCourses()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}


func GetUserCourses(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	q := c.Query("q", "")
  courses, err := database.GetUserCourses(user.ID, q)
  if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
  }
  return c.JSON(courses)
}

func CreateUserCourse(c *fiber.Ctx) error {
	userID := c.Params("userId")
	courseID := c.Params("courseId")
	err := database.CreateUserCourse(userID, courseID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}
