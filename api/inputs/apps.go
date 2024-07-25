package inputs

import (
	"fmt"
	"iron-stream/internal/database"
)

func CleanAppInput(input database.App) (database.App, error) {
	if input.Name == "" {
		return database.App{}, fmt.Errorf("El nombre es requerido.")
	}

	if len(input.Name) > 55 {
		return database.App{}, fmt.Errorf("El nombre no debe tener más de 55 caracteres.")
	}

	if input.ProcessName == "" {
		return database.App{}, fmt.Errorf("El nombre del proceso es requerido.")
	}

	if len(input.ProcessName) > 55 {
		return database.App{}, fmt.Errorf("El nombre del proceso no debe tener más de 55 caracteres.")
	}

	if input.Os == "" {
		return database.App{}, fmt.Errorf("El sistema operativo es requerido.")
	}

	if len(input.Os) > 55 {
		return database.App{}, fmt.Errorf("El sistema operativo no debe tener más de 55 caracteres.")
	}

	return database.App{
		Name:        input.Name,
		ProcessName: input.ProcessName,
		Os:          input.Os,
		IsActive:    input.IsActive,
	}, nil
}
