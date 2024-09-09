package database

import (
	"fmt"
	"iron-stream/internal/utils"
)

type Policy struct {
	ID        int64  `json:"id"`
	Content string `json:"content"`
	PType string `json:"p_type"`
	CreatedAt string `json:"created_at"`
}

func DeletePolicy(id string) error {
	result, err := DB.Exec("DELETE FROM policy WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No policy found with the id %v", id)
	}
	return nil
}

func GetPolicy() ([]Policy, error) {
	var ps []Policy
	rows, err := DB.Query("SELECT * FROM policy")
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Policy
		if err := rows.Scan(&p.ID, &p.Content, &p.PType, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return ps, nil
}

func CreatePolicy(content, p_type string) error {
	date := utils.FormattedDate()
	_, err := DB.Exec(`INSERT INTO policy 
  (content, p_type, created_at) VALUES 
    (?, ?, ?)`, content, p_type, date)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return nil
}
