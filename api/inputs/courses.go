package inputs

import (
	"fmt"
	"iron-stream/internal/database"
)

func CleanCourseInput(input database.Course) (database.Course, error) {
	if input.Title == "" {
		return database.Course{}, fmt.Errorf("El título es requerido.")
	}
	if len(input.Title) > 55 {
		return database.Course{}, fmt.Errorf("El título no debe tener más de 55 caracteres.")
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

	return database.Course{
		Title:       input.Title,
		Description: input.Description,
		Author:      input.Author,
		Duration:    input.Duration,
	}, nil
}
