package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

func DeleteUserAppsByAppIdAndUserId(userId string, appId string) error {
	_, err := DB.Exec(`DELETE FROM user_apps WHERE app_id = ? AND user_id = ?;`, appId, userId)
	if err != nil {
    return fmt.Errorf("DeleteUserAppsByAppIdAndUserId: %v", err)
	}
	return nil
}

func GetUserAppsIds(userID string) ([]int64, error) {
	var appsIDs []int64
	rows, err := DB.Query(`
		SELECT uc.app_id FROM user_apps uc
		WHERE uc.user_id = ?;
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserCourseIDs: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var appId int64
		if err := rows.Scan(&appId); err != nil {
			return nil, fmt.Errorf("GetUserCourseIDs: %v", err)
		}
		appsIDs = append(appsIDs, appId)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetUserCourseIDs: %v", err)
	}

	return appsIDs, nil
}


func CreateUserApp(userId, appId string) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("CreateUserApp: %v", err)
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
		return fmt.Errorf("CreateUserCourse: %v", err)
	}
	if !userExists {
		return fmt.Errorf("CreateUserCourse: userID %s no existe", userId)
	}

	var courseExists bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM apps WHERE id = ?)`,
		appId).Scan(&courseExists)
	if err != nil {
		return fmt.Errorf("CreateUserCourse: %v", err)
	}
	if !courseExists {
		return fmt.Errorf("CreateUserCourse: courseID %s no existe", appId)
	}

	date := utils.FormattedDate()
	_, err = tx.Exec(`
    INSERT INTO user_apps (user_id, app_id, created_at) VALUES (?, ?, ?)`,
		userId, appId, date)
	if err != nil {
		return fmt.Errorf("CreateUserCourse: %v", err)
	}

	return nil
}
