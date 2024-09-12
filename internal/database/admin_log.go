package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

type AdminLog struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	LType     string `json:"l_type"`
	CreatedAt string `json:"created_at"`
}

func GetAdminLog() ([]AdminLog, error) {
	var als []AdminLog
	rows, err := DB.Query(`SELECT *
		FROM admin_log ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var al AdminLog
		if err := rows.Scan(&al.ID, &al.Content, &al.LType, &al.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		als = append(als, al)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return als, nil
}

func CreateAdminLog(content, l_type string) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`
  INSERT INTO admin_log
  (content, l_type, created_at) 
  VALUES (?, ?, ?)`,
		content, l_type, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
