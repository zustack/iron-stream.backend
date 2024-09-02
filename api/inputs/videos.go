package inputs

import (
	"fmt"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"mime/multipart"
)

func CleanUpdateVideoInput(input database.Video) (database.Video, error) {
	if len(input.Title) > 55 {
		return database.Video{}, fmt.Errorf("El título no debe tener más de 55 caracteres.")
	}

	if len(input.Description) > 480 {
		return database.Video{}, fmt.Errorf("La descripción no debe tener más de 480 caracteres.")
	}

	if len(input.VideoHLS) > 255 {
		return database.Video{}, fmt.Errorf("El video no debe tener más de 255 caracteres.")
	}

	if len(input.Length) > 155 {
		return database.Video{}, fmt.Errorf("La duración no debe tener más de 155 caracteres.")
	}

	return database.Video{
		Title:       input.Title,
		Description: input.Description,
		VideoHLS:    input.VideoHLS,
		Length:      input.Length,
	}, nil
}

type CreateVideoInput struct {
	Title       string
	Description string
	Author      string
  CourseID    string
	Thumbnail   *multipart.FileHeader
	Video string
	Duration    string
}

func CreateVideo(input CreateVideoInput) (database.Video, error) {
  _, err := database.GetCourseById(input.CourseID)
  if err != nil {
    return database.Video{}, err
  }

	if input.Title == "" {
		return database.Video{}, fmt.Errorf("The title is required.")
	}
	if len(input.Title) > 50 {
		return database.Video{}, fmt.Errorf("The title should not have more than 50 characters.")
	}

	if input.Description == "" {
		return database.Video{}, fmt.Errorf("The description is required.")
	}
	if len(input.Description) > 150 {
		return database.Video{}, fmt.Errorf("The title should not have more than 150 characters.")
	}

	if input.Duration == "" {
		return database.Video{}, fmt.Errorf("The duration is required.")
	}
	if len(input.Duration) > 20 {
		return database.Video{}, fmt.Errorf("The title should not have more than 20 characters.")
	}

	if input.CourseID == "" {
		return database.Video{}, fmt.Errorf("The course id is required.")
	}

  length, err := utils.GetVideoLength(input.Video)
  if err != nil {
    return database.Video{}, err
  }

	if input.Thumbnail.Size > MaxFileSize {
		return database.Video{}, fmt.Errorf("The thumbnail is too large. The maximum size is 10MB.")
	}

	thumbnail, err := utils.ManageThumbnail(input.Thumbnail)
	if err != nil {
		return database.Video{}, err
	}

  if input.Video == "" {
    return database.Video{}, fmt.Errorf("The video tmp path is required.")
  }

	video, err := utils.ManageVideos(input.Video, input.CourseID)
	if err != nil {
		return database.Video{}, err
	}

	return database.Video{
		Title:       input.Title,
		Description: input.Description,
		Duration:    input.Duration,
    CourseID:    input.CourseID,
    Thumbnail:   thumbnail,
    VideoHLS:    video,
    Length:      length,
	}, nil
}
