package database

import (
	"database/sql"
	"fmt"
	"iron-stream/internal/utils"
)

type Policy struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type PolicyItem struct {
	ID       int64  `json:"id"`
	Body     string `json:"body"`
	PolicyID int64  `json:"policy_id"`
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


func DeletePolicyItem(id string) error {
	result, err := DB.Exec("DELETE FROM policy_item WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No policy item found with the id %v", id)
	}
	return nil
}

func GetPolicyItemsByPolicyId(policyId string) ([]PolicyItem, error) {
	var ps []PolicyItem
	rows, err := DB.Query("SELECT * FROM policy_item WHERE policy_id = ?", policyId)
	if err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p PolicyItem
		if err := rows.Scan(&p.ID, &p.Body, &p.PolicyID); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return ps, nil
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
		if err := rows.Scan(&p.ID, &p.Title); err != nil {
			return nil, fmt.Errorf("An unexpected error occurred: %v", err)
		}
		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("An unexpected error occurred: %v", err)
	}

	return ps, nil
}


func CreatePolicy(title string) error {
  _, err := DB.Exec(`INSERT INTO policy (title) VALUES (?)`, title)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return nil
}

func CreatePolicyItem(body string, policyId string) error {
  _, err := DB.Exec(`INSERT INTO policy_items (body, policy_id) VALUES (?, ?)`, body, policyId)
	if err != nil {
		return fmt.Errorf("An unexpected error occurred: %v", err)
	}
	return nil
}
