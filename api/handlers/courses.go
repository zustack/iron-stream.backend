package handlers

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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
