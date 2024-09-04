package inputs

import (
	"fmt"
	"iron-stream/internal/database"
	"iron-stream/internal/utils"
	"mime/multipart"
)

type FileInput struct {
	VideoID string
	Page    string
	File    *multipart.FileHeader
}

func CreateFile(input FileInput) (database.File, error) {
	if input.VideoID == "" {
		return database.File{}, fmt.Errorf("The video ID is required.")
	}
	_, err := database.GetVideoById(input.VideoID)
	if err != nil {
		return database.File{}, fmt.Errorf("No video found with the id %s.", input.VideoID)
	}
	if input.Page == "" {
		return database.File{}, fmt.Errorf("The page is required.")
	}

	if input.File.Size > MaxFileSize {
		return database.File{}, fmt.Errorf("The file is too large. The maximum size is 10MB.")
	}

	path, err := utils.ManageFile(input.File)
	if err != nil {
		return database.File{}, err
	}

	return database.File{
		VideoID: input.VideoID,
		Page:    input.Page,
		Path:    path,
	}, nil
}
