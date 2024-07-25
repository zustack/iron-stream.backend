package database

import (
	"database/sql"
	"fmt"
)

type App struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProcessName string `json:"process_name"`
	Os          string `json:"os"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

func GetAppByID(id string) (App, error) {
	var a App
	row := DB.QueryRow(`SELECT * FROM apps WHERE id = ?`, id)
	if err := row.Scan(&a.ID, &a.Name, &a.ProcessName, &a.Os, &a.IsActive, &a.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return a, fmt.Errorf("GetAppByID: %s: no such app", id)
		}
		return a, fmt.Errorf("GetAppByID: %s: %v", id, err)
	}
	return a, nil
}

func GetAppsCount() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM apps").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("GetAppsCount: %v", err)
	}
	return count, nil
}

func GetApps(searchParam, isActiveParam string, limit, cursor int) ([]App, error) {
	var apps []App
	rows, err := DB.Query(`SELECT * FROM apps WHERE 
  name LIKE ? OR process_name LIKE ? OR os LIKE ? OR is_active LIKE ? 
  ORDER BY id DESC LIMIT ? OFFSET ?`,
		searchParam, searchParam, searchParam, isActiveParam, limit, cursor)
	if err != nil {
		return nil, fmt.Errorf("GetApps: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ID, &a.Name, &a.ProcessName, &a.Os, &a.IsActive, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetApps: %v", err)
		}
		apps = append(apps, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetApps: %v", err)
	}
	return apps, nil
}

func UpdateApp(app App) error {
	result, err := DB.Exec(`UPDATE apps SET 
  name = ?, process_name = ? , is_active = ?, os = ? WHERE id = ?`,
		app.Name, app.ProcessName, app.IsActive, app.Os, app.ID)
	if err != nil {
		return fmt.Errorf("UpdateApp: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateApp: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("UpdateApp: no app found with ID: %d", app.ID)
	}
	return nil
}

func GetAppsByOsAndIsActive(os string) ([]App, error) {
	var apps []App
	rows, err := DB.Query(`SELECT process_name, name
		FROM apps WHERE os = ? AND is_active = ?`, os, true)
	if err != nil {
		return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ProcessName, &a.Name); err != nil {
			return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
		}
		apps = append(apps, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
	}

	return apps, nil
}

func DeleteAppByID(id string) error {
	result, err := DB.Exec("DELETE FROM apps WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteAppByID: app id: %s: %v", id, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteAppByID: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("DeleteAppByID: no app found with ID: %s", id)
	}
	return nil
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
