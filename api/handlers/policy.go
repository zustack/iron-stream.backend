package handlers

import (
	"fmt"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func GetPolicys(c *fiber.Ctx) error {
  data, err := database.GetPolicy()
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(data)
}

func GetPolicyItems(c *fiber.Ctx) error {
  policyId := c.Params("policyId")
  data, err := database.GetPolicyItemsByPolicyId(policyId)
  if err != nil {
    fmt.Println(err)
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.JSON(data)
}

func DeletePolicyItem(c *fiber.Ctx) error {
  id := c.Params("id")
  err := database.DeletePolicyItem(id)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(204)
}

func DeletePolicy(c *fiber.Ctx) error {
  id := c.Params("id")
  err := database.DeletePolicy(id)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(204)
}

func CreatePolicyItem(c *fiber.Ctx) error {
  policyId := c.Params("policyId")
	var p database.PolicyItem
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
  if p.Body == "" {
    return c.Status(400).JSON(fiber.Map{
      "error": "The policy item body is required.",
    })
  }
  if len(p.Body) > 255 {
    return c.Status(400).JSON(fiber.Map{
      "error": "The policy item body should not have more than 255 characters.",
    })
  }
  err := database.CreatePolicyItem(p.Body, policyId)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(200)
}

func CreatePolicy(c *fiber.Ctx) error {
	var p database.Policy
	if err := c.BodyParser(&p); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
  if p.Title == "" {
    return c.Status(400).JSON(fiber.Map{
      "error": "Title is required.",
    })
  }
  if len(p.Title) > 55 {
    return c.Status(400).JSON(fiber.Map{
      "error": "The policy title should not have more than 55 characters.",
    })
  }
  err := database.CreatePolicy(p.Title)
  if err != nil {
    return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(200)
}
