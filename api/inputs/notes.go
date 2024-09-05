package inputs

import (
	"fmt"
	"iron-stream/internal/database"
)

func CreateNote(input database.Note) (database.Note, error) {
	if input.Body == "" {
		return database.Note{}, fmt.Errorf("The body is required.")
	}
	if len(input.Body) > 255 {
		return database.Note{}, fmt.Errorf("The body should not have more than 255 characters.")
	}

	if input.VideoTitle == "" {
		return database.Note{}, fmt.Errorf("The video title is required.")
	}
	if len(input.VideoTitle) > 50 {
		return database.Note{}, fmt.Errorf("The video title should not have more than 255 characters.")
	}

	if input.Time == "" {
		return database.Note{}, fmt.Errorf("The video time is required.")
	}
	if len(input.Time) > 50 {
		return database.Note{}, fmt.Errorf("The video time should not have more than 50 characters.")
	}

	return database.Note{
    Body:       input.Body,
    VideoTitle: input.VideoTitle,
    Time:       input.Time,
	}, nil
}
