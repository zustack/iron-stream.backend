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

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createHistoryInput struct {
	Id string `json:"id"`
	VideoId string `json:"video_id"`
  CurrentVideoId int64 `json:"current_video_id"`
	CourseId string `json:"course_id"`
  Resume string `json:"resume"`
}

func WatchNewVideo(c *fiber.Ctx) error {
  user := c.Locals("user").(*database.User)
  var payload createHistoryInput 
  if err := c.BodyParser(&payload); err != nil {
    fmt.Println("fucking error", err)
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "No se pudo procesar la solicitud.",
    })
  }
  // actualizo el ultimo registro del historial con el resume del video

  // !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!! 
  // deberia actualizar el history de todos los registros del usuario que tengan el mismo video_id y user_id
  // payload.id es el history_id que viene del video actual
  // payload.resume es el video actual
  // user.ID
  // payload.CurrentVideoId
  // saco esto? y al buscar los registros del usuario solo traigo los ultimos recolectados y unicos?
  // cuanto cuesta la escritura de abajo?
  /*
  err := database.UpdateHistory(payload.Resume, user.ID, payload.CurrentVideoId)
  if err != nil {
    fmt.Println("er1", err)
    return c.SendStatus(500)
  }
  */
  // actualizo el ultimo registro del historial con el resume del video
  err := database.UpdateHistory(payload.Id, payload.Resume)
  if err != nil {
    fmt.Println("er1", err)
    return c.SendStatus(500)
  }

  // find the last record with the same user.ID, payload.VideoId, payload.CourseId
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
  user := c.Locals("user").(*database.User)
  course_id := c.Params("course_id")
  // obtener el utlimo video del historial
  record, err := database.GetLastVideoByUserIdAndCourseId(user.ID, course_id)
  // es el primer video que el usuario ve
  if err != nil {
    // obtener el primer video del curso
    video, err := database.GetFistVideoByCourseId(course_id)
    if err != nil {
      return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": err.Error(),
      })
    }
    // agrega un video view
    err = database.UpdateVideoViews(video.ID)
    if err != nil {
      return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": err.Error(),
      })
    }
    // crear nuevo record con el primer video
    // user_id int64, video_id int64, course_id int64, resume string
    videoStr := fmt.Sprintf("%d", video.ID)
    newRecord, err := database.CreateHistory(user.ID, videoStr, course_id, "")
    if err != nil {
      return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
        "error": err.Error(),
      })
    }
    fmt.Println("first record", newRecord)
    return c.JSON(video)
  }

  // nesesito el video y el record para el historial, sobre todo para saber donde estaba el video
  video, err := database.GetVideoById(record.VideoId)

  // agrega un video view
  err = database.UpdateVideoViews(video.ID)
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  return c.Status(fiber.StatusOK).JSON(fiber.Map{
    "video": video,
    "resume": record.VideoResume,
    "history_id": record.ID,
  })
}



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
  user := c.Locals("user").(*database.User)
  // este history deberia retornar items unicos, no duplicados!! mucho alloc de memoria!!!!
  // ok deberia estar bien ahora :)
  history, err := database.GetUserUniqueHistory(user.ID)
  if err != nil {
    fmt.Println("1", err)
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

  historyMap := make(map[int64]string)
  for _, h := range history {
    historyMap[h.VideoId] = h.VideoResume
    fmt.Println("resume", h.VideoResume)
  }

  id := c.Params("id")
  if id == "" {
    fmt.Println("2", err)
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": "Course ID is required",
    })
  }

  cursor, err := strconv.Atoi(c.Query("cursor", "0"))
  if err != nil {
    fmt.Println("3", err)
    return c.Status(400).SendString("Invalid cursor")
  }
  limit := 50 

  searchParam := c.Query("q", "")
  searchParam = "%" + searchParam + "%"

  videos, err := database.GetVideos(id, searchParam, cursor, limit)
  if err != nil {
    fmt.Println("4", err)
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
      "error": err.Error(),
    })
  }

for i := range videos {
		if resume, ok := historyMap[videos[i].ID]; ok {
			parts := strings.Split(resume, ".")
			resumeInt, err1 := strconv.Atoi(parts[0])
			if err1 != nil {
				fmt.Println("Error al convertir la parte entera:", err1)
				continue
			}

			lengthInt, err2 := strconv.Atoi(videos[i].Length)
			if err2 != nil {
				fmt.Println("Error al convertir la longitud:", err2)
				continue
			}

			if lengthInt == 0 {
				fmt.Println("Longitud no puede ser cero")
				continue
			}
			result := float64(resumeInt) / float64(lengthInt) * 100

			videos[i].VideoResume = fmt.Sprintf("%.0f", result)
			fmt.Println("Nuevo resumen:", videos[i].VideoResume)
		} else {
			videos[i].VideoResume = "0"
		}
	}

  totalCount, err := database.GetVideosCount()
  if err != nil {
    fmt.Println("5", err)
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

  // get the lenght of the video
	getLengthCmd := exec.Command("sh", filepath.Join(os.Getenv("ROOT_PATH"), "get-video-length.sh"), cleanInput.VideoHLS)
  output, err := getLengthCmd.Output()
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": err.Error(),
    })
  }
  length := strings.Trim(string(output), "\n")
  fmt.Println("length: ", length)
  cleanInput.Length = length

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
