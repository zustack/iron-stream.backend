package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

func DeleteUserAppsByAppIdAndUserId(userId string, appId string) error {
	result, err := DB.Exec(`DELETE FROM user_apps WHERE app_id = ? AND user_id = ?;`, appId, userId)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func GetEnrolledUserApps(userID, q string) ([]App, error) {
	query := `
		SELECT 
			a.id,
			a.name,
			a.process_name,
			CASE 
				WHEN ua.user_id IS NOT NULL THEN true 
				ELSE false 
			END AS is_user_enrolled
		FROM 
			apps a
		LEFT JOIN 
			user_apps ua
		ON 
			a.id = ua.app_id
			AND ua.user_id = $1
		WHERE 
		  LOWER(a.name) LIKE '%' || LOWER(?) || '%' 
			OR LOWER(a.process_name) LIKE '%' || LOWER(?) || '%';`

	rows, err := DB.Query(query, userID, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apps []App

	for rows.Next() {
		var app App
		err := rows.Scan(
			&app.ID,
			&app.Name,
			&app.ProcessName,
			&app.IsUserEnrolled,
		)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return apps, nil
}

func CreateUserApp(userId, appId string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var userExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`,
		userId).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if !userExists {
		return fmt.Errorf("The user with id %s does not exist", userId)
	}

	var appExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM apps WHERE id = ?)`,
		appId).Scan(&appExists)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if !appExists {
		return fmt.Errorf("The app with id %s does not exist", appId)
	}

	date := utils.FormattedDate()
	_, err = tx.Exec(`
    INSERT INTO user_apps (user_id, app_id, created_at) VALUES (?, ?, ?)`,
		userId, appId, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}

func GetUserApps(userId int64) ([]App, error) {
	var apps []App
	query := `
        SELECT a.name, a.process_name, a.execute_always 
        FROM user_apps ua
        JOIN apps a ON ua.app_id = a.id
        WHERE ua.user_id = ?;
    `
	rows, err := DB.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ProcessName, &a.Name, &a.ExecuteAlways); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		apps = append(apps, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return apps, nil
}
