package database

import "fmt"

type App struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProcessName string `json:"process_name"`
	Os          string `json:"os"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

func CreateApp(a App) (int64, error) {
	result, err := DB.Exec(`
  INSERT INTO apps
  (name, process_name, os, is_active) 
  VALUES (?, ?, ?, ?)`,
		a.Name, a.ProcessName, a.Os, a.IsActive)

	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	return id, nil
}
