package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

type UserLog struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	LType     string `json:"l_type"`
	UserID    int64  `json:"user_id"`
	CreatedAt string `json:"created_at"`
}

func GetUserLog(userID string) ([]UserLog, error) {
	var uls []UserLog
	rows, err := DB.Query(`SELECT *
		FROM user_log WHERE user_id = ?`, userID)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ul UserLog
		if err := rows.Scan(&ul.ID, &ul.Content, &ul.LType, &ul.UserID, &ul.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		uls = append(uls, ul)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return uls, nil
}

func CreateUserLog(content, l_type string, user_id int64) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`
  INSERT INTO user_log
  (content, l_type, user_id, created_at) 
  VALUES (?, ?, ?, ?)`,
		content, l_type, user_id, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
