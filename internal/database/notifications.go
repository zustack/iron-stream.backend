package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

type Notification struct {
	ID        int64  `json:"id"`
	NType     string `json:"n_type"`
	Info      string `json:"info"`
	CreatedAt string `json:"created_at"`
}

func DeleteNotification(info string) error {
	result, err := DB.Exec("DELETE FROM notifications WHERE info = ?", info)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No notification found with the info %s", info)
	}
	return nil
}

func GetInfoNotifications(nType string) ([]Notification, error) {
	var ns []Notification
	rows, err := DB.Query(`SELECT info FROM notifications WHERE n_type = ?`, nType)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.Info); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		ns = append(ns, n)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return ns, nil
}

func GetNotifications(nType string) (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM notifications WHERE n_type = ?", nType).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return count, nil
}

func CreateNotification(nType, info string) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`
  INSERT INTO notifications (n_type, info, created_at) VALUES (?, ?, ?)`,
		nType, info, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return nil
}
