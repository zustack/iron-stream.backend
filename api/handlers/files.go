package handlers

import (
	"fmt"
	"iron-stream/internal/database"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func DeleteFile(c *fiber.Ctx) error {
  id :=  c.Query("id")
  if id == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Video ID is required",
    })
  }

  id64, err := strconv.ParseInt(id, 10, 64)
  if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
      "status": "fail", "message": "video ID is invalid",})
  }

  err = database.DeleteFileByID(id64)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(204)
}

func GetFiles(c *fiber.Ctx) error {
  videoID :=  c.Query("videoID")
  page :=     c.Query("page")
  if videoID == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Video ID is required",
    })
  }
  if page == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Video ID is required",
    })
  }

  page64, err := strconv.ParseInt(page, 10, 64)
  if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
      "status": "fail", "message": "page is invalid",})
  }

  videoID64, err := strconv.ParseInt(videoID, 10, 64)
  if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
      "status": "fail", "message": "video ID is invalid",})
  }

  files, err := database.GetFiles(videoID64, page64)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.Status(fiber.StatusOK).JSON(files)
}

func CreateFile(c *fiber.Ctx) error {
  videoID :=  c.FormValue("videoID")
  page :=     c.FormValue("page")

  if videoID == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Video ID is required",
    })
  }

  if len(videoID) > 155 {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Video ID is too long",
    })
  }

  if page == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Video ID is required",
    })
  }

  if len(page) > 10 {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "only 10 pages allowed",
    })
  }

  page64, err := strconv.ParseInt(page, 10, 64)
  if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
      "status": "fail", "message": "page is invalid",})
  }

  videoID64, err := strconv.ParseInt(videoID, 10, 64)
  if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
      "status": "fail", "message": "video ID is invalid",})
  }

  path, err := c.FormFile("path")
  if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

	const MaxFileSize = 10 * 1024 * 1024 // 10MB en bytes
	if path.Size > MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El archivo es demasiado grande. El tamaño máximo permitido es 10MB.",
		})
	}
	id := uuid.New()
	ext := filepath.Ext(path.Filename)
	newFilename := fmt.Sprintf("%s%s", id, ext)
	pathPath := filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "files")
	c.SaveFile(path, fmt.Sprintf("%s/%s", pathPath, newFilename))
	pathToDB := fmt.Sprintf("/web/uploads/files/%s", newFilename)

	payloadToDB := database.File {
    VideoID: videoID64,
    Page: page64,
    Path: pathToDB,
	}


  newFileID, err := database.CreateFile(payloadToDB)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": newFileID})
}
