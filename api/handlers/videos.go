package handlers

import (
	"fmt"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"
	"time"

	"github.com/gofiber/fiber/v2"
)


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

	video, err := database.GetVideoById("1")

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
  thumbnail, _ := c.FormFile("thumbnail")
	cleanInput, err := inputs.UpdateVideo(inputs.UpdateVideoInput{
    ID:          c.FormValue("id"),
    Title:       c.FormValue("title"),
    Description: c.FormValue("description"),
    Duration:    c.FormValue("duration"),
    Thumbnail:   thumbnail,
    OldThumbnail: c.FormValue("old_thumbnail"),
    Video:       c.FormValue("video_tmp"),
    OldVideoHLS: c.FormValue("old_video_hls"),
  })

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	payloadToDB := database.Video{
		Title:       cleanInput.Title,
		Description: cleanInput.Description,
		VideoHLS:    cleanInput.VideoHLS,
		Thumbnail:   cleanInput.Thumbnail,
		Length:      cleanInput.Length,
		Duration:    cleanInput.Duration,
		ID:          cleanInput.ID,
	}

	err = database.UpdateVideo(payloadToDB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
      "error": err.Error()},
    )
	}

	return c.SendStatus(200)
}


func DeleteVideo(c *fiber.Ctx) error {
	id := c.Params("id")
	err := database.DeleteVideoByID(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(204)
}


func GetFeed(c *fiber.Ctx) error {
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



func GetAdminVideos(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	q := c.Query("q", "")
	q = "%" + q + "%"

	videos, err := database.GetAdminVideos(courseId, q)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
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

