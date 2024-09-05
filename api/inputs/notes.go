package inputs

import (
	"fmt"
	"iron-stream/internal/database"
	"math"
)

type CreateNotePayload struct {
	Body       string  `json:"body"`
	VideoTitle string  `json:"video_title"`
	Time       float64 `json:"time"`
}

func CreateNote(input CreateNotePayload) (database.Note, error) {
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

	minutes := int(math.Floor(input.Time / 60))
	remainingSeconds := int(math.Floor(math.Mod(input.Time, 60)))
	noteTime := fmt.Sprintf("%02d:%02d", minutes, remainingSeconds)

	return database.Note{
		Body:       input.Body,
		VideoTitle: input.VideoTitle,
		Time:       noteTime,
	}, nil
}
