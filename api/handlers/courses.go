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
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ACUInput struct {
	CourseID int64 `json:"course_id"`
	UserID   int64 `json:"user_id"`
}

func AddCourseToUser(c *fiber.Ctx) error {
	var payload ACUInput
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}
	err := database.AddCourseToUser(payload.UserID, payload.CourseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(200)
}

func GetAdminCourses(c *fiber.Ctx) error {
	cursor, err := strconv.Atoi(c.Query("cursor", "0"))
	if err != nil {
		return c.Status(400).SendString("Invalid cursor")
	}

	limit := 50
	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"

	isActiveParam := c.Query("a", "")

	courses, err := database.GetAdminCourses(searchParam, isActiveParam, limit, cursor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalCount, err := database.GetAppsCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var previousID, nextID *int
	if cursor > 0 {
		prev := cursor - limit
		if prev < 0 {
			prev = 0
		}
		previousID = &prev
	}
	if cursor+limit < totalCount {
		next := cursor + limit
		nextID = &next
	}

	response := struct {
		Data       []database.Course `json:"data"`
		PreviousID *int              `json:"previousId"`
		NextID     *int              `json:"nextId"`
	}{
		Data:       courses,
		PreviousID: previousID,
		NextID:     nextID,
	}

	return c.JSON(response)
}

func GetCourses(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)

	userCourseIDs := make(map[int64]bool)

	if user.Courses != "" {
		courseIDStrings := strings.Split(user.Courses, ",")
		for _, idStr := range courseIDStrings {
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course ID format"})
			}
			userCourseIDs[id] = true
		}
	}

	cursor, err := strconv.Atoi(c.Query("cursor", "0"))
	if err != nil {
		return c.Status(400).SendString("Invalid cursor")
	}
	limit := 50 

	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"

	courses, err := database.GetCourses(searchParam, cursor, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalCount, err := database.GetCoursesCount()
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var previousID, nextID *int
	if cursor > 0 {
		prev := cursor - limit
		if prev < 0 {
			prev = 0
		}
		previousID = &prev
	}
	if cursor+limit < totalCount {
		next := cursor + limit
		nextID = &next
	}

	type AllowedCourses struct {
		database.Course
		Allowed bool `json:"allowed"`
	}

	coursesWithOn := make([]AllowedCourses, len(courses))
	for i, course := range courses {
		coursesWithOn[i] = AllowedCourses{
			Course:  course,
			Allowed: userCourseIDs[course.ID],
		}
	}

	response := struct {
		Data       []AllowedCourses `json:"data"`
		PreviousID *int             `json:"previousId"`
		NextID     *int             `json:"nextId"`
	}{
		Data:       coursesWithOn,
		PreviousID: previousID,
		NextID:     nextID,
	}

	return c.JSON(response)
}

func DeleteCourse(c *fiber.Ctx) error {
	id := c.Params("id")

	err := database.DeleteCourseByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)

}

func UpdateCourse(c *fiber.Ctx) error {
	payloadToClean := database.Course{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Author:      c.FormValue("author"),
		Duration:    c.FormValue("duration"),
	}

	id := c.FormValue("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid course ID"})
	}

	// convert id to int64
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

	previewTmpDir := c.FormValue("previewTmpDir")
	var previewFinalPath string
	if previewTmpDir != "" {
		previewId := uuid.New()
		previewDir := "/web/uploads/previews/" + previewId.String()
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
	}

	err = database.UpdateCourse(database.Course{
		ID:          id64,
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		Author:      cleanInput.Author,
		Thumbnail:   thumbnailToDB,
		Preview:     previewFinalPath + "/master.m3u8",
		Duration:    cleanInput.Duration,
		IsActive:    isActiveBool,
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
