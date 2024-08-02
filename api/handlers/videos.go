package handlers

import (
	"fmt"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UpdateVideo(c *fiber.Ctx) error {
	payloadToClean := database.Video{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Length:    c.FormValue("length"),
    VideoHLS: c.FormValue("video_tmp"),
	}

	cleanInput, err := inputs.CleanVideoInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  id := c.FormValue("id")
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

  var previewDir string
  if cleanInput.VideoHLS != "" {
    videoId := uuid.New()
    previewDir = "/web/uploads/videos/" + videoId.String()
    previewFinalPath := filepath.Join(os.Getenv("ROOT_PATH"), previewDir)
    err = os.MkdirAll(previewFinalPath, 0755)
    if err != nil {
      return c.SendStatus(fiber.StatusInternalServerError)
    }

    ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")

    cmd := exec.Command("sh", ffmpegPath, cleanInput.VideoHLS, previewFinalPath)
    err = cmd.Run()
    if err != nil {
      return c.SendStatus(500)
    }
  }

	payloadToDB := database.Video {
    Title:       cleanInput.Title,
    Description: cleanInput.Description,
    VideoHLS:    previewDir + "/master.m3u8",
    Thumbnail:   thumbnailToDB,
    Length:      cleanInput.Length,
    ID:          id64,
	}

  err = database.UpdateVideo(payloadToDB)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
  }

  return c.SendStatus(fiber.StatusOK)

}

func DeleteVideo(c *fiber.Ctx) error {
  id := c.Params("id")
  err := database.DeleteVideoByID(id)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.SendStatus(fiber.StatusNoContent)
}

func GetVideos(c *fiber.Ctx) error {
  id := c.Params("id")
  if id == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Course ID is required",
    })
  }

	cursor, err := strconv.Atoi(c.Query("cursor", "0"))
	if err != nil {
		return c.Status(400).SendString("Invalid cursor")
	}
	limit := 50 

	searchParam := c.Query("q", "")
	searchParam = "%" + searchParam + "%"

	videos, err := database.GetVideos(id, searchParam, cursor, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalCount, err := database.GetVideosCount()
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

	response := struct {
		Data       []database.Video `json:"data"`
		PreviousID *int             `json:"previousId"`
		NextID     *int             `json:"nextId"`
	}{
		Data:       videos,
		PreviousID: previousID,
		NextID:     nextID,
	}

	return c.JSON(response)
}

func CreateVideo(c *fiber.Ctx) error {
	payloadToClean := database.Video{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
    VideoHLS:      c.FormValue("video_tmp"),
		Length:    c.FormValue("length"),
	}

	cleanInput, err := inputs.CleanVideoInput(payloadToClean)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

  courseID := c.FormValue("course_id")
  courseID64, err := strconv.ParseInt(courseID, 10, 64)
  if err != nil {
    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
      "status": "fail", "message": "Course ID is invalid",})
  }

  thumbnail, err := c.FormFile("thumbnail")
  if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
	const MaxFileSize = 10 * 1024 * 1024 // 10MB en bytes
	if thumbnail.Size > MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "El archivo es demasiado grande. El tamaño máximo permitido es 10MB.",
		})
	}
	id := uuid.New()
	ext := filepath.Ext(thumbnail.Filename)
	newFilename := fmt.Sprintf("%s%s", id, ext)
	thumbnailsPath := filepath.Join(os.Getenv("ROOT_PATH"), "web", "uploads", "thumbnails")
	c.SaveFile(thumbnail, fmt.Sprintf("%s/%s", thumbnailsPath, newFilename))
	thumbnailToDB := fmt.Sprintf("/web/uploads/thumbnails/%s", newFilename)

	videoId := uuid.New()
	previewDir := "/web/uploads/videos/" + videoId.String()
	previewFinalPath := filepath.Join(os.Getenv("ROOT_PATH"), previewDir)
	err = os.MkdirAll(previewFinalPath, 0755)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")

	cmd := exec.Command("sh", ffmpegPath, cleanInput.VideoHLS, previewFinalPath)
	err = cmd.Run()
	if err != nil {
		return c.SendStatus(500)
	}

	payloadToDB := database.Video {
    Title:       cleanInput.Title,
    Description: cleanInput.Description,
    VideoHLS:    previewDir + "/master.m3u8",
    Thumbnail:   thumbnailToDB,
    Length:      cleanInput.Length,
    CourseID:    courseID64,
	}

  newVideoID, err := database.CreateVideo(payloadToDB)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
  }

  return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": newVideoID})
}
