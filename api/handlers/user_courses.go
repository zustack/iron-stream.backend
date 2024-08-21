package handlers

import (
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func DeleteUserCoursesByCourseId(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	err := database.DeleteUserCoursesByCourseId(courseId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func DeleteAllUserCourses(c *fiber.Ctx) error {
	err := database.DeleteAllUserCourses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetUserCourses(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	q := c.Query("q", "")
	q = "%" + q + "%"

	userCourseIDs, err := database.GetUserCourseIds(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// the 1 means that only active courses will be returned
	courses, err := database.GetCourses("1", q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	courseIDSet := make(map[int64]bool)
	for _, id := range userCourseIDs {
		courseIDSet[id] = true
	}

	type AllowedCourse struct {
		database.Course
		Allowed bool `json:"allowed"`
	}

	coursesWithAllowed := make([]AllowedCourse, len(courses))
	for i, course := range courses {
		coursesWithAllowed[i] = AllowedCourse{
			Course:  course,
			Allowed: courseIDSet[course.ID],
		}
	}

	response := struct {
		Data []AllowedCourse `json:"data"`
	}{
		Data: coursesWithAllowed,
	}

	return c.JSON(response)
}

// admin
// /create/user/course/:user_id/:course_id
func CreateUserCourse(c *fiber.Ctx) error {
	userID := c.Params("userId")
	courseID := c.Params("courseId")
	err := database.CreateUserCourse(userID, courseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}
