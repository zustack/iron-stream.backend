package handlers

import (
	"fmt"
	"iron-stream/api/inputs"
	"iron-stream/internal/database"

	"github.com/gofiber/fiber/v2"
)

func UpdateVideoSReview(c *fiber.Ctx) error {
	id := c.Params("id")
	sReview := c.Params("s_review")
	err := database.UpdateVideoSReview(sReview, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendStatus(200)
}

func WatchNewVideo(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	type createHistoryInput struct {
		HistoryId      string `json:"history_id"`
		VideoId        int64  `json:"video_id"`
		CurrentVideoId int64  `json:"current_video_id"`
		CourseId       string `json:"course_id"`
		Resume         string `json:"resume"`
	}
	var payload createHistoryInput
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := database.UpdateHistory(payload.HistoryId, payload.Resume)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	record, err := database.GetLastVideoByUserIdAndCourseIdAndVideoId(user.ID, payload.CourseId, payload.VideoId)
	if err != nil {
		if err.Error() == "Record not found" {
			newRecord, err := database.CreateHistory(user.ID, payload.VideoId, payload.CourseId, "")
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return c.JSON(newRecord)
		} else {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	newRecord, err := database.CreateHistory(user.ID, payload.VideoId, payload.CourseId, record.VideoResume)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(newRecord)

}

func GetCurrentVideo(c *fiber.Ctx) error {
	user := c.Locals("user").(*database.User)
	courseId := c.Params("courseId")
	record, err := database.GetLastVideoByUserIdAndCourseId(user.ID, courseId)
	if err != nil {
		if err.Error() == "Record not found" {
			video, err := database.GetFistVideoByCourseId(courseId)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			err = database.UpdateVideoViews(video.ID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			_, err = database.CreateHistory(user.ID, video.ID, courseId, "")
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			isFile, err := database.FileExistsByVideoId(video.ID)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(200).JSON(fiber.Map{
				"video":      video,
				"resume":     "",
				"history_id": "",
				"isFile":     isFile,
			})

		} else {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	video, err := database.GetVideoById(record.VideoId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.UpdateVideoViews(video.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	isFile, err := database.FileExistsByVideoId(video.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"video":      video,
		"resume":     record.VideoResume,
		"history_id": record.ID,
		"isFile":     isFile,
	})

}

func UpdateVideo(c *fiber.Ctx) error {
	thumbnail, _ := c.FormFile("thumbnail")
	cleanInput, err := inputs.UpdateVideo(inputs.UpdateVideoInput{
		ID:           c.FormValue("id"),
		Title:        c.FormValue("title"),
		Description:  c.FormValue("description"),
		Duration:     c.FormValue("duration"),
		Thumbnail:    thumbnail,
		OldThumbnail: c.FormValue("old_thumbnail"),
		Video:        c.FormValue("video_tmp"),
		OldVideoHLS:  c.FormValue("old_video_hls"),
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
		fmt.Println("the f error", err)
		return c.Status(500).JSON(fiber.Map{
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
