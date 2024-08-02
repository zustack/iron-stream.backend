package inputs

import (
	"fmt"
	"iron-stream/internal/database"
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

func CleanVideoInput(input database.Video) (database.Video, error) {
	if input.Title == "" {
		return database.Video{}, fmt.Errorf("El título es requerido.")
	}
	if len(input.Title) > 55 {
		return database.Video{}, fmt.Errorf("El título no debe tener más de 55 caracteres.")
	}

	if input.Description == "" {
		return database.Video{}, fmt.Errorf("La descripción es requerido.")
	}
	if len(input.Description) > 480 {
		return database.Video{}, fmt.Errorf("La descripción no debe tener más de 480 caracteres.")
	}

	if input.VideoHLS == "" {
		return database.Video{}, fmt.Errorf("El video es requerido.")
	}
	if len(input.VideoHLS) > 255 {
		return database.Video{}, fmt.Errorf("El video no debe tener más de 255 caracteres.")
	}

	if input.Length == "" {
		return database.Video{}, fmt.Errorf("La duración es requerida.")
	}
	if len(input.VideoHLS) > 155 {
		return database.Video{}, fmt.Errorf("La duración no debe tener más de 155 caracteres.")
	}

	return database.Video{
		Title:       input.Title,
		Description: input.Description,
    VideoHLS:    input.VideoHLS,
    Length:      input.Length,
	}, nil
}
