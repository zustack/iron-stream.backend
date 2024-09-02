package handlers

import (
	"fmt"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetFeed(c *fiber.Ctx) error {
	time.Sleep(2000 * time.Millisecond)
	user := c.Locals("user").(*database.User)
	courseId := c.Params("courseId")
	searchParam := c.Query("q", "")
	videos, err := database.GetFeed(user.ID, courseId, searchParam)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(videos)
}

type updateHistoryPayload struct {
	Id     string `json:"id"`
	Resume string `json:"resume"`
}

func UpdateHistory(c *fiber.Ctx) error {
	time.Sleep(2000 * time.Millisecond)
	var payload updateHistoryPayload
	if err := c.BodyParser(&payload); err != nil {
		fmt.Println("el error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}
	err := database.UpdateHistory(payload.Id, payload.Resume)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

type createHistoryInput struct {
	Id             string `json:"id"`
	VideoId        string `json:"video_id"`
	CurrentVideoId int64  `json:"current_video_id"`
	CourseId       string `json:"course_id"`
	Resume         string `json:"resume"`
}

func WatchNewVideo(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	var payload createHistoryInput
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No se pudo procesar la solicitud.",
		})
	}
	err := database.UpdateHistory(payload.Id, payload.Resume)
	if err != nil {
		return c.SendStatus(500)
	}

	record, err := database.GetLastVideoByUserIdAndCourseIdAndVideoId(user.ID, payload.CourseId, payload.VideoId)
	if err != nil {
		newRecord, err := database.CreateHistory(user.ID, payload.VideoId, payload.CourseId, "")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(newRecord)
	}

	newRecord, err := database.CreateHistory(user.ID, payload.VideoId, payload.CourseId, record.VideoResume)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	fmt.Println("new record", newRecord)
	return c.JSON(newRecord)
}

func GetCurrentVideo(c *fiber.Ctx) error {
	time.Sleep(2000 * time.Millisecond)
	user := c.Locals("user").(*database.User)
	course_id := c.Params("course_id")
	record, err := database.GetLastVideoByUserIdAndCourseId(user.ID, course_id)
	if err != nil {
		video, err := database.GetFistVideoByCourseId(course_id)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		err = database.UpdateVideoViews(video.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		videoStr := fmt.Sprintf("%d", video.ID)
		newRecord, err := database.CreateHistory(user.ID, videoStr, course_id, "")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fmt.Println("first record", newRecord)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"video":      video,
			"resume":     "",
			"history_id": "",
		})
	}

	video, err := database.GetVideoById(record.VideoId)

	err = database.UpdateVideoViews(video.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"video":      video,
		"resume":     record.VideoResume,
		"history_id": record.ID,
	})
}

func UpdateVideo(c *fiber.Ctx) error {
	payloadToClean := database.Video{
		Title:       c.FormValue("title"),
		Description: c.FormValue("description"),
		Duration:    c.FormValue("duration"),
	}

	cleanInput, err := inputs.CleanUpdateVideoInput(payloadToClean)
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
			"status": "fail", "message": "video ID is invalid"})
	}

	const MaxFileSize = 10 * 1024 * 1024 // 10MB en bytes
	var thumbnailToDB string
	thumbnailToDB = c.FormValue("old_thumbnail")
	fmt.Println("old thumbnail", thumbnailToDB)
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
		err = c.SaveFile(thumbnail, fmt.Sprintf("%s/%s", thumbnailsPath, newFilename))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al guardar el thumbnail"})
		}

		thumbnailToDB = fmt.Sprintf("/web/uploads/thumbnails/%s", newFilename)

	}

	var previewDir string
	previewDir = c.FormValue("old_video")
	fmt.Println("old video", previewDir)
	length := c.FormValue("length")
	if c.FormValue("video_tmp") != "" {

		getLengthCmd := exec.Command("sh", filepath.Join(os.Getenv("ROOT_PATH"), "get-video-length.sh"), c.FormValue("video_tmp"))
		output, err := getLengthCmd.Output()
		if err != nil {
			fmt.Println("err", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		length = strings.Trim(string(output), "\n")

		videoId := uuid.New()
		previewDir = "/web/uploads/videos/" + videoId.String()
		previewFinalPath := filepath.Join(os.Getenv("ROOT_PATH"), previewDir)
		err = os.MkdirAll(previewFinalPath, 0755)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		ffmpegPath := filepath.Join(os.Getenv("ROOT_PATH"), "ffmpeg-convert.sh")

		cmd := exec.Command("sh", ffmpegPath, c.FormValue("video_tmp"), previewFinalPath)
		err = cmd.Run()
		if err != nil {
			return c.SendStatus(500)
		}
		previewDir = previewDir + "/master.m3u8"
	}

	payloadToDB := database.Video{
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		VideoHLS:    previewDir,
		Thumbnail:   thumbnailToDB,
		Length:      length,
		Duration:    cleanInput.Duration,
		ID:          id64,
	}

	err = database.UpdateVideo(payloadToDB)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)

}

func DeleteVideo(c *fiber.Ctx) error {
	time.Sleep(2000 * time.Millisecond)
	id := c.Params("id")
	err := database.DeleteVideoByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(id)
}

func GetAdminVideos(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	q := c.Query("q", "")
	q = "%" + q + "%"

	videos, err := database.GetVideos(courseId, q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(videos)
}

func GetVideos(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	q := c.Query("q", "")
	q = "%" + q + "%"

	videos, err := database.GetVideos(courseId, q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(videos)
}


func CreateVideo(c *fiber.Ctx) error {
	thumbnail, err := c.FormFile("thumbnail")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cleanInput, err := inputs.CreateVideo(inputs.CreateVideoInput{
    Title:       c.FormValue("title"),
    Description: c.FormValue("description"),
    Duration:    c.FormValue("duration"),
    CourseID:    c.FormValue("course_id"),
    Video:       c.FormValue("video_tmp"),
    Thumbnail:   thumbnail,
  })

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.CreateVideo(database.Video{
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		VideoHLS:    cleanInput.VideoHLS,
		Thumbnail:   cleanInput.Thumbnail,
		Duration:    cleanInput.Duration,
		Length:      cleanInput.Length,
		CourseID:    cleanInput.CourseID,
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
      "error": err.Error(),
    })
	}

  return c.SendStatus(200)
}

