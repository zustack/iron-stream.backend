package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

type App struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	ProcessName    string `json:"process_name"`
	IsActive       bool   `json:"is_active"`
	ExecuteAlways  bool   `json:"execute_always"`
	CreatedAt      string `json:"created_at"`
	IsUserEnrolled bool   `json:"is_user_enrolled"`
}

func UpdateAppStatus(id, isActive string) error {
	result, err := DB.Exec(`UPDATE apps SET is_active = ? WHERE id = ?`,
		isActive, id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No apps found with the id %v", id)
	}
	return nil
}

func GetActiveApps() ([]App, error) {
	var apps []App
	rows, err := DB.Query(`SELECT process_name, name
		FROM apps WHERE is_active = ?`, true)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ProcessName, &a.Name); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		apps = append(apps, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return apps, nil
}

func GetAdminApps(searchParam, isActiveParam string) ([]App, error) {
	var apps []App
	var args []interface{}
	query := `SELECT * FROM apps WHERE 
              (name LIKE ? OR process_name LIKE ?)`

	args = append(args, searchParam, searchParam)

	if isActiveParam != "" {
		query += ` AND is_active = ?`
		isActive := isActiveParam == "1"
		args = append(args, isActive)
	}

	query += ` ORDER BY id DESC`

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ID, &a.Name, &a.ProcessName, &a.IsActive, &a.ExecuteAlways, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		apps = append(apps, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return apps, nil
}

func UpdateApp(app App) error {
	result, err := DB.Exec(`UPDATE apps SET 
  name = ?, process_name = ? , is_active = ?, execute_always = ? WHERE id = ?`,
		app.Name, app.ProcessName, app.IsActive, app.ExecuteAlways, app.ID)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No app found with the id %v", app.ID)
	}
	return nil
}

func DeleteAppByID(id string) error {
	result, err := DB.Exec("DELETE FROM apps WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No app found with the id %s", id)
	}
	return nil
}

func CreateApp(a App) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`
  INSERT INTO apps
  (name, process_name, is_active, execute_always, created_at) 
  VALUES (?, ?, ?, ?, ?)`,
		a.Name, a.ProcessName, a.IsActive, a.ExecuteAlways, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
