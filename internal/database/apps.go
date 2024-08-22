package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type App struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ProcessName string `json:"process_name"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

func GetAppsByIds(ids []int64) ([]App, error) {
  var apps []App
  for _, id := range ids {
    var a App
    row := DB.QueryRow(`SELECT name, process_name FROM apps WHERE id = ?`, id)
    if err := row.Scan(&a.Name, &a.ProcessName); err != nil {
      if err == sql.ErrNoRows {
        return nil, fmt.Errorf("GetAppByID: %v: no such app", id)
      }
      return nil, fmt.Errorf("GetAppByID: %v: %v", id, err)
    }
    apps = append(apps, a)
  }
  return apps, nil
}

func UpdateAppStatus(id, isActive string) error {
	result, err := DB.Exec(`UPDATE apps SET is_active = ? WHERE id = ?`,
		isActive, id)
	if err != nil {
		return fmt.Errorf("UpdateApp: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateApp: error getting rows affected %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("UpdateApp: no app found with ID: %s", id)
	}
	return nil
}


func GetAppByID(id string) (App, error) {
	var a App
	row := DB.QueryRow(`SELECT * FROM apps WHERE id = ?`, id)
	if err := row.Scan(&a.ID, &a.Name, &a.ProcessName, &a.IsActive, &a.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return a, fmt.Errorf("GetAppByID: %s: no such app", id)
		}
		return a, fmt.Errorf("GetAppByID: %s: %v", id, err)
	}
	return a, nil
}

// TODO fix the is_active NOT WORKING!
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
		return nil, fmt.Errorf("GetApps: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ID, &a.Name, &a.ProcessName, &a.IsActive, &a.CreatedAt); err != nil {
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
  name = ?, process_name = ? , is_active = ? WHERE id = ?`,
		app.Name, app.ProcessName, app.IsActive, app.ID)
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

func GetActiveApps() ([]App, error) {
	var apps []App
	rows, err := DB.Query(`SELECT process_name, name
		FROM apps WHERE is_active = ?`, true)
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
  (name, process_name, is_active, created_at) 
  VALUES (?, ?, ?, ?)`,
		a.Name, a.ProcessName, a.IsActive, date)

	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateApp: %v", err)
	}

	return id, nil
}
