package handlers

import (
	"fmt"
	"io"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateCourse(c *fiber.Ctx) error {

	// (title, description, author, thumbnail, preview, duration, is_active, sort_order)
	payloadToClean := database.Course{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Author:      c.FormValue("author"),
		Duration:    c.FormValue("duration"),
	}

	cleanInput, err := inputs.CleanCourseInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isActive := c.FormValue("is_active")
	isActiveBool, err := strconv.ParseBool(isActive)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid boolean value for is_active"})
	}

	thumbnail, err := c.FormFile("thumbnail")
	const MaxFileSize = 10 * 1024 * 1024 // 10MB en bytes
	if thumbnail.Size > MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El archivo es demasiado grande. El tamaño máximo permitido es 10MB.",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	id := uuid.New()
	ext := filepath.Ext(thumbnail.Filename)
	newFilename := fmt.Sprintf("%s%s", id, ext)
	thumbnailsPath := filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "thumbnails")
	c.SaveFile(thumbnail, fmt.Sprintf("%s/%s", thumbnailsPath, newFilename))
	thumbnailToDB := fmt.Sprintf("/web/uploads/thumbnails/%s", newFilename)

	previewTmp := c.FormValue("preview_tmp")
	previewDir := "/web/uploads/previews/" + id.String()
	previewFinalPath := filepath.Join(os.Getenv("ROOT_PATH"), previewDir)
	err = os.MkdirAll(previewFinalPath, 0755)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")
	cmd := exec.Command("sh", ffmpegPath, previewTmp, previewFinalPath)
	err = cmd.Run()
	if err != nil {
		return c.SendStatus(500)
	}

	payloadToDB := database.Course{
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		Author:      cleanInput.Author,
		Thumbnail:   thumbnailToDB,
		Preview:     previewDir + "/master.m3u8",
		Duration:    cleanInput.Duration,
		IsActive:    isActiveBool,
	}

	newCourseID, err := database.CreateCourse(payloadToDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": newCourseID,
	})
}

func ChunkUpload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El archivo preview es requerido.",
		})
	}

	id := c.FormValue("id")
	if id == "" {
		id = uuid.New().String()
	}

	tmpPath := filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "tmp", id)

	chunkNumber, err := strconv.Atoi(c.FormValue("chunkNumber"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al procesar el chunk number.",
		})
	}

	totalChunks, err := strconv.Atoi(c.FormValue("totalChunks"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Error al procesar el total chunks.",
		})
	}

	if err := os.MkdirAll(tmpPath, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al crear el directorio temporal.",
		})
	}

	finalPath := filepath.Join(tmpPath, file.Filename)

	finalFile, err := os.OpenFile(finalPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al abrir el archivo final.",
		})
	}
	defer finalFile.Close()

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al abrir el archivo cargado.",
		})
	}
	defer src.Close()

	_, err = io.Copy(finalFile, src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error al guardar el chunk.",
		})
	}

	if chunkNumber == totalChunks-1 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Archivo cargado con éxito",
			"path":    finalPath,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Chunk cargado con éxito",
		"id":      id,
	})
}
