package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type SpecialApp struct {
	ID          int64  `json:"id"`
	UserId      int64  `json:"user_id"`
	Name        string `json:"name"`
	ProcessName string `json:"process_name"`
	Os          string `json:"os"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

func DeleteAllSpecialAppsByUserId(userId int64) error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM special_apps WHERE user_id = ?", userId).Scan(&count)
	if err != nil {
		return fmt.Errorf("GetAppsCount: %v", err)
	}
	if count == 0 {
		return nil
	}

	result, err := DB.Exec("DELETE FROM special_apps WHERE user_id = ?", userId)
	if err != nil {
		return fmt.Errorf("DeleteAppByID: app id: %v: %v", userId, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteAppByID: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("DeleteAppByID: no app found with ID: %v", userId)
	}
	return nil
}

func GetSpecialAppsByUserId(os string, userId int64) ([]SpecialApp, error) {
	var apps []SpecialApp
	rows, err := DB.Query(`SELECT id, user_id, process_name, name
		FROM special_apps WHERE user_id = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a SpecialApp
		if err := rows.Scan(&a.ID, &a.UserId, &a.ProcessName, &a.Name); err != nil {
			return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
		}
		apps = append(apps, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
	}

	return apps, nil
}

func CreateSpecialApp(a SpecialApp) (int64, error) {
	date := utils.FormattedDate()
	result, err := DB.Exec(`
  INSERT INTO special_apps
  (user_id, name, process_name, os, is_active, created_at) 
  VALUES (?, ?, ?, ?, ?, ?)`,
		a.UserId, a.Name, a.ProcessName, a.Os, a.IsActive, date)

	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	return id, nil
}

type App struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProcessName string `json:"process_name"`
	Os          string `json:"os"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

func GetAppsByOs(os string) ([]App, error) {
	var apps []App
	rows, err := DB.Query(`SELECT * FROM apps WHERE os = ?`, os)
	if err != nil {
		return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ID, &a.Name, &a.ProcessName, &a.Os, &a.IsActive, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
		}
		apps = append(apps, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetAppsByOsAndIsActive: %v", err)
	}

	return apps, nil
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

func GetApps(searchParam, isActiveParam string) ([]App, error) {
	var apps []App
	var args []interface{}
	query := `SELECT * FROM apps WHERE 
              (name LIKE ? OR process_name LIKE ? OR os LIKE ?)`

	args = append(args, searchParam, searchParam, searchParam)

	if isActiveParam != "" {
		query += ` AND is_active = ?`
		isActive := isActiveParam == "1"
		args = append(args, isActive)
	}

	query += ` ORDER BY id DESC `

	rows, err := DB.Query(query, args...)
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
	date := utils.FormattedDate()
	result, err := DB.Exec(`
  INSERT INTO apps
  (name, process_name, os, is_active, created_at) 
  VALUES (?, ?, ?, ?, ?)`,
		a.Name, a.ProcessName, a.Os, a.IsActive, date)

	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	return id, nil
}
