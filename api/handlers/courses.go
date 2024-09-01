package handlers

import (
	"fmt"
	"io"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAdminCourses(c *fiber.Ctx) error {
	q := c.Query("q", "")
	q = "%" + q + "%"

	a := c.Query("a", "")

	courses, err := database.GetCourses(a, q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(courses)
}

// TODO: pass is active from the frontend!
func UpdateCourseActiveStatus(c *fiber.Ctx) error {
	time.Sleep(2000 * time.Millisecond)
	id := c.Params("id")
	err := database.UpdateCourseActiveStatus(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(fiber.StatusOK)
}

type SortCoursesInput struct {
	SortCourses []SortPayload `json:"sort_courses"`
}

type SortPayload struct {
	ID        int64  `json:"id"`
	SortOrder string `json:"sort_order"`
}

// TODO: check que los ids existan
func SortCourse(c *fiber.Ctx) error {
	var payload SortCoursesInput
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}

	for _, item := range payload.SortCourses {
		if item.SortOrder == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Todos lo cursos deben tener un sort order.",
			})
		}
		err := database.EditSortCourses(item.ID, item.SortOrder)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.SendStatus(200)
}

func GetSoloCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	course, err := database.GetCourseById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(course)
}


func DeleteCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	err := database.DeleteCourseByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(204)
}


func UpdateCourse(c *fiber.Ctx) error {
	thumbnail, _ := c.FormFile("thumbnail")
	cleanInput, err := inputs.CleanUpdateCourse(inputs.UpdateCourseInput{
    ID:          c.FormValue("id"),
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Author:      c.FormValue("author"),
		Duration:    c.FormValue("duration"),
		IsActive:    c.FormValue("is_active"),
		Thumbnail:   thumbnail,
    OldThumbnail: c.FormValue("old_thumbnail"),
		Preview:     c.FormValue("preview_tmp"),
    OldPreview:  c.FormValue("old_preview"),
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.UpdateCourse(database.Course{
    ID:          cleanInput.ID,
    Title:       cleanInput.Title,
    Description: cleanInput.Description,
    Author:      cleanInput.Author,
    Thumbnail:   cleanInput.Thumbnail,
    Preview:     cleanInput.Preview,
    Duration:    cleanInput.Duration,
    IsActive:    cleanInput.IsActive,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  return c.SendStatus(200)
}


func CreateCourse(c *fiber.Ctx) error {
	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "The thumbnail is required.",
		})
	}

	cleanInput, err := inputs.CleanCreateCourse(inputs.CreateCourseInput{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Author:      c.FormValue("author"),
		Duration:    c.FormValue("duration"),
		IsActive:    c.FormValue("is_active"),
		Thumbnail:   thumbnail,
		Preview:     c.FormValue("preview_tmp"),
	})

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateCourse(database.Course{
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		Author:      cleanInput.Author,
		Thumbnail:   cleanInput.Thumbnail,
		Preview:     cleanInput.Preview,
		Duration:    cleanInput.Duration,
		IsActive:    cleanInput.IsActive,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  return c.SendStatus(200)
}

func ChunkUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "The file is required",
		})
	}

	chunkNumber, err := strconv.Atoi(c.FormValue("chunkNumber"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid chunk number",
		})
	}

	totalChunks, err := strconv.Atoi(c.FormValue("totalChunks"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid total chunks",
		})
	}

	uuid := c.FormValue("uuid")
	if chunkNumber == 0 {
		err = os.MkdirAll(filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "tmp", uuid), 0755)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Invalid total chunks",
			})
		}
	}

	uploadDir := filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "tmp", uuid)

	tempFilePath := filepath.Join(uploadDir, fmt.Sprintf("%s.part", file.Filename))
	tempFile, err := os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error opening temporary file",
		})
	}
	defer tempFile.Close()

	src, err := file.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error opening uploaded file",
		})
	}
	defer src.Close()

	_, err = io.Copy(tempFile, src)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error saving chunk",
		})
	}

	if chunkNumber == totalChunks-1 {
		finalPath := filepath.Join(uploadDir, file.Filename)
		err = os.Rename(tempFilePath, finalPath)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Error finalizing file",
			})
		}
		return c.SendString(finalPath)
	}

	return c.SendStatus(204)
}
