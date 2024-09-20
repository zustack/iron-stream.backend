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
	Price       string
}

func CleanCreateCourse(input CreateCourseInput) (database.Course, error) {
	if input.Title == "" {
		return database.Course{}, fmt.Errorf("The title is required.")
	}
	if len(input.Title) > 50 {
		return database.Course{}, fmt.Errorf("The title should not have more than 50 characters.")
	}

	price, err := strconv.Atoi(input.Price)
	if err != nil {
		return database.Course{}, fmt.Errorf(err.Error())
	}

	if price <= 0 {
		return database.Course{}, fmt.Errorf("The price is required to be greater than 0.")
	}

	if input.Description == "" {
		return database.Course{}, fmt.Errorf("The description is required.")
	}
	if len(input.Description) > 270 {
		return database.Course{}, fmt.Errorf("The description should not have more than 270 characters.")
	}

	if input.Author == "" {
		return database.Course{}, fmt.Errorf("The author is required.")
	}
	if len(input.Author) > 30 {
		return database.Course{}, fmt.Errorf("The author should not have more than 30 characters.")
	}

	if input.Duration == "" {
		return database.Course{}, fmt.Errorf("The duration is required.")
	}
	if len(input.Duration) > 30 {
		return database.Course{}, fmt.Errorf("The duration should not have more than 30 characters.")
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
		Price:       price,
		Thumbnail:   thumbnail,
		Preview:     input.Preview,
	}, nil
}

type UpdateCourseInput struct {
	ID           string
	Title        string
	Description  string
	Author       string
	Thumbnail    *multipart.FileHeader
	OldThumbnail string
	Preview      string
	OldPreview   string
	Duration     string
	IsActive     string
	Price        string
}

func CleanUpdateCourse(input UpdateCourseInput) (database.Course, error) {
	if input.ID == "" {
		return database.Course{}, fmt.Errorf("The id is required.")
	}
	id, err := strconv.ParseInt(input.ID, 10, 64)
	if err != nil {
		return database.Course{}, fmt.Errorf(err.Error())
	}

	price, err := strconv.Atoi(input.Price)
	if err != nil {
		return database.Course{}, fmt.Errorf(err.Error())
	}

	if price <= 0 {
		return database.Course{}, fmt.Errorf("The price is required to be greater than 0.")
	}

	if input.Title == "" {
		return database.Course{}, fmt.Errorf("The title is required.")
	}
	if len(input.Title) > 50 {
		return database.Course{}, fmt.Errorf("The title should not have more than 50 characters.")
	}

	if input.Description == "" {
		return database.Course{}, fmt.Errorf("The description is required.")
	}
	if len(input.Description) > 270 {
		return database.Course{}, fmt.Errorf("The description should not have more than 270 characters.")
	}

	if input.Author == "" {
		return database.Course{}, fmt.Errorf("The author is required.")
	}
	if len(input.Author) > 30 {
		return database.Course{}, fmt.Errorf("The author should not have more than 30 characters.")
	}

	if input.Duration == "" {
		return database.Course{}, fmt.Errorf("The duration is required.")
	}
	if len(input.Duration) > 30 {
		return database.Course{}, fmt.Errorf("The duration should not have more than 30 characters.")
	}

	isActiveBool, err := strconv.ParseBool(input.IsActive)
	if err != nil {
		return database.Course{}, fmt.Errorf(err.Error())
	}

	var thumbnail string
	if input.Thumbnail != nil {
		if input.Thumbnail.Size > MaxFileSize {
			return database.Course{}, fmt.Errorf("The thumbnail is too large. The maximum size is 10MB.")
		}
		thumbnail, err = utils.ManageThumbnail(input.Thumbnail)
		if err != nil {
			return database.Course{}, err
		}
	} else {
		thumbnail = input.OldThumbnail
	}

	var previewToDB string
	if input.Preview != "" {
		preview, err := utils.ManagePreviews(input.Preview)
		if err != nil {
			return database.Course{}, err
		}
		previewToDB = preview
	} else {
		previewToDB = input.OldPreview
	}

	return database.Course{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Author:      input.Author,
		Duration:    input.Duration,
		IsActive:    isActiveBool,
		Thumbnail:   thumbnail,
		Preview:     previewToDB,
		Price:       price,
	}, nil
}
