package handlers

import (
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func DeleteReview(c *fiber.Ctx) error {
  id := c.Params("id")
  err := database.DeleteReview(id)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(204)
}

func UpdatePublicStatus(c *fiber.Ctx) error {
  public := c.Params("public")
  id := c.Params("id")
  err := database.UpdatePublicStatus(public, id)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(200)
}

func GetAdminReviews(c *fiber.Ctx) error {
	searchParam := c.Query("q", "")
	publicParam := c.Query("p", "")
  reviews, err := database.GetAdminReviews(searchParam, publicParam)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(reviews)
}

func GetPublicReviewsByCourseId(c *fiber.Ctx) error {
  courseId := c.Params("courseId")
  reviews, err := database.GetPublicReviewsByCourseId(courseId)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(reviews)
}

func CreateReview(c *fiber.Ctx) error {
  user := c.Locals("user").(*database.User)

  type payload struct {
    CourseId       string `json:"course_id"`
    Description    string `json:"description"`
    Rating         float64 `json:"rating"`
  }

	var p payload
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  cleanInput, err := inputs.CreateReview(database.Review{
    Description: p.Description,
    Rating: p.Rating,
  })

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  _, err = database.GetCourseById(p.CourseId)
  if err != nil {
    if err.Error() == "No course found with the id " + p.CourseId {
      return c.Status(404).JSON(fiber.Map{
        "error": err.Error(),
      })
    }
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  exists, err := database.UserReviewExists(user.ID, p.CourseId)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  if exists {
    return c.Status(400).JSON(fiber.Map{
      "error": "You have already left a review for this course.",
    })
  }

	allowed, err := database.UserCourseExists(user.ID, p.CourseId)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !allowed {
		return c.Status(401).JSON(fiber.Map{
			"error": "You don't have permission to create reviews for this course.",
		})
	}

  author := user.Name + " " + user.Surname
  err = database.CreateReview(user.ID, p.CourseId, author, cleanInput.Description, cleanInput.Rating)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(200)
}

