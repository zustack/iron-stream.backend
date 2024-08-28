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
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	time.Sleep(2000 * time.Millisecond)
	id := c.Params("id")

	err := database.DeleteCourseByID(id)
	if err != nil {
		fmt.Println("el error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendString(id)
}

func UpdateCourse(c *fiber.Ctx) error {
	payloadToClean := database.Course{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Author:      c.FormValue("author"),
		Duration:    c.FormValue("duration"),
	}

	sortOrder := c.FormValue("sortOrder")
	sortOrderInt, err := strconv.ParseInt(sortOrder, 10, 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course ID"})
	}

	id := c.FormValue("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course ID"})
	}

	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course ID"})
	}

	cleanInput, err := inputs.CleanUpdateCourseInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isActive := c.FormValue("is_active")
	isActiveBool, err := strconv.ParseBool(isActive)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid boolean value for is_active",
		})
	}

	const MaxFileSize = 10 * 1024 * 1024 // 10MB en bytes
	var thumbnailToDB string
	thumbnailToDB = c.FormValue("old_thumbnail")
	thumbnail, err := c.FormFile("thumbnail")
	if err == nil {
		if thumbnail.Size > MaxFileSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "El archivo es demasiado grande. El tamaño máximo permitido es 10MB.",
			})
		}
		thumbnail_id := uuid.New()
		ext := filepath.Ext(thumbnail.Filename)
		newFilename := fmt.Sprintf("%s%s", thumbnail_id, ext)
		thumbnailsPath := filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "thumbnails")
		// save the file in local disk
		err = c.SaveFile(thumbnail, fmt.Sprintf("%s/%s", thumbnailsPath, newFilename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al guardar el thumbnail"})
		}
		thumbnailToDB = fmt.Sprintf("/web/uploads/thumbnails/%s", newFilename)
	}

	previewTmpDir := c.FormValue("preview_tmp")
	var previewDir string
	previewDir = c.FormValue("old_video")
	if previewTmpDir != "" {
		previewId := uuid.New()
		previewDir = "/web/uploads/previews/" + previewId.String()
		previewFinalPath := filepath.Join(os.Getenv("ROOT_PATH"), previewDir)
		err = os.MkdirAll(previewFinalPath, 0755)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")
		cmd := exec.Command("sh", ffmpegPath, previewTmpDir, previewFinalPath)
		err = cmd.Run()
		if err != nil {
			return c.SendStatus(500)
		}
		previewDir = previewDir + "/master.m3u8"
	}

	fmt.Println("isVideo", c.FormValue("isVideo"))
	if c.FormValue("isVideo") == "false" {
		previewDir = ""
	}

	err = database.UpdateCourse(database.Course{
		ID:          id64,
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		Author:      cleanInput.Author,
		Thumbnail:   thumbnailToDB,
		Preview:     previewDir,
		Duration:    cleanInput.Duration,
		IsActive:    isActiveBool,
		SortOrder:   int(sortOrderInt),
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func CreateCourse(c *fiber.Ctx) error {
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

	if previewTmp != "" {
		ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")
		cmd := exec.Command("sh", ffmpegPath, previewTmp, previewFinalPath)
		err = cmd.Run()
		if err != nil {
			fmt.Println("the error22", err)
			return c.SendStatus(500)
		}
	}

	if previewTmp == "" {
		previewDir = ""
	} else {
		previewDir = previewDir + "/master.m3u8"
	}

	payloadToDB := database.Course{
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		Author:      cleanInput.Author,
		Thumbnail:   thumbnailToDB,
		Preview:     previewDir,
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
		return c.Status(400).SendString("Error retrieving the file")
	}

	chunkNumber, err := strconv.Atoi(c.FormValue("chunkNumber"))
	if err != nil {
		return c.Status(400).SendString("Invalid chunk number")
	}

	totalChunks, err := strconv.Atoi(c.FormValue("totalChunks"))
	if err != nil {
		return c.Status(400).SendString("Invalid total chunks")
	}

	uuid := c.FormValue("uuid")
	if chunkNumber == 0 {
		err = os.MkdirAll(filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "tmp", uuid), 0755)
		if err != nil {
			fmt.Println("the error", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	uploadDir := filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "tmp", uuid)

	tempFilePath := filepath.Join(uploadDir, fmt.Sprintf("%s.part", file.Filename))
	tempFile, err := os.OpenFile(tempFilePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return c.Status(500).SendString("Error opening temporary file")
	}
	defer tempFile.Close()

	src, err := file.Open()
	if err != nil {
		return c.Status(500).SendString("Error opening uploaded file")
	}
	defer src.Close()

	_, err = io.Copy(tempFile, src)
	if err != nil {
		return c.Status(500).SendString("Error saving chunk")
	}

	if chunkNumber == totalChunks-1 {
		finalPath := filepath.Join(uploadDir, file.Filename)
		err = os.Rename(tempFilePath, finalPath)
		if err != nil {
			return c.Status(500).SendString("Error finalizing file")
		}
		return c.SendString(finalPath)
	}

	return c.SendStatus(204)
}
