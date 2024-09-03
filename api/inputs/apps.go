package inputs

import (
	"fmt"
	"iron-stream/internal/database"
)

func UpdateApp(input database.App) (database.App, error) {
	if input.ID == 0 {
		return database.App{}, fmt.Errorf("The ID is required.")
	}

	if input.Name == "" {
		return database.App{}, fmt.Errorf("The name is required.")
	}
	if len(input.Name) > 55 {
		return database.App{}, fmt.Errorf("The name should not have more than 55 characters.")
	}

	if input.ProcessName == "" {
		return database.App{}, fmt.Errorf("The process name is required.")
	}
	if len(input.ProcessName) > 55 {
		return database.App{}, fmt.Errorf("The process name should not have more than 55 characters.")
	}

	return database.App{
		ID:          input.ID,
		Name:        input.Name,
		ProcessName: input.ProcessName,
	}, nil
}

func CreateApp(input database.App) (database.App, error) {
	if input.Name == "" {
		return database.App{}, fmt.Errorf("The name is required.")
	}
	if len(input.Name) > 55 {
		return database.App{}, fmt.Errorf("The name should not have more than 55 characters.")
	}

	if input.ProcessName == "" {
		return database.App{}, fmt.Errorf("The process name is required.")
	}
	if len(input.ProcessName) > 55 {
		return database.App{}, fmt.Errorf("The process name should not have more than 55 characters.")
	}

	return database.App{
		Name:        input.Name,
		ProcessName: input.ProcessName,
	}, nil
}
