package inputs

import (
	"fmt"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"mime/multipart"
	"strconv"
)

const MaxFileSize = 10 * 1024 * 1024

type CreateCourseInput struct {
	Title       string
	Description string
	Author      string
	Thumbnail   *multipart.FileHeader
	Preview     string
	Duration    string
	IsActive    string
}

func CourseInput(input CreateCourseInput) (database.Course, error) {
	if input.Title == "" {
		return database.Course{}, fmt.Errorf("The title is required.")
	}
	if len(input.Title) > 55 {
		return database.Course{}, fmt.Errorf("The title should not have more than 55 characters.")
	}

	if input.Description == "" {
		return database.Course{}, fmt.Errorf("La descripción es requerido.")
	}
	if len(input.Description) > 480 {
		return database.Course{}, fmt.Errorf("La descripción no debe tener más de 480 caracteres.")
	}

	if input.Author == "" {
		return database.Course{}, fmt.Errorf("La autor es requerido.")
	}
	if len(input.Author) > 25 {
		return database.Course{}, fmt.Errorf("El autor no debe tener más de 25 caracteres.")
	}

	if input.Duration == "" {
		return database.Course{}, fmt.Errorf("La duración es requerido.")
	}
	if len(input.Duration) > 25 {
		return database.Course{}, fmt.Errorf("La duración no debe tener más de 25 caracteres.")
	}

	isActiveBool, err := strconv.ParseBool(input.IsActive)
	if err != nil {
		return database.Course{}, fmt.Errorf(err.Error())
	}

	if input.Thumbnail.Size > MaxFileSize {
		return database.Course{}, fmt.Errorf("The thumbnail is too large. The maximum size is 10MB.")
	}

	thumbnail, err := utils.ManageThumbnail(input.Thumbnail)
	if err != nil {
		return database.Course{}, err
	}

	if len(input.Preview) > 1000 {
		return database.Course{}, fmt.Errorf("The preview is too long. Make sure that the file path existis.")
	}

	if input.Preview != "" {
		preview, err := utils.ManagePreviews(input.Preview)
		if err != nil {
			return database.Course{}, err
		}
		input.Preview = preview
	}

	return database.Course{
		Title:       input.Title,
		Description: input.Description,
		Author:      input.Author,
		Duration:    input.Duration,
		IsActive:    isActiveBool,
		Thumbnail:   thumbnail,
		Preview:     input.Preview,
	}, nil
}

func CleanUpdateCourseInput(input database.Course) (database.Course, error) {
	if len(input.Title) > 55 {
		return database.Course{}, fmt.Errorf("El título no debe tener más de 55 caracteres.")
	}
	if len(input.Description) > 480 {
		return database.Course{}, fmt.Errorf("La descripción no debe tener más de 480 caracteres.")
	}
	if len(input.Author) > 25 {
		return database.Course{}, fmt.Errorf("El autor no debe tener más de 25 caracteres.")
	}
	if len(input.Duration) > 25 {
		return database.Course{}, fmt.Errorf("La duración no debe tener más de 25 caracteres.")
	}
	return database.Course{
		Title:       input.Title,
		Description: input.Description,
		Author:      input.Author,
		Duration:    input.Duration,
	}, nil
}
