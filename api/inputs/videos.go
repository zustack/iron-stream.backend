package inputs

import (
	"fmt"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"mime/multipart"
)

type UpdateVideoInput struct {
  ID string
	Title       string
	Description string
  Duration    string
	Thumbnail   *multipart.FileHeader
  OldThumbnail string
	Video string
  OldVideoHLS string
}

func UpdateVideo(input UpdateVideoInput) (database.Video, error) {
	if input.ID == "" {
		return database.Video{}, fmt.Errorf("The video id is required.")
	}

  v, err := database.GetVideoById(input.ID)
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
		return database.Video{}, fmt.Errorf("The duration should not have more than 20 characters.")
	}

  var length string
  var video string
  if input.Video != "" {
    length, err = utils.GetVideoLength(input.Video)
    if err != nil {
      return database.Video{}, err
    }

	  video, err = utils.ManageVideos(input.Video, v.CourseID)
	  if err != nil {
		  return database.Video{}, err
	  }
  } else {
    video = input.OldVideoHLS
    length = v.Length
  }

	var thumbnail string
	if input.Thumbnail != nil {
		if input.Thumbnail.Size > MaxFileSize {
			return database.Video{}, fmt.Errorf("The thumbnail is too large. The maximum size is 10MB.")
		}
		thumbnail, err = utils.ManageThumbnail(input.Thumbnail)
		if err != nil {
			return database.Video{}, err
		}
	} else {
		thumbnail = input.OldThumbnail
	}

	return database.Video{
    ID:          v.ID,
    Title:       input.Title,
    Description: input.Description,
    Duration:    input.Duration,
    Thumbnail:   thumbnail,
    VideoHLS:    video,
    Length:      length,
	}, nil
}

type CreateVideoInput struct {
	Title       string
	Description string
  CourseID    string
	Thumbnail   *multipart.FileHeader
	Video string
	Duration    string
}

func CreateVideo(input CreateVideoInput) (database.Video, error) {
	if input.CourseID == "" {
		return database.Video{}, fmt.Errorf("The course id is required.")
	}

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
